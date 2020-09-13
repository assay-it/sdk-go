//
// Copyright (C) 2018 - 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

package recv_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/assay-it/sdk-go/assay"
	µ "github.com/assay-it/sdk-go/http"
	ƒ "github.com/assay-it/sdk-go/http/recv"
	ø "github.com/assay-it/sdk-go/http/send"
)

func TestCodeOk(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.Join(
		ø.GET(ts.URL+"/json"),
		ø.AcceptJSON(),
		ƒ.Code(µ.StatusCodeOK),
	)
	cat := assay.IO(µ.Default())

	if cat = req(cat); cat.Fail != nil {
		t.Error("fail to handle request")
	}
}

func TestCodeNoMatch(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.Join(
		ø.GET(ts.URL+"/other"),
		ø.AcceptJSON(),
		ƒ.Code(µ.StatusCodeOK),
	)
	cat := assay.IO(µ.Default())

	switch req(cat).Fail.(type) {
	case *µ.StatusBadRequest:
		return
	default:
		t.Error("fail to detect code mismatch")
	}
}

func TestHeaderOk(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.Join(
		ø.GET(ts.URL+"/json"),
		ø.AcceptJSON(),
		ƒ.Code(µ.StatusCodeOK),
		ƒ.Header("content-type").Is("application/json"),
	)
	cat := assay.IO(µ.Default())

	if cat = req(cat); cat.Fail != nil {
		t.Error("fail to match header value")
	}
}

func TestHeaderAny(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.Join(
		ø.GET(ts.URL+"/json"),
		ø.AcceptJSON(),
		ƒ.Code(µ.StatusCodeOK),
		ƒ.Header("content-type").Any(),
	)
	cat := assay.IO(µ.Default())

	if cat = req(cat); cat.Fail != nil {
		t.Error("fail to match header value")
	}
}

func TestHeaderVal(t *testing.T) {
	ts := mock()
	defer ts.Close()

	var content string
	req := µ.Join(
		ø.GET(ts.URL+"/json"),
		ø.AcceptJSON(),
		ƒ.Code(µ.StatusCodeOK),
		ƒ.Header("content-type").String(&content),
	)
	cat := assay.IO(µ.Default())

	if cat = req(cat); cat.Fail != nil || content != "application/json" {
		t.Error("fail to match header value")
	}
}

func TestHeaderMismatch(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.Join(
		ø.GET(ts.URL+"/json"),
		ø.AcceptJSON(),
		ƒ.Code(µ.StatusCodeOK),
		ƒ.Header("content-type").Is("foo/bar"),
	)
	cat := assay.IO(µ.Default())

	if cat = req(cat); cat.Fail == nil {
		t.Error("fail to detect header mismatch")
	}
}

//
func mock() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/json":
				w.Header().Add("Content-Type", "application/json")
				w.Write([]byte(`{"site": "example.com"}`))
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		}),
	)
}
