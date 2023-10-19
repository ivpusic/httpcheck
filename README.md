[![Go](https://github.com/ikawaha/httpcheck/actions/workflows/test.yml/badge.svg)](https://github.com/ikawaha/httpcheck/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/ikawaha/httpcheck.svg)](https://pkg.go.dev/github.com/ikawaha/httpcheck)
# httpcheck

[supertest](https://github.com/visionmedia/supertest) inspired library for testing HTTP servers.

A Fork from [ivpusic/httpcheck](https://github.com/ivpusic/httpcheck) with following changes:

* Change to set testing.T when generating the request instead of the constructor,
* Fix to prevent incorrect method chain,
* Add to the timeout option of the client to the checker.


## How to install?
```
go get github.com/ikawaha/httpcheck
```

## API Documentation
[![Go Reference](https://pkg.go.dev/badge/github.com/ikawaha/httpcheck.svg)](https://pkg.go.dev/github.com/ikawaha/httpcheck)

## How to use?

### Basic example
```Go
package main

import (
	"github.com/ikawaha/httpcheck"
)

func TestExample(t *testing.T) {
	// testHandler should be instance of http.Handler
	checker := httpcheck.New(&testHandler{})

	checker.Test(t, "GET", "/some/url").
		WithHeader("key", "value").
		WithCookie("key", "value").
		Check().
		HasStatus(200).
		HasCookie("key", "expectedValue").
		HasHeader("key", "expectedValue").
		HasJSON(&someType{})
}
```

### Include body

#### String
```Go
package main

import (
	"github.com/ivpusic/httpcheck"
)

func TestExample(t *testing.T) {
	checker := httpcheck.New(&testHandler{})

	checker.Test(t, "GET", "/some/url").
		WithString("Hello!")
		Check().
		HasStatus(200)
}
```

#### JSON
```Go
package main

import (
	"github.com/ivpusic/httpcheck"
)

func TestExample(t *testing.T) {
	checker := httpcheck.New(&testHandler{})

	data := &someStruct{
		field1: "hi",
	}

	checker.Test(t, "GET", "/some/url").
		WithJSON(data)
		Check().
		HasStatus(200)
}
```

#### XML
```Go
package main

import (
	"github.com/ivpusic/httpcheck"
)

func TestExample(t *testing.T) {
	checker := httpcheck.New(&testHandler{})

	data := &someStruct{
		field1: "hi",
	}

	checker.Test(t, "GET", "/some/url").
		WithXML(data)
		Check().
		HasStatus(200)
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
	checker := httpcheck.New(&testHandler{})

	checker.TestRequest(t, &http.Request{ /* fields */ }).
		Check().
		HasStatus(200)
}
```

### Define callback
```Go
package main

import (
	"net/http"
	"github.com/ikawaha/httpcheck"
)

func TestExample(t *testing.T) {
	checker := httpcheck.New(&testHandler{})

	checker.Test(t, "GET", "/some/url").
		Check().
		HasStatus(200).
		HasBody([]byte("some body")).
		Cb(func(response *http.Response) { /* do something */ })
}
```

---
License MIT
