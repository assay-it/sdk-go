package assay

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
casual composition ${BUILD_ID} and ${CONFIG_ENDPOINT}
*/
func Host(defaultValue string) *url.URL {
	if host := hostBuildEnv(); host != nil {
		return host
	}

	if host := hostCasual(); host != nil {
		return host
	}

	uri, err := url.Parse(defaultValue)
	if err != nil {
		return &url.URL{}
	}

	return uri
}

func hostBuildEnv() *url.URL {
	endpoint := os.Getenv("BUILD_ENDPOINT")
	if endpoint == "" {
		return nil
	}

	uri, err := url.Parse(endpoint)
	if err != nil {
		return nil
	}
	return uri
}

func hostCasual() *url.URL {
	build := os.Getenv("BUILD_ID")
	if build == "" {
		return nil
	}

	format := os.Getenv("CONFIG_ENDPOINT")
	if format == "" {
		return nil
	}

	uri, err := url.Parse(fmt.Sprintf(format, build))
	if err != nil {
		return nil
	}

	return uri
}
