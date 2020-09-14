# assay.it SDK for Go

[assay.it](https://assay.it) is the service to **confirm quality** and **eliminate production risks** in Serverless, Microservices and other applications. `assay-it/sdk-go` is the official SDK for the Go programming language.

[![Documentation](https://pkg.go.dev/badge/github.com/assay-it/sdk-go)](https://pkg.go.dev/github.com/assay-it/sdk-go)
[![Build Status](https://github.com/assay-it/sdk-go/workflows/build/badge.svg)](https://github.com/assay-it/sdk-go/actions/)
[![Git Hub](https://img.shields.io/github/last-commit/assay-it/sdk-go.svg)](http://github.com/assay-it/sdk-go)
[![Coverage Status](https://coveralls.io/repos/github/assay-it/sdk-go/badge.svg?branch=master)](https://coveralls.io/github/assay-it/sdk-go?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/assay-it/sdk-go)](https://goreportcard.com/report/github.com/assay-it/sdk-go)
[![Maintainability](https://api.codeclimate.com/v1/badges/9e6bbe00ab093e5465ff/maintainability)](https://codeclimate.com/github/assay-it/sdk-go/maintainability)


Use this SDK to develop **pure functional** and **typesafe** Behavior as a Code to validate applications. [assay.it](https://assay.it) automates the validation process along CI/CD pipelines. Continuous proofs of the quality helps to eliminate defects at earlier phases of loosely coupled topologies such as serverless applications, microservices and other systems that rely on interface syntaxes and its behaviors.


## Getting Started

### Installing

Use `go get` to retrieve the SDK to add it to your `GOPATH` workspace.

```bash
go get github.com/assay-it/sdk-go
```

Use `go get -u` to update SDK with latest version

```bash
go get -u github.com/assay-it/sdk-go
```

### Quick Example

This example shows a minimal Behavior as a Code suite, which pings the website and ensures the response is correct.

```go
package main

import (
  "github.com/assay-it/sdk-go/assay"
  "github.com/assay-it/sdk-go/http"
  ƒ "github.com/assay-it/sdk-go/http/recv"
  ø "github.com/assay-it/sdk-go/http/send"
)

func TestOk() assay.Arrow {
  return http.Join(
    ø.GET("https://assay.it"),
    ƒ.Code(http.StatusCodeOK),
    ƒ.Header("Content-Type").Is("text/html"),
  )
}

func main() {
  cat := assay.IO(
    http.Default(),
    assay.Logging(assay.LogLevelDebug),
  )
  TestOk()(cat)
}
```

## Further Reading

[Minimal example](https://github.com/assay-it/sample.assay.it) - an annotated walk-through and a minimal Behavior as a Code suite. Use it as getting started tutorial.

[Advanced example](https://github.com/assay-it/example.assay.it) - an annotated walk-through about [CI/CD workflows](https://assay.it/2020/07/01/everything-is-continuos/), shows advanced usage of Behavior as a Code paradigm.

[Developer Guide](https://assay.it/doc/core) - The documentation is a general introduction on how to configure and use the SDK.


## How To Contribute

SDK is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

### bugs

If you experience any issues with SDK, please let us know via [GitHub issues](https://github.com/assay-it/sdk-go/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 


## License

[![See LICENSE](https://img.shields.io/github/license/assay-it/sdk-go.svg?style=for-the-badge)](LICENSE)
