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

func TestDefined(t *testing.T) {
	type Site struct {
		Site string
		Host string
	}
	var site Site

	f := assay.Join(
		ç.FMap(func() error {
			site = Site{"site", "host"}
			return nil
		}),
		ç.Defined(&site),
		ç.Defined(&site.Site),
	)
	c := assay.IO()

	if c = f(c); c.Fail != nil {
		t.Error("unable to check is value is defined")
	}
}

func TestNotDefined(t *testing.T) {
	type Site struct {
		Site string
		Host string
	}
	var site Site

	f := assay.Join(
		ç.FMap(func() error {
			site = Site{"site", ""}
			return nil
		}),
		ç.Defined(&site),
		ç.Defined(&site.Host),
	)
	c := assay.IO()

	if c = f(c); c.Fail == nil {
		t.Error("unable to catch undefined value")
	}
}

func TestValue(t *testing.T) {
	type Site struct {
		Site string
		Host []byte
	}
	var site Site

	f := assay.Join(
		ç.FMap(func() error {
			site = Site{"site", []byte("abc")}
			return nil
		}),
		ç.Value(&site).Is(&Site{"site", []byte("abc")}),
		ç.Value(&site.Site).String("site"),
		ç.Value(&site.Host).Bytes([]byte("abc")),
	)
	c := assay.IO()

	if c = f(c); c.Fail != nil {
		t.Error("unable to match value")
	}
}

func TestValueNoMatch(t *testing.T) {
	type Site struct {
		Site string
		Host []byte
	}
	var site Site

	f := assay.Join(
		ç.FMap(func() error {
			site = Site{"site", []byte("abc")}
			return nil
		}),
		ç.Value(&site).Is(&Site{"site1", []byte("abc")}),
	)
	c := assay.IO()

	if c = f(c); c.Fail == nil {
		t.Error("unable to detect mismatched value(s)")
	}
}

type E struct{ Site string }

type Seq []E

func (seq Seq) Len() int                { return len(seq) }
func (seq Seq) Swap(i, j int)           { seq[i], seq[j] = seq[j], seq[i] }
func (seq Seq) Less(i, j int) bool      { return seq[i].Site < seq[j].Site }
func (seq Seq) String(i int) string     { return seq[i].Site }
func (seq Seq) Value(i int) interface{} { return seq[i] }

var seqMock Seq = Seq{
	{Site: "q.example.com"},
	{Site: "a.example.com"},
	{Site: "z.example.com"},
	{Site: "w.example.com"},
	{Site: "s.example.com"},
	{Site: "x.example.com"},
	{Site: "e.example.com"},
	{Site: "d.example.com"},
	{Site: "c.example.com"},
}

func TestSeqHas(t *testing.T) {
	var seq Seq
	expectS := E{Site: "s.example.com"}
	expectZ := E{Site: "z.example.com"}

	f := assay.Join(
		ç.FMap(func() error {
			seq = seqMock
			return nil
		}),
		ç.Seq(&seq).Has(expectS.Site),
		ç.Seq(&seq).Has(expectS.Site, expectS),
		ç.Seq(&seq).Has(expectZ.Site, expectZ),
	)
	c := assay.IO()

	if c = f(c); c.Fail != nil {
		t.Error("unable to match seq")
	}
}

func TestSeqHasNotFound(t *testing.T) {
	var seq Seq
	expect0 := E{Site: "0.example.com"}

	f := assay.Join(
		ç.FMap(func() error {
			seq = seqMock
			return nil
		}),
		ç.Seq(&seq).Has(expect0.Site),
	)
	c := assay.IO()

	if c = f(c); c.Fail == nil {
		t.Error("unable to detect missing element")
	}
}

func TestSeqHasNoMatch(t *testing.T) {
	var seq Seq
	expectS := E{Site: "s.example.com"}
	expectZ := E{Site: "z.example.com"}

	f := assay.Join(
		ç.FMap(func() error {
			seq = seqMock
			return nil
		}),
		ç.Seq(&seq).Has(expectS.Site, expectZ),
	)
	c := assay.IO()

	if c = f(c); c.Fail == nil {
		t.Error("unable to detect mismatched element")
	}
}
