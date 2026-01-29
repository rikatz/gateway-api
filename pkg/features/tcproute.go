/*
Copyright 2026 The Kubernetes Authors.

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

package features

// -----------------------------------------------------------------------------
// Features - TCPRoute Conformance (Core)
// -----------------------------------------------------------------------------

const (
	// This option indicates support for TCPRoute
	SupportTCPRoute FeatureName = "TCPRoute"
)

// TCPRouteFeature contains metadata for the TCPRoute feature.
var TCPRouteFeature = Feature{
	Name:        SupportTCPRoute,
	Channel:     FeatureChannelExperimental,
	Description: "Implements the capability for TCPRoute, providing TCP-based request routing to backend services",
	GEPNumber:   0,
}

// TCPRouteFeatures includes all SupportedFeatures needed to be conformant with
// the TCPRoute resource.
var TCPRouteFeatures = map[FeatureName]Feature{
	SupportTCPRoute: TCPRouteFeature,
}

// -----------------------------------------------------------------------------
// Features - TCPRoute Conformance (Extended)
// -----------------------------------------------------------------------------

const (
	// This option indicates support for the name field in the TCPRouteRule (extended conformance)
	SupportTCPRouteNamedRouteRule FeatureName = "TCPRouteNamedRouteRule"
)

// TCPRouteNamedRouteRule contains metadata for the SupportTCPRouteNamedRouteRule feature.
var TCPRouteNamedRouteRule = Feature{
	Name:        SupportTCPRouteNamedRouteRule,
	Channel:     FeatureChannelExperimental,
	Description: "Implements support for the name field in TCPRouteRule, allowing individual route rules to be referenced by name from other resources",
	GEPNumber:   995,
}
