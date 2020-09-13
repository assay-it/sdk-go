//
// Copyright (C) 2018 - 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

package http

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/assay-it/sdk-go/assay"
)

/*

Arrow is a morphism applied to HTTP
*/
type Arrow func(*assay.IOCat) *assay.IOCat

/*

Join composes HTTP arrows to high-order function
(a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
*/
func Join(arrows ...Arrow) assay.Arrow {
	return func(cat *assay.IOCat) *assay.IOCat {
		cat.HTTP = &assay.IOCatHTTP{}

		for _, f := range arrows {
			if cat = f(cat); cat.Fail != nil {
				return cat
			}
		}

		return cat
	}
}

/*

Stack configures custom HTTP stack for the category.
*/
func Stack(client *http.Client) assay.Config {
	pool := pool{client}
	return assay.SideEffect(pool.Unsafe)
}

/*

Default configures default HTTP stack for the category.
*/
func Default() assay.Config {
	pool := pool{defaultClient()}
	return assay.SideEffect(pool.Unsafe)
}

func defaultClient() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			ReadBufferSize: 128 * 1024,
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

type pool struct{ *http.Client }

func (p pool) Unsafe(cat *assay.IOCat) *assay.IOCat {
	if cat.Fail != nil {
		return cat
	}

	var eg *http.Request
	eg, cat.Fail = http.NewRequest(
		cat.HTTP.Send.Method,
		cat.HTTP.Send.URL.String(),
		cat.HTTP.Send.Payload,
	)
	if cat.Fail != nil {
		return cat
	}

	for head, value := range cat.HTTP.Send.Header {
		eg.Header.Set(head, *value)
	}

	var in *http.Response
	in, cat.Fail = p.Client.Do(eg)
	if cat.Fail != nil {
		return cat
	}

	cat.HTTP.Recv = &assay.HTTPRecver{Response: in}

	logSend(cat.LogLevel, eg)
	logRecv(cat.LogLevel, in)

	return cat
}

func logSend(level int, eg *http.Request) {
	if level >= assay.LogLevelEgress {
		if msg, err := httputil.DumpRequest(eg, level == assay.LogLevelDebug); err == nil {
			log.Printf(">>>>\n%s\n", msg)
		}
	}
}

func logRecv(level int, in *http.Response) {
	if level >= assay.LogLevelIngress {
		if msg, err := httputil.DumpResponse(in, level == assay.LogLevelDebug); err == nil {
			log.Printf("<<<<\n%s\n", msg)
		}
	}
}
