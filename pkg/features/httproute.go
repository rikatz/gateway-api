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
// Features - HTTPRoute Conformance (Core)
// -----------------------------------------------------------------------------

const (
	// This option indicates support for HTTPRoute
	SupportHTTPRoute FeatureName = "HTTPRoute"
)

// HTTPRouteFeature contains metadata for the HTTPRoute feature.
var HTTPRouteFeature = Feature{
	Name:        SupportHTTPRoute,
	Channel:     FeatureChannelStandard,
	Description: "Implements the core HTTPRoute resource, providing HTTP-based request routing and traffic management capabilities",
	GEPNumber:   0,
}

// HTTPRouteCoreFeatures includes all SupportedFeatures needed to be conformant with
// the HTTPRoute resource.
var HTTPRouteCoreFeatures = sets.New(
	HTTPRouteFeature,
)

// -----------------------------------------------------------------------------
// Features - HTTPRoute Conformance (Extended)
// -----------------------------------------------------------------------------

const (
	// This option indicates support for Destination Port matching.
	SupportHTTPRouteDestinationPortMatching FeatureName = "HTTPRouteDestinationPortMatching"

	// This option indicates support for HTTPRoute backend request header modification
	SupportHTTPRouteBackendRequestHeaderModification FeatureName = "HTTPRouteBackendRequestHeaderModification"

	// This option indicates support for HTTPRoute query param matching (extended conformance).
	SupportHTTPRouteQueryParamMatching FeatureName = "HTTPRouteQueryParamMatching"

	// This option indicates support for HTTPRoute method matching (extended conformance).
	SupportHTTPRouteMethodMatching FeatureName = "HTTPRouteMethodMatching"

	// This option indicates support for HTTPRoute response header modification (extended conformance).
	SupportHTTPRouteResponseHeaderModification FeatureName = "HTTPRouteResponseHeaderModification"

	// This option indicates support for HTTPRoute port redirect (extended conformance).
	SupportHTTPRoutePortRedirect FeatureName = "HTTPRoutePortRedirect"

	// This option indicates support for HTTPRoute scheme redirect (extended conformance).
	SupportHTTPRouteSchemeRedirect FeatureName = "HTTPRouteSchemeRedirect"

	// This option indicates support for HTTPRoute path redirect (extended conformance).
	SupportHTTPRoutePathRedirect FeatureName = "HTTPRoutePathRedirect"

	// This option indicates support for HTTPRoute host rewrite (extended conformance)
	SupportHTTPRouteHostRewrite FeatureName = "HTTPRouteHostRewrite"

	// This option indicates support for HTTPRoute path rewrite (extended conformance)
	SupportHTTPRoutePathRewrite FeatureName = "HTTPRoutePathRewrite"

	// This option indicates support for HTTPRoute request mirror (extended conformance).
	SupportHTTPRouteRequestMirror FeatureName = "HTTPRouteRequestMirror"

	// This option indicates support for multiple RequestMirror filters within the same HTTPRoute rule (extended conformance).
	SupportHTTPRouteRequestMultipleMirrors FeatureName = "HTTPRouteRequestMultipleMirrors"

	// This option indicates support for HTTPRoute request mirror filter with percentage based mirroring (extended conformance).
	SupportHTTPRouteRequestPercentageMirror FeatureName = "HTTPRouteRequestPercentageMirror"

	// This option indicates support for HTTPRoute request timeouts (extended conformance).
	SupportHTTPRouteRequestTimeout FeatureName = "HTTPRouteRequestTimeout"

	// This option indicates support for HTTPRoute backendRequest timeouts (extended conformance).
	SupportHTTPRouteBackendTimeout FeatureName = "HTTPRouteBackendTimeout"

	// This option indicates support for HTTPRoute parentRef port (extended conformance).
	SupportHTTPRouteParentRefPort FeatureName = "HTTPRouteParentRefPort"

	// This option indicates support for HTTPRoute with a backendref with an appProtocol 'kubernetes.io/h2c' (extended conformance)
	SupportHTTPRouteBackendProtocolH2C FeatureName = "HTTPRouteBackendProtocolH2C"

	// This option indicates support for HTTPRoute with a backendref with an appProtocol 'kubernetes.io/ws' (extended conformance)
	SupportHTTPRouteBackendProtocolWebSocket FeatureName = "HTTPRouteBackendProtocolWebSocket"

	// This option indicates support for the name field in the HTTPRouteRule (extended conformance)
	SupportHTTPRouteNamedRouteRule FeatureName = "HTTPRouteNamedRouteRule"

	// This option indicates support for the cors filter in the HTTPRouteFilter (extended conformance)
	SupportHTTPRouteCORS FeatureName = "HTTPRouteCORS"
	// This option indicates support for HTTPRoute additional redirect status code 303 (extended conformance)
	SupportHTTPRoute303RedirectStatusCode FeatureName = "HTTPRoute303RedirectStatusCode"

	// This option indicates support for HTTPRoute additional redirect status code 303 (extended conformance)
	SupportHTTPRoute307RedirectStatusCode FeatureName = "HTTPRoute307RedirectStatusCode"

	// This option indicates support for HTTPRoute additional redirect status code 303 (extended conformance)
	SupportHTTPRoute308RedirectStatusCode FeatureName = "HTTPRoute308RedirectStatusCode"

	// This option indicates support for HTTPRoute retries (extended conformance)
	SupportHTTPRouteRetries FeatureName = "HTTPRouteRetries"

	// This option indicates support for HTTPRoute retry budgets (extended conformance)
	SupportHTTPRouteRetryBudget FeatureName = "HTTPRouteRetryBudget"
)

var (
	// HTTPRouteDestinationPortMatchingFeature contains metadata for the HTTPRouteDestinationPortMatching feature.
	HTTPRouteDestinationPortMatchingFeature = Feature{
		Name:        SupportHTTPRouteDestinationPortMatching,
		Channel:     FeatureChannelExperimental,
		Description: "Implements the capability for HTTPRoute to match traffic based on the destination port of incoming requests",
		GEPNumber:   957,
	}
	// HTTPRouteBackendRequestHeaderModificationFeature contains metadata for the HTTPRouteBackendRequestHeaderModification feature.
	HTTPRouteBackendRequestHeaderModificationFeature = Feature{
		Name:        SupportHTTPRouteBackendRequestHeaderModification,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to modify request headers before forwarding to specific backends",
		GEPNumber:   1323,
	}
	// HTTPRouteQueryParamMatchingFeature contains metadata for the HTTPRouteQueryParamMatching feature.
	HTTPRouteQueryParamMatchingFeature = Feature{
		Name:        SupportHTTPRouteQueryParamMatching,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to match and route traffic based on URL query parameters",
		GEPNumber:   0,
	}
	// HTTPRouteMethodMatchingFeature contains metadata for the HTTPRouteMethodMatching feature.
	HTTPRouteMethodMatchingFeature = Feature{
		Name:        SupportHTTPRouteMethodMatching,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to match and route traffic based on HTTP methods (GET, POST, PUT, DELETE, etc.)",
		GEPNumber:   0,
	}
	// HTTPRouteResponseHeaderModificationFeature contains metadata for the HTTPRouteResponseHeaderModification feature.
	HTTPRouteResponseHeaderModificationFeature = Feature{
		Name:        SupportHTTPRouteResponseHeaderModification,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to modify response headers before returning them to clients",
		GEPNumber:   1323,
	}
	// HTTPRoutePortRedirectFeature contains metadata for the HTTPRoutePortRedirect feature.
	HTTPRoutePortRedirectFeature = Feature{
		Name:        SupportHTTPRoutePortRedirect,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to redirect requests to a different port",
		GEPNumber:   726,
	}
	// HTTPRouteSchemeRedirectFeature contains metadata for the HTTPRouteSchemeRedirect feature.
	HTTPRouteSchemeRedirectFeature = Feature{
		Name:        SupportHTTPRouteSchemeRedirect,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to redirect requests between HTTP and HTTPS schemes",
		GEPNumber:   726,
	}
	// HTTPRoutePathRedirectFeature contains metadata for the HTTPRoutePathRedirect feature.
	HTTPRoutePathRedirectFeature = Feature{
		Name:        SupportHTTPRoutePathRedirect,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to redirect requests to a different URL path",
		GEPNumber:   726,
	}
	// HTTPRouteHostRewriteFeature contains metadata for the HTTPRouteHostRewrite feature.
	HTTPRouteHostRewriteFeature = Feature{
		Name:        SupportHTTPRouteHostRewrite,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to rewrite the host header before forwarding requests to backends",
		GEPNumber:   726,
	}
	// HTTPRoutePathRewriteFeature contains metadata for the HTTPRoutePathRewrite feature.
	HTTPRoutePathRewriteFeature = Feature{
		Name:        SupportHTTPRoutePathRewrite,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to rewrite the request path before forwarding to backends",
		GEPNumber:   726,
	}
	// HTTPRouteRequestMirrorFeature contains metadata for the HTTPRouteRequestMirror feature.
	HTTPRouteRequestMirrorFeature = Feature{
		Name:        SupportHTTPRouteRequestMirror,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to mirror/duplicate requests to additional backends for testing or monitoring purposes",
		GEPNumber:   0,
	}
	// HTTPRouteRequestMultipleMirrorsFeature contains metadata for the HTTPRouteRequestMultipleMirrors feature.
	HTTPRouteRequestMultipleMirrorsFeature = Feature{
		Name:        SupportHTTPRouteRequestMultipleMirrors,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to mirror requests to multiple backends simultaneously",
		GEPNumber:   0,
	}
	// HTTPRouteRequestPercentageMirrorFeature contains metadata for the HTTPRouteRequestMultipleMirrors feature.
	HTTPRouteRequestPercentageMirrorFeature = Feature{
		Name:        SupportHTTPRouteRequestPercentageMirror,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to mirror a percentage of requests rather than all requests, enabling gradual traffic sampling",
		GEPNumber:   3171,
	}
	// HTTPRouteRequestTimeoutFeature contains metadata for the HTTPRouteRequestTimeout feature.
	HTTPRouteRequestTimeoutFeature = Feature{
		Name:        SupportHTTPRouteRequestTimeout,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to configure timeouts for the entire request lifecycle from client to backend and back",
		GEPNumber:   1742,
	}
	// HTTPRouteBackendTimeoutFeature contains metadata for the HTTPRouteBackendTimeout feature.
	HTTPRouteBackendTimeoutFeature = Feature{
		Name:        SupportHTTPRouteBackendTimeout,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to configure timeouts specifically for backend request processing",
		GEPNumber:   1742,
	}
	// HTTPRouteParentRefPortFeature contains metadata for the HTTPRouteParentRefPort feature.
	HTTPRouteParentRefPortFeature = Feature{
		Name:        SupportHTTPRouteParentRefPort,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to specify a port in parentRef, allowing routes to attach to specific Gateway listener ports",
		GEPNumber:   0,
	}
	// HTTPRouteBackendProtocolH2CFeature contains metadata for the HTTPRouteBackendProtocolH2C feature.
	HTTPRouteBackendProtocolH2CFeature = Feature{
		Name:        SupportHTTPRouteBackendProtocolH2C,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to communicate with backends using HTTP/2 Cleartext (h2c) protocol via appProtocol annotation",
		GEPNumber:   1911,
	}
	// HTTPRouteBackendProtocolWebSocketFeature contains metadata for the HTTPRouteBackendProtocolWebSocket feature.
	HTTPRouteBackendProtocolWebSocketFeature = Feature{
		Name:        SupportHTTPRouteBackendProtocolWebSocket,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to communicate with backends using WebSocket protocol via appProtocol annotation",
		GEPNumber:   1911,
	}
	// HTTPRouteNamedRouteRule contains metadata for the SupportHTTPRouteNamedRouteRule feature.
	HTTPRouteNamedRouteRule = Feature{
		Name:        SupportHTTPRouteNamedRouteRule,
		Channel:     FeatureChannelStandard,
		Description: "Implements support for the name field in HTTPRouteRule, allowing individual route rules to be referenced by name from other resources",
		GEPNumber:   995,
	}
	// HTTPRouteCORS contains metadata for the SupportHTTPRouteCORS feature.
	HTTPRouteCORS = Feature{
		Name:        SupportHTTPRouteCORS,
		Channel:     FeatureChannelExperimental,
		Description: "Implements the capability for HTTPRoute to configure Cross-Origin Resource Sharing (CORS) policies for handling browser cross-origin requests",
		GEPNumber:   0,
	}
	// HTTPRoute303RedirectStatusCodeFeature contains metadata for the HTTPRoute303RedirectStatusCode feature.
	HTTPRoute303RedirectStatusCodeFeature = Feature{
		Name:        SupportHTTPRoute303RedirectStatusCode,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to use HTTP 303 (See Other) status code for redirects",
		GEPNumber:   726,
	}
	// HTTPRoute307RedirectStatusCodeFeature contains metadata for the HTTPRoute307RedirectStatusCode feature.
	HTTPRoute307RedirectStatusCodeFeature = Feature{
		Name:        SupportHTTPRoute307RedirectStatusCode,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to use HTTP 307 (Temporary Redirect) status code for redirects while preserving the request method",
		GEPNumber:   726,
	}
	// HTTPRoute308RedirectStatusCodeFeature contains metadata for the HTTPRoute308RedirectStatusCode feature.
	HTTPRoute308RedirectStatusCodeFeature = Feature{
		Name:        SupportHTTPRoute308RedirectStatusCode,
		Channel:     FeatureChannelStandard,
		Description: "Implements the capability for HTTPRoute to use HTTP 308 (Permanent Redirect) status code for redirects while preserving the request method",
		GEPNumber:   726,
	}
	// HTTPRouteRetriesFeature contains metadata for the HTTPRouteRetries feature.
	HTTPRouteRetriesFeature = Feature{
		Name:        SupportHTTPRouteRetries,
		Channel:     FeatureChannelExperimental,
		Description: "Implements the capability for HTTPRoute to automatically retry failed requests to backends with configurable retry conditions, maximum attempts, and backoff intervals",
		GEPNumber:   1731,
	}
	// HTTPRouteRetryBudgetFeature contains metadata for the HTTPRouteRetryBudget feature.
	HTTPRouteRetryBudgetFeature = Feature{
		Name:        SupportHTTPRouteRetryBudget,
		Channel:     FeatureChannelExperimental,
		Description: "Implements the capability for HTTPRoute to configure retry budgets that limit retry attempts based on available budget to prevent retry storms",
		GEPNumber:   3388,
	}
)

// HTTPRouteExtendedFeatures includes all extended features for HTTPRoute
// conformance and can be used to opt-in to run all HTTPRoute extended features tests.
// This does not include any Core Features.
var HTTPRouteExtendedFeatures = sets.New(
	HTTPRouteDestinationPortMatchingFeature,
	HTTPRouteBackendRequestHeaderModificationFeature,
	HTTPRouteQueryParamMatchingFeature,
	HTTPRouteMethodMatchingFeature,
	HTTPRouteResponseHeaderModificationFeature,
	HTTPRoutePortRedirectFeature,
	HTTPRouteSchemeRedirectFeature,
	HTTPRoutePathRedirectFeature,
	HTTPRouteHostRewriteFeature,
	HTTPRoutePathRewriteFeature,
	HTTPRouteRequestMirrorFeature,
	HTTPRouteRequestMultipleMirrorsFeature,
	HTTPRouteRequestPercentageMirrorFeature,
	HTTPRouteRequestTimeoutFeature,
	HTTPRouteBackendTimeoutFeature,
	HTTPRouteParentRefPortFeature,
	HTTPRouteBackendProtocolH2CFeature,
	HTTPRouteBackendProtocolWebSocketFeature,
	HTTPRouteNamedRouteRule,
	HTTPRouteCORS,
	HTTPRoute303RedirectStatusCodeFeature,
	HTTPRoute307RedirectStatusCodeFeature,
	HTTPRoute308RedirectStatusCodeFeature,
	HTTPRouteRetriesFeature,
	HTTPRouteRetryBudgetFeature,
)
