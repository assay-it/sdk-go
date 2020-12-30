//
// Copyright (C) 2018 - 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay-it/sdk-go
//

/*

Package oauth2 support OAuth2 protocol https://tools.ietf.org/html/rfc6749,
use this package to build quality assessment pipelines against protected
interface. The package supports following grants:
* Client Credentials Grant

  import (
    "github.com/assay-it/sdk-go/assay"
    "github.com/assay-it/sdk-go/http"
    "github.com/assay-it/sdk-go/oauth2"
    ø "github.com/assay-it/sdk-go/http/send"
  )

  var auth oauth2.Client = oauth2.Client{
    // obtain oauth2 client ID and Secret from your authorization service
    ID:     "xxx",
    Secret: "xxx-xxx",

    // define scopes of the client
    Scopes: []string{"my/scope"},

    //
    Server: oauth2.Server{
      TokenURL: "https://xyz.auth.eu-west-1.amazoncognito.com/oauth2/token",
    },
  }

  // compose Client Credentials flow with API I/O into HoC function
  func MyScenario() assay.Arrow {
    return oauth2.Join(
      auth.ClientCredentials,
      MyRequest,
    )
  }

  func MyRequest(token *oauth2.Token) assay.Arrow {
    return http.Join(
      // ...
      ø.Header("Authorization").Val(&t.AccessToken)
    )
  }

*/
package oauth2
