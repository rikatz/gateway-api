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

package experimental

import (
	"strings"
	"testing"

	"k8s.io/utils/ptr"

	v1 "sigs.k8s.io/gateway-api/apis/v1"
)

// These are manual tests that demonstrate the generated validation code works correctly.
// They serve as examples for how to construct test objects with experimental fields.

func TestHTTPRoute_WithRetry(t *testing.T) {
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

	err := v1.HasExperimentalFields(route)
	if err == nil {
		t.Error("Expected error for HTTPRoute with experimental Retry field, got nil")
	}
	if !strings.Contains(err.Error(), "retry") {
		t.Errorf("Expected error to mention 'retry', got: %v", err)
	}
	t.Logf("✓ Correctly detected experimental retry field: %v", err)
}

func TestHTTPRoute_WithCORS(t *testing.T) {
	route := &v1.HTTPRoute{
		Spec: v1.HTTPRouteSpec{
			Rules: []v1.HTTPRouteRule{
				{
					Filters: []v1.HTTPRouteFilter{
						{
							Type: v1.HTTPRouteFilterCORS,
							CORS: &v1.HTTPCORSFilter{
								AllowOrigins: []v1.CORSOrigin{"*"},
							},
						},
					},
				},
			},
		},
	}

	err := v1.HasExperimentalFields(route)
	if err == nil {
		t.Error("Expected error for HTTPRoute with experimental CORS filter, got nil")
	}
	if !strings.Contains(err.Error(), "cors") {
		t.Errorf("Expected error to mention 'cors', got: %v", err)
	}
	t.Logf("✓ Correctly detected experimental CORS field: %v", err)
}

func TestHTTPRoute_NoExperimental(t *testing.T) {
	path := "/"
	route := &v1.HTTPRoute{
		Spec: v1.HTTPRouteSpec{
			Rules: []v1.HTTPRouteRule{
				{
					Matches: []v1.HTTPRouteMatch{
						{
							Path: &v1.HTTPPathMatch{
								Value: &path,
							},
						},
					},
				},
			},
		},
	}

	err := v1.HasExperimentalFields(route)
	if err != nil {
		t.Errorf("Expected no error for HTTPRoute without experimental fields, got: %v", err)
	}
	t.Log("✓ Correctly identified no experimental fields")
}

func TestGateway_WithDefaultScope(t *testing.T) {
	gw := &v1.Gateway{
		Spec: v1.GatewaySpec{
			DefaultScope: "namespace",
		},
	}

	err := v1.HasExperimentalFields(gw)
	if err == nil {
		t.Error("Expected error for Gateway with experimental DefaultScope field, got nil")
	}
	if !strings.Contains(err.Error(), "defaultScope") {
		t.Errorf("Expected error to mention 'defaultScope', got: %v", err)
	}
	t.Logf("✓ Correctly detected experimental defaultScope field: %v", err)
}

func TestGateway_NoExperimental(t *testing.T) {
	gw := &v1.Gateway{
		Spec: v1.GatewaySpec{
			GatewayClassName: "example",
		},
	}

	err := v1.HasExperimentalFields(gw)
	if err != nil {
		t.Errorf("Expected no error for Gateway without experimental fields, got: %v", err)
	}
	t.Log("✓ Correctly identified no experimental fields")
}
