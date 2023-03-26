package mainimpl

import (
	"strings"

	"github.com/unrolled/secure"
)

func SetupSecureMiddleware(httpsDomain string) *secure.Secure {
	secureMiddleware := secure.New(secure.Options{
		// When developing, the AllowedHosts, SSL, and STS options can cause some unwanted effects. Usually testing happens on http, not https, and on localhost, not your production domain... so set this to true for dev environment.
		// If you would like your development environment to mimic production with complete Host blocking, SSL redirects, and STS headers, leave this as false. Default if false.
		IsDevelopment: development,
		// AllowedHosts is a slice of fully qualified domain names that are allowed. Default is an empty slice, which allows any and all host names.
		AllowedHosts: []string{
			// the normal domain
			httpsDomain,
			// the domain but as a wildcard match with *. infront
			`*\.` + strings.Replace(httpsDomain, ".", `\.`, -1),
		},
		// AllowedHostsAreRegex determines, if the provided `AllowedHosts` slice contains valid regular expressions. If this flag is set to true, every request's host will be checked against these expressions. Default is false.
		// for the wildcard matching
		AllowedHostsAreRegex: true,

		// TLS stuff
		// If SSLRedirect is set to true, then only allow https requests. Default is false.
		SSLRedirect: true,
		// SSLHost is the host name that is used to redirect http requests to https. Default is "", which indicates to use the same host.
		SSLHost: httpsDomain,

		// Important for reverse-proxy setups (when nginx or similar does the TLS termination)
		// SSLProxyHeaders is set of header keys with associated values that would indicate a valid https request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		// HostsProxyHeaders is a set of header keys that may hold a proxied hostname value for the request.
		HostsProxyHeaders: []string{"X-Forwarded-Host"},

		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security
		// STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
		STSSeconds: 2592000, // 30 days in seconds
		// If STSPreload is set to true, the `preload` flag will be appended to the Strict-Transport-Security header. Default is false.
		STSPreload: false, // don't submit to googles list service
		// If STSIncludeSubdomains is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false.
		STSIncludeSubdomains: false,

		// See for more https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP
		// helpful: https://report-uri.com/home/generate
		// ContentSecurityPolicy allows the Content-Security-Policy header value to be set with a custom value. Default is "".
		ContentSecurityPolicy: "default-src 'self'; img-src 'self' data:", // enforce no external content

		// If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
		BrowserXssFilter: true,
		// If FrameDeny is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false.
		FrameDeny: true,
		// If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
		//ContentTypeNosniff: true,
	})

	return secureMiddleware
}
