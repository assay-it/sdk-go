//
// Copyright (C) 2018 - 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

package recv

import (
	"fmt"
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

Code is a mandatory statement to match expected HTTP Status Code against
received one. The execution fails with BadMatchCode if service responds
with other value then specified one.
*/
func Code(code ...http.StatusCodeAny) http.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		if cat = cat.Unsafe(); cat.Fail != nil {
			return cat
		}

		status := cat.HTTP.Recv.StatusCode
		if !hasCode(code, status) {
			cat.Fail = http.NewStatusCode(status, code[0])
		}
		return cat
	}
}

func hasCode(s []http.StatusCodeAny, e int) bool {
	for _, a := range s {
		if a.Value() == e {
			return true
		}
	}
	return false
}

// THeader is tagged string, represents HTTP Header
type THeader struct{ string }

/*

Header matches presence of header in the response or match its entire content.
The execution fails with BadMatchHead if the matched value do not meet expectations.
*/
func Header(header string) THeader {
	return THeader{header}
}

// Is matches value of HTTP header, Use wildcard string ("*") to match any header value
func (header THeader) Is(value string) http.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		h := cat.HTTP.Recv.Header.Get(header.string)
		if h == "" {
			cat.Fail = &assay.Mismatch{
				Diff:    fmt.Sprintf("- %s: %s", header.string, value),
				Payload: nil,
			}
			return cat
		}

		if value != "*" && !strings.HasPrefix(h, value) {
			cat.Fail = &assay.Mismatch{
				Diff:    fmt.Sprintf("+ %s: %s\n- %s: %s", header.string, h, header.string, value),
				Payload: map[string]string{header.string: h},
			}
			return cat
		}

		return cat
	}
}

// String matches a header value to closed variable of string type.
func (header THeader) String(value *string) http.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		val := cat.HTTP.Recv.Header.Get(header.string)
		if val == "" {
			cat.Fail = &assay.Mismatch{
				Diff:    fmt.Sprintf("- %s: *", header.string),
				Payload: nil,
			}
		} else {
			*value = val
		}

		return cat
	}
}

// Any matches a header value, syntax sugar of Header(...).Is("*")
func (header THeader) Any() http.Arrow {
	return header.Is("*")
}
