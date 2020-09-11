package tk

import (
	"fmt"
	"net/url"
	"os"
)

/*

Env returns value of environment variable.
It is a convenient wrapper over os.LookupEnv
*/
func Env(key, defaultValue string) string {
	value, defined := os.LookupEnv(key)
	if !defined {
		return defaultValue
	}

	return value
}

/*

Host deducts a target host from build environments.
It either uses BUILD_ENDPOINT variable to read the value or
casual composition of v${BUILD_ID}.${BUILD_DOMAIN}
*/
func Host(defaultValue string) string {
	if host := hostBuildEnv(); host != "" {
		return host
	}

	if host := hostCasual(); host != "" {
		return host
	}

	return defaultValue
}

func hostBuildEnv() string {
	endpoint := os.Getenv("BUILD_ENDPOINT")
	if endpoint == "" {
		return ""
	}

	uri, err := url.Parse(endpoint)
	if err != nil {
		return ""
	}
	return uri.Host
}

func hostCasual() string {
	build := os.Getenv("BUILD_ID")
	if build == "" {
		return ""
	}

	domain := os.Getenv("BUILD_DOMAIN")
	if domain == "" {
		return ""
	}

	return fmt.Sprintf("v%s.%s", build, domain)
}
