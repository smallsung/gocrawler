package gocrawler

import (
	"net/http"
)

//const DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

type HttpMethod string

const (
	HttpMethodGet  HttpMethod = http.MethodGet
	HttpMethodPost HttpMethod = http.MethodPost
	//HttpMethodHead    HttpMethod = http.MethodHead
	//HttpMethodPut     HttpMethod = http.MethodPut
	//HttpMethodPatch   HttpMethod = http.MethodPatch
	//HttpMethodDelete  HttpMethod = http.MethodDelete
	//HttpMethodConnect HttpMethod = http.MethodConnect
	//HttpMethodOptions HttpMethod = http.MethodOptions
	//HttpMethodTrace   HttpMethod = http.MethodTrace
)

type HttpStatus int

const (
//StatusContinue           HttpStatus = http.StatusContinue           // RFC 7231, 6.2.1
//StatusSwitchingProtocols HttpStatus = http.StatusSwitchingProtocols // RFC 7231, 6.2.2
//StatusProcessing         HttpStatus = http.StatusProcessing         // RFC 2518, 10.1
//StatusEarlyHints         HttpStatus = http.StatusEarlyHints         // RFC 8297
//
//StatusOK                   HttpStatus = http.StatusOK                   // RFC 7231, 6.3.1
//StatusCreated              HttpStatus = http.StatusCreated              // RFC 7231, 6.3.2
//StatusAccepted             HttpStatus = http.StatusAccepted             // RFC 7231, 6.3.3
//StatusNonAuthoritativeInfo HttpStatus = http.StatusNonAuthoritativeInfo // RFC 7231, 6.3.4
//StatusNoContent            HttpStatus = http.StatusNoContent            // RFC 7231, 6.3.5
//StatusResetContent         HttpStatus = http.StatusResetContent         // RFC 7231, 6.3.6
//StatusPartialContent       HttpStatus = http.StatusPartialContent       // RFC 7233, 4.1
//StatusMultiStatus          HttpStatus = http.StatusMultiStatus          // RFC 4918, 11.1
//StatusAlreadyReported      HttpStatus = http.StatusAlreadyReported      // RFC 5842, 7.1
//StatusIMUsed               HttpStatus = http.StatusIMUsed               // RFC 3229, 10.4.1
//
//StatusMultipleChoices   HttpStatus = http.StatusMultipleChoices   // RFC 7231, 6.4.1
//StatusMovedPermanently  HttpStatus = http.StatusMovedPermanently  // RFC 7231, 6.4.2
//StatusFound             HttpStatus = http.StatusFound             // RFC 7231, 6.4.3
//StatusSeeOther          HttpStatus = http.StatusSeeOther          // RFC 7231, 6.4.4
//StatusNotModified       HttpStatus = http.StatusNotModified       // RFC 7232, 4.1
//StatusUseProxy          HttpStatus = http.StatusUseProxy          // RFC 7231, 6.4.5
//_                       HttpStatus = 306                          // RFC 7231, 6.4.6 (Unused)
//StatusTemporaryRedirect HttpStatus = http.StatusTemporaryRedirect // RFC 7231, 6.4.7
//StatusPermanentRedirect HttpStatus = http.StatusPermanentRedirect // RFC 7538, 3
//
//StatusBadRequest                   HttpStatus = http.StatusBadRequest                   // RFC 7231, 6.5.1
//StatusUnauthorized                 HttpStatus = http.StatusUnauthorized                 // RFC 7235, 3.1
//StatusPaymentRequired              HttpStatus = http.StatusPaymentRequired              // RFC 7231, 6.5.2
//StatusForbidden                    HttpStatus = http.StatusForbidden                    // RFC 7231, 6.5.3
//StatusNotFound                     HttpStatus = http.StatusNotFound                     // RFC 7231, 6.5.4
//StatusMethodNotAllowed             HttpStatus = http.StatusMethodNotAllowed             // RFC 7231, 6.5.5
//StatusNotAcceptable                HttpStatus = http.StatusNotAcceptable                // RFC 7231, 6.5.6
//StatusProxyAuthRequired            HttpStatus = http.StatusProxyAuthRequired            // RFC 7235, 3.2
//StatusRequestTimeout               HttpStatus = http.StatusRequestTimeout               // RFC 7231, 6.5.7
//StatusConflict                     HttpStatus = http.StatusConflict                     // RFC 7231, 6.5.8
//StatusGone                         HttpStatus = http.StatusGone                         // RFC 7231, 6.5.9
//StatusLengthRequired               HttpStatus = http.StatusLengthRequired               // RFC 7231, 6.5.10
//StatusPreconditionFailed           HttpStatus = http.StatusPreconditionFailed           // RFC 7232, 4.2
//StatusRequestEntityTooLarge        HttpStatus = http.StatusRequestEntityTooLarge        // RFC 7231, 6.5.11
//StatusRequestURITooLong            HttpStatus = http.StatusRequestURITooLong            // RFC 7231, 6.5.12
//StatusUnsupportedMediaType         HttpStatus = http.StatusUnsupportedMediaType         // RFC 7231, 6.5.13
//StatusRequestedRangeNotSatisfiable HttpStatus = http.StatusRequestedRangeNotSatisfiable // RFC 7233, 4.4
//StatusExpectationFailed            HttpStatus = http.StatusExpectationFailed            // RFC 7231, 6.5.14
//StatusTeapot                       HttpStatus = http.StatusTeapot                       // RFC 7168, 2.3.3
//StatusMisdirectedRequest           HttpStatus = http.StatusMisdirectedRequest           // RFC 7540, 9.1.2
//StatusUnprocessableEntity          HttpStatus = http.StatusUnprocessableEntity          // RFC 4918, 11.2
//StatusLocked                       HttpStatus = http.StatusLocked                       // RFC 4918, 11.3
//StatusFailedDependency             HttpStatus = http.StatusFailedDependency             // RFC 4918, 11.4
//StatusTooEarly                     HttpStatus = http.StatusTooEarly                     // RFC 8470, 5.2.
//StatusUpgradeRequired              HttpStatus = http.StatusUpgradeRequired              // RFC 7231, 6.5.15
//StatusPreconditionRequired         HttpStatus = http.StatusPreconditionRequired         // RFC 6585, 3
//StatusTooManyRequests              HttpStatus = http.StatusTooManyRequests              // RFC 6585, 4
//StatusRequestHeaderFieldsTooLarge  HttpStatus = http.StatusRequestHeaderFieldsTooLarge  // RFC 6585, 5
//StatusUnavailableForLegalReasons   HttpStatus = http.StatusUnavailableForLegalReasons   // RFC 7725, 3
//
//StatusInternalServerError           HttpStatus = http.StatusInternalServerError           // RFC 7231, 6.6.1
//StatusNotImplemented                HttpStatus = http.StatusNotImplemented                // RFC 7231, 6.6.2
//StatusBadGateway                    HttpStatus = http.StatusBadGateway                    // RFC 7231, 6.6.3
//StatusServiceUnavailable            HttpStatus = http.StatusServiceUnavailable            // RFC 7231, 6.6.4
//StatusGatewayTimeout                HttpStatus = http.StatusGatewayTimeout                // RFC 7231, 6.6.5
//StatusHTTPVersionNotSupported       HttpStatus = http.StatusHTTPVersionNotSupported       // RFC 7231, 6.6.6
//StatusVariantAlsoNegotiates         HttpStatus = http.StatusVariantAlsoNegotiates         // RFC 2295, 8.1
//StatusInsufficientStorage           HttpStatus = http.StatusInsufficientStorage           // RFC 4918, 11.5
//StatusLoopDetected                  HttpStatus = http.StatusLoopDetected                  // RFC 5842, 7.2
//StatusNotExtended                   HttpStatus = http.StatusNotExtended                   // RFC 2774, 7
//StatusNetworkAuthenticationRequired HttpStatus = http.StatusNetworkAuthenticationRequired // RFC 6585, 6
)
