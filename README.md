# go-form3

Łukasz Mikołajczak - no prior commercial experience with Go.

However, I try to contribute to open-source Go projects and maintain some personal ones. Besides that I learned Go by studying books and checking the most popular Go projects on GitHub. Therefore, code in this repo is based on the good practices and is influenced by the patterns that I observed in various open source projects and in Go's standard library.

### Test:

```
docker-compose up
```

### Usage:

```go
package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/lmikolajczak/go-form3/form3"
)

func main() {
	// Initialise new client
	f3 := form3.NewClient("http://localhost:8080")

	// Create new account with the given attributes
	account, err := f3.CreateAccount(
		uuid.NewString(),
		&form3.AccountAttributes{
			Country: form3.String("NL"),
			Name:    []string{"L. Mikolajczak"},
		},
	)
	// ...handle the error...

	// Fetch account with the given ID
	account, err = f3.FetchAccount(account.ID)
	// ...handle the error...

	// Delete account with the given ID and version
	if err = f3.DeleteAccount(account.ID, 123); err != nil {
		// If the error originates from Form3 REST API then form3.F3Error is returned
		// It contains additional context (if available) about what was the cause of
		// the error.
		// API documentation points out that error_code and error_message fields are
		// available for responses with 400 http status code but looks like responses
		// with other codes (e.g 404, 409) also provide these fields.
		fmt.Println(err) // -> http 409: code: 0, message=invalid version
	}
}
```

### Notes:

1. `form3_test.go` contains some general tests that do not run against provided fake account API.
2. `account_test.go` tests run against provided fake account API.

Possible improvements:

1. Introduce `context.Context` in `Client.Request` to allow granular control over HTTP request timeouts. Currently, it's not possible to specify the HTTP request timeout (default 15sec), unless you plug in custom HTTP client with a different timout via `WithHTTPClient` option.
2. `CreateAccount`, `FetchAccount`, `DeleteAccount` methods are registered directly on the client. If the project grows (support different resources) it would be better to separate the resources implementation from the client.
3. If the test suite starts to grow then something like `testify` could help to organise it and help with assertions in general.
