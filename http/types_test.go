//
// Copyright (C) 2018 - 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/assay-it/sdk-go/assay"
	µ "github.com/assay-it/sdk-go/http"
	ƒ "github.com/assay-it/sdk-go/http/recv"
	ø "github.com/assay-it/sdk-go/http/send"
)

func TestJoin(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := µ.Join(
		ø.URL("GET", ts.URL+"/ok"),
		ƒ.Code(µ.StatusOK),
	)

	cat := assay.IO(µ.Default())
	if cat = req(cat); cat.Fail != nil {
		t.Error("http.Join failed")
	}
}

func TestJoinAssay(t *testing.T) {
	ts := mock()
	defer ts.Close()

	req := assay.Join(
		µ.Join(
			ø.URL("GET", ts.URL+"/ok"),
			ƒ.Code(µ.StatusOK),
		),
		µ.Join(
			ø.URL("GET", ts.URL),
			ƒ.Code(µ.StatusBadRequest),
		),
	)

	cat := assay.IO(µ.Default())
	if cat = req(cat); cat.Fail != nil {
		t.Error("http.Join failed")
	}
}

//
func mock() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.URL.Path == "/ok":
				w.WriteHeader(http.StatusOK)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		}),
	)
}
