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
	for _, f := range []assay.Arrow{
		assay.Join(identity(), identity()),
		identity().Then(identity()),
	} {
		if f(assay.IO()).Fail != nil {
			t.Error("join is failed")
		}
	}
}

func TestJoinFail(t *testing.T) {
	for _, f := range []assay.Arrow{
		assay.Join(identity(), fail()),
		assay.Join(fail(), identity()),
		identity().Then(fail()),
		fail().Then(identity()),
	} {
		if f(assay.IO()).Fail == nil {
			t.Error("join with fail is failed")
		}
	}
}
