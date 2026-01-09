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

package testv1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

// TestResource is a test resource for validating experimental field detection
type TestResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestResourceSpec   `json:"spec,omitempty"`
	Status TestResourceStatus `json:"status,omitempty"`
}

// TestResourceSpec defines the desired state of TestResource
type TestResourceSpec struct {
	// SimpleField is a regular field
	SimpleField string `json:"simpleField,omitempty"`

	// ExperimentalDirectField is a direct experimental field
	// +optional
	// <gateway:experimental>
	ExperimentalDirectField string `json:"experimentalDirectField,omitempty"`

	// ExperimentalPointer is a pointer experimental field
	// +optional
	// <gateway:experimental>
	ExperimentalPointer *ExperimentalConfig `json:"experimentalPointer,omitempty"`

	// Items contains a list of test items
	// +optional
	Items []TestItem `json:"items,omitempty"`
}

// TestItem represents an item in a list
type TestItem struct {
	// Name of the item
	Name string `json:"name,omitempty"`

	// ExperimentalItemField is experimental
	// +optional
	// <gateway:experimental>
	ExperimentalItemField *string `json:"experimentalItemField,omitempty"`

	// NestedItems contains nested items
	// +optional
	NestedItems []NestedItem `json:"nestedItems,omitempty"`
}

// NestedItem represents a nested item
type NestedItem struct {
	// Value of the nested item
	Value string `json:"value,omitempty"`

	// ExperimentalNested is experimental
	// +optional
	// <gateway:experimental>
	ExperimentalNested *NestedConfig `json:"experimentalNested,omitempty"`
}

// ExperimentalConfig is an experimental configuration
type ExperimentalConfig struct {
	// Enabled indicates if the feature is enabled
	Enabled bool `json:"enabled,omitempty"`
}

// NestedConfig is a nested experimental configuration
type NestedConfig struct {
	// Mode specifies the mode
	Mode string `json:"mode,omitempty"`
}

// TestResourceStatus defines the observed state of TestResource
type TestResourceStatus struct {
	// Conditions represent the latest available observations
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}
