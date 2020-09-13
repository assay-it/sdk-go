//
// Copyright (C) 2018 - 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

package recv

import (
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
