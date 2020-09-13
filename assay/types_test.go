//
// Copyright (C) 2018 - 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

package assay_test

import (
	"errors"
	"testing"

	"github.com/assay-it/sdk-go/assay"
)

func a() assay.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		return cat
	}
}

func f() assay.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		cat.Fail = errors.New("fail")
		return cat
	}
}

func TestJoinAB(t *testing.T) {
	f := assay.Join(a(), a())

	if f(assay.IO()).Fail != nil {
		t.Error("a . a is failed")
	}
}

func TestJoinAF(t *testing.T) {
	f := assay.Join(a(), f())

	if f(assay.IO()).Fail == nil {
		t.Error("a . f is failed")
	}
}

func TestJoinFA(t *testing.T) {
	f := assay.Join(f(), a())

	if f(assay.IO()).Fail == nil {
		t.Error("f . a is failed")
	}
}
