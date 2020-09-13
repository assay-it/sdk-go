//
// Copyright (C) 2018 - 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

package send

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/assay-it/sdk-go/assay"
	"github.com/assay-it/sdk-go/http"
)

//-------------------------------------------------------------------
//
// core arrows
//
//-------------------------------------------------------------------

/*
URL defines a mandatory parameters to the request such as
HTTP method and destination URL, use Params arrow if you
need to supply URL query params.
*/
func URL(method, uri string, args ...interface{}) http.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		var addr *url.URL
		if addr, cat.Fail = mkURL(uri, args...); cat.Fail != nil {
			return cat
		}

		if cat.HTTP == nil {
			cat.HTTP = &assay.IOCatHTTP{}
		}

		switch addr.Scheme {
		case "http", "https":
			cat.HTTP.Send = &assay.UpStreamHTTP{
				Method:  method,
				URL:     addr,
				Header:  make(map[string]*string),
				Payload: bytes.NewBuffer(nil),
			}
		default:
			cat.Fail = errors.New("Not supported")
			// io.Fail = xxxx.ProtocolNotSupported(io.URL.String())
		}
		return cat
	}
}

func mkURL(uri string, args ...interface{}) (*url.URL, error) {
	opts := []interface{}{}
	for _, x := range args {
		switch v := x.(type) {
		case *url.URL:
			v.Path = strings.TrimSuffix(v.Path, "/")
			opts = append(opts, v.String())
		default:
			val := reflect.ValueOf(x)
			if val.Kind() == reflect.Ptr {
				opts = append(opts, url.PathEscape(fmt.Sprintf("%v", val.Elem())))
			} else {
				opts = append(opts, url.PathEscape(fmt.Sprintf("%v", val)))
			}
		}
	}

	return url.Parse(fmt.Sprintf(uri, opts...))
}

/*

HtHeader is tagged string, represents HTTP Header
*/
type HtHeader struct{ string }

/*

Header defines HTTP headers to the request, use combinator
to define multiple header values.

  http.HTTP(
		ø.Header("Accept").Is(...),
		ø.Header("Content-Type").Is(...),
	)
*/
func Header(header string) HtHeader {
	return HtHeader{header}
}

// Is sets a literval value of HTTP header
func (header HtHeader) Is(value string) http.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		cat.HTTP.Send.Header[header.string] = &value
		return cat
	}
}

// Val sets a value of HTTP header from variable
func (header HtHeader) Val(value *string) http.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		cat.HTTP.Send.Header[header.string] = value
		return cat
	}
}
