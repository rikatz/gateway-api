/*
Copyright 2024 The Kubernetes Authors.

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

import "k8s.io/apimachinery/pkg/util/sets"

// -----------------------------------------------------------------------------
// Features - TLSRoute Conformance (Core)
// -----------------------------------------------------------------------------

const (
	// This option indicates support for TLSRoute
	SupportTLSRoute FeatureName = "TLSRoute"

	// This option indicates support for TLSRoute mode Terminate (extended conformance)
	SupportTLSRouteModeTerminate FeatureName = "TLSRouteModeTerminate"
)

var (
	// TLSRouteFeature contains metadata for the TLSRoute feature.
	TLSRouteFeature = Feature{
		Name:        SupportTLSRoute,
		Channel:     FeatureChannelExperimental,
		Description: "Implements the capability for TLSRoute, allowing a user to route traffic based on SNI attributes",
		GEPNumber:   2643,
	}
	// TLSRouteModeTerminate contains metadata for the TLSRouteModeTerminate feature.
	TLSRouteModeTerminateFeature = Feature{
		Name:        SupportTLSRouteModeTerminate,
		Channel:     FeatureChannelExperimental,
		Description: "Implements the capability for a TLSRoute to be terminated on the Gateway, making the proxy pass the unencrypted packets to the backend without caring about attributes other than a TCP communication.",
		GEPNumber:   2643,
	}
)

// TLSCoreFeatures includes all the supported features for the TLSRoute API at
// a Core level of support.
var TLSRouteCoreFeatures = sets.New(
	TLSRouteFeature,
)

// TLSRouteExtendedFeatures includes all extended features for TLSRoute
// conformance and can be used to opt-in to run all TLSRoute extended features tests.
// This does not include any Core Features.
var TLSRouteExtendedFeatures = sets.New(
	TLSRouteModeTerminateFeature,
)
