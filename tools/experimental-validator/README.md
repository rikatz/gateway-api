# Experimental Field Validator Generator

A Go code generator that creates validation functions to detect when Gateway API objects use experimental fields.

## Overview

This tool:
1. Scans Gateway API types in `apis/` directories
2. Identifies types marked with `// +kubebuilder:object:root=true`
3. Finds fields marked with `// <gateway:experimental>`
4. Generates validation code and tests

## Generated Code

**Output files:**
- `pkg/experimental/zz_generated.validation.go` - Validation functions
- `pkg/experimental/zz_generated.validation_test.go` - Generated test skeletons

**Main function:**
```go
func HasExperimentalFields(obj any) error
```

Returns an error listing all experimental fields in use, or `nil` if none are used.

## Usage

### Running the Generator

```bash
go run ./tools/experimental-validator
```

This generates/updates files in `pkg/experimental/`.

### Using the Validation Function

```go
package main

import (
	"fmt"
	v1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/pkg/experimental"
	"k8s.io/utils/ptr"
)

func main() {
	route := &v1.HTTPRoute{
		Spec: v1.HTTPRouteSpec{
			Rules: []v1.HTTPRouteRule{
				{
					Retry: &v1.HTTPRouteRetry{
						Attempts: ptr.To(3),
					},
				},
			},
		},
	}

	if err := experimental.HasExperimentalFields(route); err != nil {
		fmt.Printf("⚠️  %v\n", err)
		// Output: ⚠️  experimental fields in use: [spec.rules.retry]
	}
}
```

## Features

### Implemented ✓

- ✅ Uses `any` instead of `interface{}`
- ✅ Detects experimental fields at all nesting levels
- ✅ Handles pointer and non-pointer fields correctly
- ✅ Supports nested slices with proper loop generation
- ✅ Returns clear error messages with JSON paths
- ✅ Only validates Spec (ignores Status and Metadata)
- ✅ Supports both value and pointer object types
- ✅ Generates tests (with manual implementation required)
- ✅ No TODOs in production code
- ✅ No external dependencies beyond go.mod

### Generated Loop Code

For nested slices like `spec.rules.filters.cors`, the generator creates:

```go
for _, item0 := range obj.Spec.Rules {
    for _, item1 := range item0.Filters {
        if item1.CORS != nil {
            experimentalFields = append(experimentalFields, "spec.rules.filters.cors")
            break // Only report once
        }
    }
}
```

## Testing

### Generator Tests

```bash
cd tools/experimental-validator
go test -v
```

Tests include:
- Parsing test APIs
- Detecting experimental fields
- Nested slice handling
- Real API detection (HTTPRoute, Gateway, etc.)

### Validation Tests

```bash
cd pkg/experimental
go test -v
```

Includes manual tests demonstrating:
- HTTPRoute with Retry (experimental)
- HTTPRoute with CORS filter (experimental)
- Gateway with DefaultScope (experimental)
- Objects without experimental fields

## Architecture

### Type Information Tracking

The generator uses `PathSegment` to track full type information:

```go
type PathSegment struct {
    FieldName string    // Go field name
    JSONName  string    // JSON tag name
    Type      FieldType // Direct/Pointer/Slice
    TypeName  string    // Go type name
}
```

This enables correct code generation for:
- Direct fields: `if obj.Spec.Field != "" {`
- Pointer fields: `if obj.Spec.Field != nil {`
- Slice fields: `for _, item := range obj.Spec.Fields {`

### Generator Flow

1. **Parse** - Load API definitions with AST parser
2. **Analyze** - Find root types and experimental fields
3. **Track** - Build complete path with type information
4. **Generate** - Create validation functions and tests

## Files

```
tools/experimental-validator/
├── main.go              # Generator implementation
├── generator_test.go    # Generator tests
├── README.md            # This file
└── testdata/
    └── apis/testv1/     # Test API for validation
        ├── doc.go
        └── types.go     # Contains experimental fields

pkg/experimental/
├── zz_generated.validation.go       # Generated validation (156 lines)
├── zz_generated.validation_test.go  # Generated test skeletons
└── validation_manual_test.go        # Working test examples
```

## Detected Experimental Fields

**HTTPRoute:**
- `spec.useDefaultGateways`
- `spec.rules[].retry`
- `spec.rules[].sessionPersistence`
- `spec.rules[].filters[].cors`
- `spec.rules[].filters[].externalAuth`
- (and more...)

**Gateway:**
- `spec.allowedListeners`
- `spec.tls`
- `spec.defaultScope`

## Integration

To integrate into the build process:

1. Add to `Makefile`:
```makefile
.PHONY: generate-experimental-validator
generate-experimental-validator:
	go run ./tools/experimental-validator
```

2. Update `hack/update-codegen.sh`:
```bash
echo "Generating experimental field validators"
go run ./tools/experimental-validator
```

3. Add verification:
```bash
# In hack/verify-codegen.sh
git diff --exit-code pkg/experimental/zz_generated.validation.go
```

## Development

### Adding New API Versions

The generator automatically scans:
- `apis/v1`
- `apis/v1alpha2`
- `apis/v1alpha3`
- `apis/v1beta1`

New versions are picked up automatically when added to the `apiVersions` slice in `main.go`.

### Debugging

Enable verbose output by modifying the generator to log field discovery:

```go
log.Printf("Found experimental field: %s", field.JSONPath)
```

## License

Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0.
