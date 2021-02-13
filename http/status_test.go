//
// Copyright (C) 2018 - 2021 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

package http_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/assay-it/sdk-go/http"
)

func TestStatusCodeNew(t *testing.T) {
	code := http.NewStatusCode(200)

	if code.Value() != 200 {
		t.Error("NewStatusCode fails")
	}

	if code.Error() != "HTTP 200 OK" {
		t.Error("NewStatusCode fails")
	}
}

func TestStatusCodeAsError(t *testing.T) {
	var err error = http.NewStatusCode(200)

	if err.Error() != "HTTP 200 OK" {
		t.Error("StatusCode is not an error")
	}

	if !errors.Is(err, http.StatusOK) {
		t.Error("StatusCode unable to match to StatusOK")
	}

	if errors.Is(err, http.StatusCreated) {
		t.Error("StatusCode matches to StatusCreated")
	}

	if errors.Is(err, fmt.Errorf("error")) {
		t.Error("StatusCode matches to any error")
	}
}

func TestStatusCodeRequired(t *testing.T) {
	var err error = http.NewStatusCode(200, 201)

	if err.Error() != "HTTP Status `200 OK`, required `201 Created`." {
		t.Error("StatusCode invalid error")
	}
}
