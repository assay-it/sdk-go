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

func identity() assay.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		return cat
	}
}

func fail() assay.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		cat.Fail = errors.New("fail")
		return cat
	}
}

func TestJoin(t *testing.T) {
	f := assay.Join(identity(), identity())

	if f(assay.IO()).Fail != nil {
		t.Error("join is failed")
	}
}

func TestJoinAF(t *testing.T) {
	f := assay.Join(identity(), fail())

	if f(assay.IO()).Fail == nil {
		t.Error("join with fail is failed")
	}
}

func TestJoinFA(t *testing.T) {
	f := assay.Join(fail(), identity())

	if f(assay.IO()).Fail == nil {
		t.Error("join with fail is failed")
	}
}

func TestThen(t *testing.T) {
	f := identity().Then(identity())

	if f(assay.IO()).Fail != nil {
		t.Error("join is failed")
	}
}

func TestThenAF(t *testing.T) {
	f := identity().Then(fail())

	if f(assay.IO()).Fail == nil {
		t.Error("join with fail is failed")
	}
}

func TestThenFA(t *testing.T) {
	f := fail().Then(identity())

	if f(assay.IO()).Fail == nil {
		t.Error("join with fail is failed")
	}
}
