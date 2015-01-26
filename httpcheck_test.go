package httpcheck

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type testPerson struct {
	Name string
	Age  int
}

type testHandler struct{}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/some":
		http.SetCookie(w, &http.Cookie{
			Name:  "some",
			Value: "cookie",
		})
		w.Header().Add("some", "header")
		w.WriteHeader(204)
	case "/json":
		body, err := json.Marshal(testPerson{
			Name: "Some",
			Age:  30,
		})

		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	case "/byte":
		w.Write([]byte("hello world"))
	}
}

func makeTestChecker(t *testing.T) *Checker {
	handler := &testHandler{}
	port := ":3000"
	return New(t, handler, port)
}

func TestNew(t *testing.T) {
	handler := &testHandler{}
	port := ":3000"
	checker := New(t, handler, port)

	assert.NotNil(t, checker)
	assert.Exactly(t, t, checker.t)
	assert.Exactly(t, handler, checker.handler)
	assert.Exactly(t, port, checker.port)
	assert.NotNil(t, checker.server)
}

func TestTest(t *testing.T) {
	checker := makeTestChecker(t)
	checker.Test("GET", "http://localhost:3000/some")

	assert.NotNil(t, checker.request)
	assert.Exactly(t, "GET", checker.request.Method)
	assert.Exactly(t, "/some", checker.request.URL.Path)
}

func TestRequest(t *testing.T) {
	checker := makeTestChecker(t)
	request := &http.Request{
		Method: "GET",
	}

	checker.TestRequest(request)
	assert.NotNil(t, checker.request)
	assert.Exactly(t, "GET", checker.request.Method)
	assert.Nil(t, checker.request.URL)
}

func TestWithHeader(t *testing.T) {
	checker := makeTestChecker(t)
	checker.Test("GET", "http://localhost:3000/some")

	checker.WithHeader("key", "value")

	assert.Equal(t, checker.request.Header.Get("key"), "value")
	assert.Equal(t, "", checker.request.Header.Get("unknown"))
}

func TestWithCookie(t *testing.T) {
	checker := makeTestChecker(t)
	checker.Test("GET", "http://localhost:3000/some")

	checker.WithCookie("key", "value")

	cookie, err := checker.request.Cookie("key")
	assert.Nil(t, err)
	assert.Equal(t, cookie.Value, "value")

	cookie, err = checker.request.Cookie("unknown")
	assert.NotNil(t, err)
}

func TestCheck(t *testing.T) {
	checker := makeTestChecker(t)
	checker.Test("GET", "http://localhost:3000/some")
	checker.Check()

	assert.NotNil(t, checker.response)
	assert.Exactly(t, 204, checker.response.StatusCode)
}

func TestHasStatus(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "http://localhost:3000/some")
	checker.Check()

	checker.HasStatus(202)
	assert.True(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker(mockT)
	checker.Test("GET", "http://localhost:3000/some")
	checker.Check()

	checker.HasStatus(204)
	assert.False(t, mockT.Failed())
}

func TestHasHeader(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "http://localhost:3000/some")
	checker.Check()

	checker.HasHeader("some", "header")
	assert.False(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker(mockT)
	checker.Test("GET", "http://localhost:3000/some")
	checker.Check()

	checker.HasHeader("some", "unknown")
	assert.True(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker(mockT)
	checker.Test("GET", "http://localhost:3000/some")
	checker.Check()

	checker.HasHeader("unknown", "header")
	assert.True(t, mockT.Failed())
}

func TestHasCookie(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "http://localhost:3000/some")
	checker.Check()

	checker.HasCookie("some", "cookie")
	assert.False(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker(mockT)
	checker.Test("GET", "http://localhost:3000/some")
	checker.Check()

	checker.HasCookie("some", "unknown")
	assert.True(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker(mockT)
	checker.Test("GET", "http://localhost:3000/some")
	checker.Check()

	checker.HasCookie("unknown", "cookie")
	assert.True(t, mockT.Failed())
}

func TestHasJson(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "http://localhost:3000/json")
	checker.Check()

	person := &testPerson{
		Name: "Some",
		Age:  30,
	}
	checker.HasJson(person)
	assert.False(t, mockT.Failed())

	person = &testPerson{
		Name: "Unknown",
		Age:  30,
	}
	checker.HasJson(person)
	assert.True(t, mockT.Failed())
}

func TestHasBody(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "http://localhost:3000/byte")
	checker.Check()

	checker.HasBody([]byte("hello world"))
}

func TestCb(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "http://localhost:3000/json")
	checker.Check()

	called := false
	checker.Cb(func(response *http.Response) {
		called = true
	})

	assert.True(t, called)
}
