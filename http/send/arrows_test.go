//
// Copyright (C) 2018 - 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

package send_test

import (
	"net/url"
	"testing"

	"github.com/assay-it/sdk-go/assay"
	"github.com/assay-it/sdk-go/http"
	ø "github.com/assay-it/sdk-go/http/send"
)

func TestSchemaHTTP(t *testing.T) {
	req := ø.URL("GET", "http://example.com")
	cat := assay.IO(http.Default())

	if cat = req(cat); cat.Fail != nil {
		t.Error("failed to support http schema")
	}
}

func TestSchemaHTTPS(t *testing.T) {
	req := ø.URL("GET", "https://example.com")
	cat := assay.IO(http.Default())

	if cat = req(cat); cat.Fail != nil {
		t.Error("failed to support https schema")
	}
}

func TestSchemaUnsupported(t *testing.T) {
	req := ø.URL("GET", "other://example.com")
	cat := assay.IO(http.Default())

	if cat = req(cat); cat.Fail == nil {
		t.Error("failed to reject unsupported schema")
	}
}

func TestURL(t *testing.T) {
	req := ø.URL("GET", "https://example.com/%s/%v", "a", 1)
	cat := assay.IO(http.Default())

	if cat = req(cat); cat.HTTP.Send.URL.String() != "https://example.com/a/1" {
		t.Error("unable to assign params to url")
	}
}

func TestURLByRef(t *testing.T) {
	a := "a"
	b := 1
	req := ø.URL("GET", "https://example.com/%s/%v", &a, &b)
	cat := assay.IO(http.Default())

	if cat = req(cat); cat.HTTP.Send.URL.String() != "https://example.com/a/1" {
		t.Error("unable to assign params to url")
	}
}

func TestURLEscape(t *testing.T) {
	a := "a b"
	b := 1
	req := ø.URL("GET", "https://example.com/%s/%v", &a, &b)
	cat := assay.IO(http.Default())

	if cat = req(cat); cat.HTTP.Send.URL.String() != "https://example.com/a%20b/1" {
		t.Error("unable to assign params to url")
	}
}

func TestURLType(t *testing.T) {
	a := "a b"
	b := 1
	p, _ := url.Parse("https://example.com")
	req := ø.URL("GET", "%s/%s/%v", p, &a, &b)
	cat := assay.IO(http.Default())

	if cat = req(cat); cat.HTTP.Send.URL.String() != "https://example.com/a%20b/1" {
		t.Error("unable to assign params to url")
	}
}

func TestHeaderByLit(t *testing.T) {
	req := http.Join(
		ø.URL("GET", "http://example.com"),
		ø.Header("Accept").Is("text/plain"),
	)
	cat := assay.IO(http.Default())

	if cat = req(cat); *cat.HTTP.Send.Header["accept"] != "text/plain" {
		t.Error("unable to set header")
	}
}

func TestHeaderByVal(t *testing.T) {
	val := "text/plain"

	req := http.Join(
		ø.URL("GET", "http://example.com"),
		ø.Header("Accept").Val(&val),
	)
	cat := assay.IO(http.Default())

	if cat = req(cat); *cat.HTTP.Send.Header["accept"] != "text/plain" {
		t.Error("unable to set header")
	}
}

func TestParams(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
		Host string `json:"host,omitempty"`
	}

	req := http.Join(
		ø.URL("GET", "https://example.com"),
		ø.Params(Site{"host", "site"}),
	)
	cat := assay.IO(http.Default())

	if cat = req(cat); cat.HTTP.Send.URL.String() != "https://example.com?host=site&site=host" {
		t.Error("failed to pass query params")
	}
}

func TestParamsInvalidFormat(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
		Host int    `json:"host,omitempty"`
	}

	req := http.Join(
		ø.URL("GET", "https://example.com"),
		ø.Params(Site{"host", 100}),
	)
	cat := assay.IO(http.Default())

	if cat = req(cat); cat.Fail == nil {
		t.Error("failed to reject invalid query params")
	}
}

func TestSendJSON(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
		Host string `json:"host,omitempty"`
	}

	req := http.Join(
		ø.URL("GET", "https://example.com"),
		ø.Header("Content-Type").Is("application/json"),
		ø.Send(Site{"host", "site"}),
	)
	cat := assay.IO(http.Default())

	if cat = req(cat); cat.HTTP.Send.Payload.String() != "{\"site\":\"host\",\"host\":\"site\"}" {
		t.Error("failed to encode JSON")
	}
}

func TestSendForm(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
		Host string `json:"host,omitempty"`
	}

	req := http.Join(
		ø.URL("GET", "https://example.com"),
		ø.Header("Content-Type").Is("application/x-www-form-urlencoded"),
		ø.Send(Site{"host", "site"}),
	)
	cat := assay.IO(http.Default())

	if cat = req(cat); cat.HTTP.Send.Payload.String() != "host=site&site=host" {
		t.Error("failed to encode forms")
	}
}

func TestSendUnknown(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
		Host string `json:"host,omitempty"`
	}

	req := http.Join(
		ø.URL("GET", "https://example.com"),
		ø.Send(Site{"host", "site"}),
	)
	cat := assay.IO(http.Default())

	if cat = req(cat); cat.Fail == nil {
		t.Error("failed to complain about missing Content-Type")
	}
}

func TestSendNotSupported(t *testing.T) {
	type Site struct {
		Site string `json:"site"`
		Host string `json:"host,omitempty"`
	}

	req := http.Join(
		ø.URL("GET", "https://example.com"),
		ø.Header("Content-Type").Is("foo/bar"),
		ø.Send(Site{"host", "site"}),
	)
	cat := assay.IO(http.Default())

	if cat = req(cat); cat.Fail == nil {
		t.Error("failed to complain about unsupported format")
	}
}
