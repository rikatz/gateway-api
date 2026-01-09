/*
Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

//go:embed validation.go.tmpl
var validationTemplate string

// FieldType represents the type characteristics of a field
type FieldType int

const (
	TypeDirect  FieldType = iota // Direct field (not a pointer or slice)
	TypePointer                  // Pointer field
	TypeSlice                    // Slice field
)

// PathSegment represents one segment in a field path with type information
type PathSegment struct {
	FieldName string
	JSONName  string
	Type      FieldType
	TypeName  string // The Go type name
	IsStruct  bool   // Whether this type is a struct (not a primitive or alias)
}

// ExperimentalField represents a complete experimental field with full path information
type ExperimentalField struct {
	Path     []PathSegment
	JSONPath string
}

// RootType represents a Kubernetes root object type
type RootType struct {
	TypeName           string
	PackageName        string
	ExperimentalFields []ExperimentalField
}

func main() {
	apiVersions := []string{"v1", "v1alpha2", "v1alpha3", "v1beta1"}
	packageMap := make(map[string][]RootType) // packageName -> rootTypes

	for _, version := range apiVersions {
		apiPath := filepath.Join("apis", version)
		if _, err := os.Stat(apiPath); os.IsNotExist(err) {
			continue
		}

		rootTypes, err := parseAPIPackage(apiPath, version)
		if err != nil {
			log.Fatalf("Error parsing %s: %v", apiPath, err)
		}

		if len(rootTypes) > 0 {
			packageMap[version] = rootTypes
		}
	}

	// Generate validation files in each API package
	totalTypes := 0
	for pkgName, rootTypes := range packageMap {
		if err := generatePackageValidationCode(pkgName, rootTypes); err != nil {
			log.Fatalf("Error generating validation code for %s: %v", pkgName, err)
		}
		totalTypes += len(rootTypes)
	}

	log.Printf("Successfully generated experimental validation code for %d types across %d packages", totalTypes, len(packageMap))
}

func parseAPIPackage(pkgPath, pkgName string) ([]RootType, error) {
	fset := token.NewFileSet()

	// Read directory and parse individual files
	entries, err := os.ReadDir(pkgPath)
	if err != nil {
		return nil, err
	}

	var files []*ast.File
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") ||
		   strings.HasSuffix(name, "_test.go") ||
		   strings.HasPrefix(name, "zz_generated") {
			continue
		}

		filePath := filepath.Join(pkgPath, name)
		file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
		}
		files = append(files, file)
	}

	if len(files) == 0 {
		return nil, nil
	}

	var rootTypes []RootType

	// Process files as a single package
	{
		// Collect all type definitions
		typeMap := make(map[string]*TypeInfo)
		for _, file := range files {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							typeMap[typeSpec.Name.Name] = &TypeInfo{
								Spec:       typeSpec,
								IsPointer:  false,
								IsSlice:    false,
								StructType: extractStructType(typeSpec),
							}
						}
					}
				}
			}
		}

		// Find root types
		for _, file := range files {
			roots := findRootTypes(file, fset)
			for _, rootSpec := range roots {
				rootType := analyzeRootType(rootSpec, typeMap, pkgName, fset)
				if len(rootType.ExperimentalFields) > 0 {
					rootTypes = append(rootTypes, rootType)
				}
			}
		}
	}

	return rootTypes, nil
}

// TypeInfo holds information about a type
type TypeInfo struct {
	Spec       *ast.TypeSpec
	IsPointer  bool
	IsSlice    bool
	StructType *ast.StructType
}

func extractStructType(typeSpec *ast.TypeSpec) *ast.StructType {
	switch t := typeSpec.Type.(type) {
	case *ast.StructType:
		return t
	case *ast.StarExpr:
		if st, ok := t.X.(*ast.StructType); ok {
			return st
		}
	}
	return nil
}

func findRootTypes(file *ast.File, fset *token.FileSet) []*ast.TypeSpec {
	var rootTypes []*ast.TypeSpec

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		isRoot := false
		if genDecl.Doc != nil {
			for _, comment := range genDecl.Doc.List {
				if strings.Contains(comment.Text, "+kubebuilder:object:root=true") {
					isRoot = true
					break
				}
			}
		}

		if !isRoot {
			declPos := fset.Position(genDecl.Pos())
			for _, cg := range file.Comments {
				cgPos := fset.Position(cg.End())
				if cgPos.Line < declPos.Line && declPos.Line-cgPos.Line < 50 {
					for _, comment := range cg.List {
						if strings.Contains(comment.Text, "+kubebuilder:object:root=true") {
							isRoot = true
							break
						}
					}
					if isRoot {
						break
					}
				}
			}
		}

		if !isRoot {
			continue
		}

		for _, spec := range genDecl.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				rootTypes = append(rootTypes, typeSpec)
			}
		}
	}

	return rootTypes
}

func analyzeRootType(rootSpec *ast.TypeSpec, typeMap map[string]*TypeInfo, pkgName string, fset *token.FileSet) RootType {
	rootType := RootType{
		TypeName:    rootSpec.Name.Name,
		PackageName: pkgName,
	}

	structType, ok := rootSpec.Type.(*ast.StructType)
	if !ok {
		return rootType
	}

	// Find Spec field
	var specTypeName string
	for _, field := range structType.Fields.List {
		if len(field.Names) > 0 && field.Names[0].Name == "Spec" {
			specTypeName = extractTypeName(field.Type)
			break
		}
	}

	if specTypeName == "" {
		return rootType
	}

	// Find experimental fields in Spec
	if specInfo, ok := typeMap[specTypeName]; ok && specInfo.StructType != nil {
		fields := findExperimentalFieldsWithTypes(specInfo.StructType, typeMap, []PathSegment{}, "spec")
		rootType.ExperimentalFields = fields
	}

	return rootType
}

func findExperimentalFieldsWithTypes(structType *ast.StructType, typeMap map[string]*TypeInfo, path []PathSegment, jsonPath string) []ExperimentalField {
	var fields []ExperimentalField

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			// Embedded struct
			typeName := extractTypeName(field.Type)
			if embeddedInfo, ok := typeMap[typeName]; ok && embeddedInfo.StructType != nil {
				fields = append(fields, findExperimentalFieldsWithTypes(embeddedInfo.StructType, typeMap, path, jsonPath)...)
			}
			continue
		}

		fieldName := field.Names[0].Name
		if !ast.IsExported(fieldName) {
			continue
		}

		jsonTag := extractJSONTag(field)
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Determine field type
		fieldType, baseTypeName := analyzeFieldType(field.Type)

		// Determine if this is a struct type
		isStruct := false
		if nestedInfo, ok := typeMap[baseTypeName]; ok && nestedInfo.StructType != nil {
			isStruct = true
		}

		newSegment := PathSegment{
			FieldName: fieldName,
			JSONName:  jsonTag,
			Type:      fieldType,
			TypeName:  baseTypeName,
			IsStruct:  isStruct,
		}
		newPath := append(path, newSegment)
		newJSONPath := jsonPath + "." + jsonTag

		// Check if experimental
		if hasExperimentalMarker(field) {
			fields = append(fields, ExperimentalField{
				Path:     newPath,
				JSONPath: newJSONPath,
			})
		}

		// ALWAYS recursively check nested types, even if the field itself is experimental
		// This allows us to find experimental fields nested inside other experimental fields
		// (e.g., Gateway.TLS.ClientCertificateRef where both TLS and ClientCertificateRef are experimental)
		if isStruct {
			if nestedInfo, ok := typeMap[baseTypeName]; ok && nestedInfo.StructType != nil {
				nested := findExperimentalFieldsWithTypes(nestedInfo.StructType, typeMap, newPath, newJSONPath)
				fields = append(fields, nested...)
			}
		}
	}

	return fields
}

func analyzeFieldType(expr ast.Expr) (FieldType, string) {
	switch t := expr.(type) {
	case *ast.Ident:
		return TypeDirect, t.Name
	case *ast.StarExpr:
		_, typeName := analyzeFieldType(t.X)
		return TypePointer, typeName
	case *ast.ArrayType:
		_, typeName := analyzeFieldType(t.Elt)
		return TypeSlice, typeName
	case *ast.SelectorExpr:
		return TypeDirect, ""
	}
	return TypeDirect, ""
}

func hasExperimentalMarker(field *ast.Field) bool {
	if field.Doc == nil {
		return false
	}

	for _, comment := range field.Doc.List {
		text := comment.Text
		if strings.Contains(text, "<gateway:experimental>") &&
			!strings.Contains(text, "<gateway:experimental:") {
			return true
		}
	}
	return false
}

func extractTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return extractTypeName(t.X)
	case *ast.ArrayType:
		return extractTypeName(t.Elt)
	case *ast.SelectorExpr:
		return ""
	}
	return ""
}

func extractJSONTag(field *ast.Field) string {
	if field.Tag == nil {
		return ""
	}

	tag := strings.Trim(field.Tag.Value, "`")
	for _, part := range strings.Split(tag, " ") {
		if strings.HasPrefix(part, "json:") {
			jsonTag := strings.TrimPrefix(part, "json:")
			jsonTag = strings.Trim(jsonTag, "\"")
			parts := strings.Split(jsonTag, ",")
			if len(parts) > 0 && parts[0] != "" {
				return parts[0]
			}
		}
	}
	return ""
}

// generatePackageValidationCode generates zz_generated.validation.go in the API package directory
func generatePackageValidationCode(pkgName string, rootTypes []RootType) error {
	sort.Slice(rootTypes, func(i, j int) bool {
		return rootTypes[i].TypeName < rootTypes[j].TypeName
	})

	// Generate validation functions to a buffer
	var validationBuf bytes.Buffer
	for _, rt := range rootTypes {
		writeValidationFunction(&validationBuf, rt)
	}

	// Prepare template data
	templateData := struct {
		PackageName         string
		RootTypes           []RootType
		ValidationFunctions string
	}{
		PackageName:         pkgName,
		RootTypes:           rootTypes,
		ValidationFunctions: validationBuf.String(),
	}

	// Parse and execute template
	tmpl, err := template.New("validation").Parse(validationTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, templateData); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Write to API package directory
	outputPath := filepath.Join("apis", pkgName, "zz_generated.validation.go")
	if err := os.WriteFile(outputPath, output.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputPath, err)
	}

	return nil
}

func writeValidationFunction(w io.Writer, rt RootType) {
	checkFuncName := fmt.Sprintf("check%sFields", rt.TypeName)
	validateFuncName := fmt.Sprintf("Validate%sExperimentalFields", rt.TypeName)

	// Write the check function that returns []string of experimental fields
	fmt.Fprintf(w, "// %s returns a list of experimental fields found in %s\n", checkFuncName, rt.TypeName)
	fmt.Fprintf(w, "func %s(obj *%s) []string {\n", checkFuncName, rt.TypeName)
	fmt.Fprintf(w, "\tvar experimentalFields []string\n\n")

	for _, field := range rt.ExperimentalFields {
		lastSeg := field.Path[len(field.Path)-1]
		// Skip non-pointer struct fields - their nested experimental fields are checked separately
		if lastSeg.Type != TypePointer && lastSeg.IsStruct {
			continue
		}
		writeFieldCheckCode(w, field)
	}

	fmt.Fprintf(w, "\treturn experimentalFields\n")
	fmt.Fprintf(w, "}\n\n")

	// Write the exported validation function that calls the generic validator
	fmt.Fprintf(w, "// %s checks for experimental fields in %s.\n", validateFuncName, rt.TypeName)
	fmt.Fprintf(w, "// Returns an error listing all experimental fields that are in use, or nil if no experimental fields are used.\n")
	fmt.Fprintf(w, "func %s(obj *%s) error {\n", validateFuncName, rt.TypeName)
	fmt.Fprintf(w, "\treturn validateExperimentalFields(obj, %s)\n", checkFuncName)
	fmt.Fprintf(w, "}\n\n")
}

func writeFieldCheckCode(w io.Writer, field ExperimentalField) {
	fmt.Fprintf(w, "\t// Check %s\n", field.JSONPath)

	// Find all slice positions in the path
	sliceIndices := []int{}
	for i, seg := range field.Path {
		if seg.Type == TypeSlice {
			sliceIndices = append(sliceIndices, i)
		}
	}

	if len(sliceIndices) == 0 {
		// No slices, direct access
		accessor := buildAccessor("obj.Spec", field.Path)
		lastSeg := field.Path[len(field.Path)-1]

		// Build nil guards for all parent pointers
		nilGuards := buildNilGuards("obj.Spec", field.Path[:len(field.Path)-1])

		if lastSeg.Type == TypePointer {
			if len(nilGuards) > 0 {
				fmt.Fprintf(w, "\tif %s && %s != nil {\n", strings.Join(nilGuards, " && "), accessor)
			} else {
				fmt.Fprintf(w, "\tif %s != nil {\n", accessor)
			}
			fmt.Fprintf(w, "\t\texperimentalFields = append(experimentalFields, \"%s\")\n", field.JSONPath)
			fmt.Fprintf(w, "\t}\n\n")
		} else if !lastSeg.IsStruct {
			// Non-pointer primitive/alias field - check for non-zero value
			if len(nilGuards) > 0 {
				fmt.Fprintf(w, "\tif %s && %s != \"\" { // Non-zero check\n", strings.Join(nilGuards, " && "), accessor)
			} else {
				fmt.Fprintf(w, "\tif %s != \"\" { // Non-zero check\n", accessor)
			}
			fmt.Fprintf(w, "\t\texperimentalFields = append(experimentalFields, \"%s\")\n", field.JSONPath)
			fmt.Fprintf(w, "\t}\n\n")
		} else {
			// Non-pointer struct - skip, nested experimental fields will be checked separately
			fmt.Fprintf(w, "\t// Skipping %s - non-pointer struct, nested fields checked separately\n\n", field.JSONPath)
		}
	} else {
		// Has slices - generate nested loops
		writeNestedLoopCheck(w, field, sliceIndices)
	}
}

func writeNestedLoopCheck(w io.Writer, field ExperimentalField, sliceIndices []int) {
	indent := "\t"
	loopVars := make([]string, len(sliceIndices))
	hasNilGuard := false

	// Generate nested loops
	for depth, sliceIdx := range sliceIndices {
		loopVar := fmt.Sprintf("item%d", depth)
		loopVars[depth] = loopVar

		// Build path to this slice - use loop variable from previous iteration if exists
		var sliceAccessor string
		if depth == 0 {
			// First loop - access from obj.Spec
			// Add nil guards for parent pointers before the first loop
			nilGuards := buildNilGuards("obj.Spec", field.Path[:sliceIdx])
			if len(nilGuards) > 0 {
				fmt.Fprintf(w, "%sif %s {\n", indent, strings.Join(nilGuards, " && "))
				indent += "\t"
				hasNilGuard = true
			}
			sliceAccessor = buildAccessor("obj.Spec", field.Path[:sliceIdx+1])
		} else {
			// Nested loop - access from previous loop variable
			prevLoopVar := loopVars[depth-1]
			prevSliceIdx := sliceIndices[depth-1]
			// Build path from previous loop var to this slice
			sliceAccessor = prevLoopVar
			for i := prevSliceIdx + 1; i <= sliceIdx; i++ {
				sliceAccessor += "." + field.Path[i].FieldName
			}
		}

		fmt.Fprintf(w, "%sfor _, %s := range %s {\n", indent, loopVar, sliceAccessor)
		indent += "\t"
	}

	// Build accessor for the experimental field from the innermost loop variable
	lastSliceIdx := sliceIndices[len(sliceIndices)-1]
	lastLoopVar := loopVars[len(loopVars)-1]

	var accessor string
	if lastSliceIdx == len(field.Path)-1 {
		// The slice element itself is the experimental field
		accessor = lastLoopVar
	} else {
		// Build path from last loop variable to the experimental field
		accessor = lastLoopVar
		for i := lastSliceIdx + 1; i < len(field.Path); i++ {
			accessor += "." + field.Path[i].FieldName
		}
	}

	// Generate the check
	lastSeg := field.Path[len(field.Path)-1]
	if lastSeg.Type == TypePointer {
		fmt.Fprintf(w, "%sif %s != nil {\n", indent, accessor)
		fmt.Fprintf(w, "%s\texperimentalFields = append(experimentalFields, \"%s\")\n", indent, field.JSONPath)
		fmt.Fprintf(w, "%s\tbreak // Only report once\n", indent)
		fmt.Fprintf(w, "%s}\n", indent)
	} else if !lastSeg.IsStruct {
		fmt.Fprintf(w, "%sif %s != \"\" { // Non-zero check\n", indent, accessor)
		fmt.Fprintf(w, "%s\texperimentalFields = append(experimentalFields, \"%s\")\n", indent, field.JSONPath)
		fmt.Fprintf(w, "%s\tbreak // Only report once\n", indent)
		fmt.Fprintf(w, "%s}\n", indent)
	} else {
		// Non-pointer struct in slice - skip, nested fields checked separately
		fmt.Fprintf(w, "%s// Skipping %s - non-pointer struct, nested fields checked separately\n", indent, field.JSONPath)
	}

	// Close all loops
	for range sliceIndices {
		indent = indent[:len(indent)-1]
		fmt.Fprintf(w, "%s}\n", indent)
	}

	// Close nil guard if we added one
	if hasNilGuard {
		indent = indent[:len(indent)-1]
		fmt.Fprintf(w, "%s}\n", indent)
	}
	fmt.Fprintf(w, "\n")
}

func buildAccessor(base string, path []PathSegment) string {
	accessor := base
	for _, seg := range path {
		accessor += "." + seg.FieldName
	}
	return accessor
}

func buildNilGuards(base string, path []PathSegment) []string {
	var guards []string
	accessor := base
	for _, seg := range path {
		accessor += "." + seg.FieldName
		if seg.Type == TypePointer {
			guards = append(guards, accessor+" != nil")
		}
	}
	return guards
}
