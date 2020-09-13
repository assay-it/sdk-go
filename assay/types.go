//
// Copyright (c) 2018 - 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

package assay

/*

IOCat defines the category for abstract I/O with a side-effects
*/
type IOCat struct {
	Fail       error
	HTTP       *IOCatHTTP
	LogLevel   int
	sideEffect Arrow
}

/*

Unsafe applies a side effect on the category
*/
func (cat *IOCat) Unsafe() *IOCat {
	return cat.sideEffect(cat)
}

/*

Config defines configuration for the IO category
*/
type Config func(*IOCat) *IOCat

/*

Arrow is a morphism applied to IO category.
The library supports various protocols through definitions of morphisms
*/
type Arrow func(*IOCat) *IOCat

/*

Join composes arrows to high-order function
(a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
*/
func Join(arrows ...Arrow) Arrow {
	return func(cat *IOCat) *IOCat {
		for _, f := range arrows {
			if cat = f(cat); cat.Fail != nil {
				return cat
			}
		}
		return cat
	}
}

/*

IO creates the instance of I/O category use Config type to parametrize
the behavior. The returned value is used to evaluate program.
*/
func IO(opts ...Config) *IOCat {
	cat := &IOCat{}
	for _, opt := range opts {
		cat = opt(cat)
	}
	return cat
}
