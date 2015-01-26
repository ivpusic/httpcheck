# httpcheck

[supertest](https://github.com/tj/supertest) inspired library for testing HTTP servers.

## How to install?
```
go get github.com/ivpusic/httpcheck
```

## How to use?

### Basic example
```Go
package main

import (
	"github.com/ivpusic/httpcheck"
)

func TestExample(t *testing.T) {
	checker := httpcheck.New(t, &testHandler{}, ":3000")

	checker.Test("GET", "/some/url").
		WithHeader("key", "value").
		WithCookie("key", "value").
		Check().
		HasStatus(200).
		HasCookie("key", "expectedValue").
		HasHeader("key", "expectedValue").
		HasJson(&someType{})
}
```

### Provide ``*http.Request`` instance
```Go
package main

import (
	"net/http"
	"github.com/ivpusic/httpcheck"
)

func TestExample(t *testing.T) {
	checker := httpcheck.New(t, &testHandler{}, ":3000")

	checker.TestRequest(&http.Request{ /* fields */ }).
		Check().
		HasStatus(200)
}
```

### Define callback
```Go
package main

import (
	"net/http"
	"github.com/ivpusic/httpcheck"
)

func TestExample(t *testing.T) {
	checker := httpcheck.New(t, &testHandler{}, ":3000")

	checker.Test("GET", "/some/url").
		Check().
		HasStatus(200).
		HasBody([]byte("some body")).
		Cb(func(response *http.Response) { /* do something */ })
}
```

## Contribution Guidelines
- Provide your fix/new feature
- Write tests
- Make sure all tests are passing
- Send Pull Request

# License
*MIT*
