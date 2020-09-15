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
	"sort"

	"github.com/assay-it/sdk-go/assay"
	"github.com/google/go-cmp/cmp"
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

// TValue is tagged type, represent matchers
type TValue struct{ actual interface{} }

/*
Value checks if the value equals to defined one.
Supply the pointer to actual value
*/
func Value(val interface{}) TValue {
	return TValue{val}
}

// Is matches a value
func (val TValue) Is(require interface{}) assay.Arrow {
	return func(io *assay.IOCat) *assay.IOCat {
		if diff := cmp.Diff(val.actual, require); diff != "" {
			io.Fail = &assay.Mismatch{
				Diff:    diff,
				Payload: val.actual,
			}
		}
		return io
	}
}

// String matches a literal value
func (val TValue) String(require string) assay.Arrow {
	return val.Is(&require)
}

// Bytes matches a literal value of bytes
func (val TValue) Bytes(require []byte) assay.Arrow {
	return val.Is(&require)
}

// TSeq is tagged type, represents Sequence of elements
type TSeq struct{ assay.Ord }

/*

Seq matches presence of element in the sequence.
*/
func Seq(seq assay.Ord) TSeq {
	return TSeq{seq}
}

/*

Has lookups element using key and matches expected value
*/
func (seq TSeq) Has(key string, expect ...interface{}) assay.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		sort.Sort(seq)
		i := sort.Search(seq.Len(), func(i int) bool { return seq.String(i) >= key })
		if i < seq.Len() && seq.String(i) == key {
			if len(expect) > 0 {
				if diff := cmp.Diff(seq.Value(i), expect[0]); diff != "" {
					cat.Fail = &assay.Mismatch{
						Diff:    diff,
						Payload: seq.Value(i),
					}
				}
			}
			return cat
		}
		cat.Fail = &assay.Undefined{Type: key}
		return cat
	}
}
