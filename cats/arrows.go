//
// Copyright (C) 2018 - 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

package cats

import (
	"reflect"

	"github.com/assay-it/sdk-go/assay"
)

/*

FMap applies clojure to category.
The function lifts any computation to the category and make it composable
with the "program".
*/
func FMap(f func() error) assay.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		cat.Fail = f()
		return cat
	}
}

/*

FlatMap applies closure to matched HTTP request.
It returns an arrow, which continue evaluation.
*/
func FlatMap(f func() assay.Arrow) assay.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		if g := f(); g != nil {
			return g(cat)
		}
		return cat
	}
}

/*

Defined checks if the value is defined, use a pointer to the value.
*/
func Defined(value interface{}) assay.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		va := reflect.ValueOf(value)
		if va.Kind() == reflect.Ptr {
			va = va.Elem()
		}

		if !va.IsValid() {
			cat.Fail = &assay.Undefined{Type: va.Type().Name()}
		}

		if va.IsValid() && va.IsZero() {
			cat.Fail = &assay.Undefined{Type: va.Type().Name()}
		}
		return cat
	}
}
