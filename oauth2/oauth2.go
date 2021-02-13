//
// Copyright (C) 2018 - 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

package oauth2

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/assay-it/sdk-go/assay"
	c "github.com/assay-it/sdk-go/cats"
	"github.com/assay-it/sdk-go/http"
	ƒ "github.com/assay-it/sdk-go/http/recv"
	ø "github.com/assay-it/sdk-go/http/send"
)

// Server is "Authorization Server" configuration
type Server struct {
	TokenURL string
}

// Client is OAuth2 Client configuration
type Client struct {
	ID     string
	Secret string
	Scopes []string
	Server Server
}

// Token is oauth2 access token
type Token struct {
	AccessToken string `json:"access_token"`
	Type        string `json:"token_type,omitempty"`
	Expires     int64  `json:"expires_in,omitempty"`
}

type tAccessTokenRequest struct {
	Type  string `json:"grant_type"`
	Scope string `json:"scope"`
}

// ClientCredentials implements OAuth2 Client Credentials Grant
func (client *Client) ClientCredentials(token *Token) assay.Arrow {
	return http.Join(
		ø.POST(client.Server.TokenURL),
		ø.Header("Authorization").Val(client.basicHTTPAuthorization()),
		ø.ContentForm(),
		ø.Send(tAccessTokenRequest{
			Type:  "client_credentials",
			Scope: strings.Join(client.Scopes, " "),
		}),
		ƒ.Code(http.StatusOK),
		ƒ.Recv(token),
	).Then(
		c.FMap(func() error {
			token.AccessToken = token.Type + " " + token.AccessToken
			return nil
		}),
	)
}

func (client *Client) basicHTTPAuthorization() *string {
	pair := fmt.Sprintf("%s:%s", client.ID, client.Secret)
	digest := fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(pair)))
	return &digest
}

/*

Arrow is a morphism over OAuth2 protocol. The arrow is a function takes access
token as input and returns I/O
*/
type Arrow func(*Token) assay.Arrow

/*

Join composes OAuth2 arrows to high-order function, use grant flow as head of
composition
(a ⟼ b, b ⟼ c, c ⟼ d) ⤇ a ⟼ d
*/
func Join(arrows ...Arrow) assay.Arrow {
	var token Token

	seq := []assay.Arrow{}
	for _, f := range arrows {
		seq = append(seq, f(&token))
	}

	return assay.Join(seq...)
}
