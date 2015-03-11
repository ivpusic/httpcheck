# httpcheck
[![Build Status](https://travis-ci.org/ivpusic/httpcheck.svg?branch=master)](https://travis-ci.org/ivpusic/httpcheck)

[supertest](https://github.com/tj/supertest) inspired library for testing HTTP servers.

## How to install?
```
go get github.com/ivpusic/httpcheck
```

## API Documentation
[godoc](https://godoc.org/github.com/ivpusic/httpcheck)

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

### Include body

#### String
```Go
package main

import (
	"github.com/ivpusic/httpcheck"
)

func TestExample(t *testing.T) {
	checker := httpcheck.New(t, &testHandler{}, ":3000")

	checker.Test("GET", "/some/url").
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
	checker := httpcheck.New(t, &testHandler{}, ":3000")

	data := &someStruct{
		field1: "hi",
	}

	checker.Test("GET", "/some/url").
		WithJson(data)
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
	checker := httpcheck.New(t, &testHandler{}, ":3000")

	data := &someStruct{
		field1: "hi",
	}

	checker.Test("GET", "/some/url").
		WithXml(data)
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
- Implement fix/feature
- Write tests for fix/feature
- Make sure all tests are passing
- Send Pull Request

# License
*MIT*
