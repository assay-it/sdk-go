//
// Copyright (C) 2018 - 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

package cats_test

import (
	"errors"
	"testing"

	"github.com/assay-it/sdk-go/assay"
	ç "github.com/assay-it/sdk-go/cats"
)

func identity() assay.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		return cat
	}
}

func TestFMap(t *testing.T) {
	var s string

	f := assay.Join(
		identity(),
		ç.FMap(func() error {
			s = "value"
			return nil
		}),
	)
	c := assay.IO()

	if c = f(c); c.Fail != nil || s != "value" {
		t.Error("unable to FMap")
	}
}

func TestFMapError(t *testing.T) {
	f := assay.Join(
		identity(),
		ç.FMap(func() error {
			return errors.New("fail")
		}),
	)
	c := assay.IO()

	if c = f(c); c.Fail == nil {
		t.Error("unable to handle error at FMap")
	}
}

func TestFlatMap(t *testing.T) {
	seq := ""

	f := ç.FMap(func() error {
		seq = seq + "a"
		return nil
	})

	g := assay.Join(
		f,
		ç.FlatMap(func() assay.Arrow { return f }),
	)
	c := assay.IO()

	if c = g(c); c.Fail != nil || seq != "aa" {
		t.Error("unable to FlatMap")
	}
}
