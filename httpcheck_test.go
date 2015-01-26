package httpcheck

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type testHandler struct{}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(204)
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

func TestHasStatusFailed(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "http://localhost:3000/some")
	checker.Check()

	checker.HasStatus(202)
	assert.True(t, mockT.Failed())
}

func TestHasStatusOk(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "http://localhost:3000/some")
	checker.Check()

	checker.HasStatus(204)
	assert.False(t, mockT.Failed())
}
