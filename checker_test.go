package httpcheck

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		w.Header().Add("hello", "goodbye")
		w.WriteHeader(http.StatusAccepted)
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

	case "/xml":
		body, err := xml.Marshal(testPerson{
			Name: "Some",
			Age:  30,
		})
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/xml")
		w.Write(body)
	case "/byte":
		w.Write([]byte("hello world"))
	case "/mirrorbody":
		body, _ := io.ReadAll(req.Body)
		w.Write(body)
	case "/cookies":
		http.SetCookie(w, &http.Cookie{
			Name:  "some",
			Value: "cookie",
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "other",
			Value: "secondcookie",
		})
	case "/redirect":
		w.Header().Set("Location", "https://localhost/redirect-test")
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "/timeout":
		time.Sleep(time.Second)
	case "/nothing":
	/* noop */
	default:
		panic("unexpected access")
	}
}

func newTestChecker(opts ...Option) *Checker {
	return New(&testHandler{}, opts...)
}

func TestNew(t *testing.T) {
	checker := New(&testHandler{})
	require.NotNil(t, checker)
	assert.Equal(t, DefaultClientTimeout, checker.client.Timeout)
}

func TestCheckerOption_ClientTimeout(t *testing.T) {
	t.Skip("skip this test because it expects failure.")
	timeout := 10 * time.Millisecond
	checker := New(&testHandler{}, ClientTimeout(timeout))
	assert.Equal(t, timeout, checker.client.Timeout)
	checker.Test(t, http.MethodGet, "/timeout").Check()
}

func TestCheckerOption_NoRedirect(t *testing.T) {
	checker := newTestChecker(NoRedirect())
	checker.Test(t, http.MethodGet, "/redirect").
		Check().
		HasStatus(http.StatusTemporaryRedirect)
	assert.Exactly(t, "/redirect", checker.request.URL.Path)
}

func TestChecker_Test(t *testing.T) {
	checker := newTestChecker()
	checker.Test(t, http.MethodGet, "/some")

	require.NotNil(t, checker.request)
	assert.Exactly(t, http.MethodGet, checker.request.Method)
	assert.Exactly(t, "/some", checker.request.URL.Path)
}

func TestChecker_TestRequest(t *testing.T) {
	checker := newTestChecker()
	request := &http.Request{
		Method: http.MethodGet,
	}

	checker.TestRequest(t, request)
	require.NotNil(t, checker.request)
	assert.Exactly(t, http.MethodGet, checker.request.Method)
	assert.Nil(t, checker.request.URL)
}
