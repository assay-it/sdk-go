//
// Copyright (C) 2018 - 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

package assay

import (
	"bytes"
	"net/http"
	"net/url"
)

/*

IOCatHTTP defines the category of HTTP I/O
*/
type IOCatHTTP struct {
	Send *HTTPSender
	Recv *HTTPRecver
}

type HTTPSender struct {
	Method  string
	URL     *url.URL
	Header  map[string]*string
	Payload *bytes.Buffer
}

type HTTPRecver struct {
	*http.Response
	Body interface{}
}
