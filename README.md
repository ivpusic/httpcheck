# httpcheck

[supertest](https://github.com/tj/supertest) inspired library for testing HTTP servers.

## How to install?
```
go get github.com/ivpusic/httpcheck
```

## How to use?

### Basic example
```Go
func TestExample(t *testing.T) {
	checker := New(t, &testHandler{}, ":3000")

	checker.Test("GET", "http://localhost:3000/some/url").
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
func TestExample(t *testing.T) {
	checker := New(t, &testHandler{}, ":3000")

	checker.TestRequest(&http.Request{ /* fields */ }).
		Check().
		HasStatus(200)
}
```

### Define callback
```Go
func TestExample(t *testing.T) {
	checker := New(t, &testHandler{}, ":3000")

	checker.Test("GET", "http://localhost:3000/some/url").
		Check().
		HasStatus(200).
		HasBody([]byte("some body")).
		Cb(func(response *http.Response) { /* do something */ })
}
```
