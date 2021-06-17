package httpcheck

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
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
		body, _ := ioutil.ReadAll(req.Body)
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
	case "/nothing":

	}
}

func makeTestChecker(opts ...Option) *Checker {
	return New(&testHandler{}, opts...)
}

func TestNew(t *testing.T) {
	checker := New(&testHandler{})
	require.NotNil(t, checker)
	assert.Equal(t, DefaultClientTimeout, checker.client.Timeout)
}

func TestClientTimeout(t *testing.T) {
	timeout := 30 * time.Second
	checker := New(&testHandler{}, ClientTimeout(timeout))
	assert.Equal(t, timeout, checker.client.Timeout)
}

func TestNoRedirect(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(NoRedirect())
	checker.Test(mockT, http.MethodGet, "/redirect").
		Check().
		HasStatus(http.StatusTemporaryRedirect)

	assert.Exactly(t, "/redirect", checker.request.URL.Path)
	assert.False(t, mockT.Failed())
}

func TestTest(t *testing.T) {
	checker := makeTestChecker()
	checker.Test(t, "GET", "/some")

	require.NotNil(t, checker.request)
	assert.Exactly(t, "GET", checker.request.Method)
	assert.Exactly(t, "/some", checker.request.URL.Path)
}

func TestRequest(t *testing.T) {
	checker := makeTestChecker()
	request := &http.Request{
		Method: "GET",
	}

	checker.TestRequest(t, request)
	require.NotNil(t, checker.request)
	assert.Exactly(t, "GET", checker.request.Method)
	assert.Nil(t, checker.request.URL)
}

func TestWithHeader(t *testing.T) {
	checker := makeTestChecker()
	checker.Test(t, "GET", "/some").
		WithHeader("key", "value")

	assert.Equal(t, checker.request.Header.Get("key"), "value")
	assert.Equal(t, "", checker.request.Header.Get("unknown"))
}

func TestWithHeaders(t *testing.T) {
	headers := map[string]string{
		"key":             "value",
		"Authorization":   "Token abce-1234",
		"X-Custom-Header": "custom_value_000",
	}
	checker := makeTestChecker()
	checker.Test(t, "GET", "/some").
		WithHeaders(headers)

	for k, v := range headers {
		assert.Equal(t, checker.request.Header.Get(k), v)
	}

	assert.Equal(t, "", checker.request.Header.Get("unknown"))
}

func TestWithBasicAuth(t *testing.T) {
	checker := makeTestChecker()
	checker.Test(t, "GET", "/some").
		WithBasicAuth("alice", "secret")

	user, pass, ok := checker.request.BasicAuth()
	assert.True(t, ok)
	assert.Equal(t, "alice", user)
	assert.Equal(t, "secret", pass)

	h := base64.StdEncoding.EncodeToString([]byte("alice:secret"))
	assert.Equal(t, checker.request.Header.Get("Authorization"), "Basic "+h)
}

func TestWithCookie(t *testing.T) {
	checker := makeTestChecker()
	checker.Test(t, "GET", "/some").
		WithCookie("key", "value")

	cookie, err := checker.request.Cookie("key")
	require.NoError(t, err)
	assert.Equal(t, cookie.Value, "value")

	_, err = checker.request.Cookie("unknown")
	assert.Error(t, err)
}

func TestCheck(t *testing.T) {
	checker := makeTestChecker()
	checker.Test(t, "GET", "/some").
		Check()

	assert.NotNil(t, checker.response)
	assert.Exactly(t, 204, checker.response.StatusCode)
}

func TestHasStatus(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()
	checker.Test(mockT, "GET", "/some").
		Check().
		HasStatus(202)
	assert.True(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker()
	checker.Test(mockT, "GET", "/some").
		Check().
		HasStatus(204)
	assert.False(t, mockT.Failed())
}

func TestHasHeader(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()
	checker.Test(mockT, "GET", "/some").
		Check().
		HasHeader("some", "header")
	assert.False(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker()
	checker.Test(mockT, "GET", "/some").
		Check().
		HasHeader("some", "unknown")
	assert.True(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker()
	checker.Test(mockT, "GET", "/some").
		Check().
		HasHeader("unknown", "header")
	assert.True(t, mockT.Failed())
}

func TestHasHeaders(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()
	checker.Test(mockT, "GET", "/some").
		Check().
		HasHeaders(map[string]string{
			"some": "header",
		})
	assert.False(t, mockT.Failed())

	checker.Test(mockT, "GET", "/some").
		Check().
		HasHeaders(map[string]string{
			"unknown":   "header",
			"X-Unknown": "abc",
		})
	assert.True(t, mockT.Failed())
}

func TestHasNilHeaders(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()

	// nil is the zero value for maps, so isn't a problem passing nil as parameter to WithHeaders and HasHeaders
	checker.Test(t, "GET", "/some").
		WithHeaders(nil).
		Check().
		HasHeaders(nil)
	assert.False(t, mockT.Failed())
}

func TestHasCookie(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()
	checker.Test(mockT, "GET", "/some").
		Check().
		HasCookie("some", "cookie")
	assert.False(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker()
	checker.Test(mockT, "GET", "/some").
		Check().
		HasCookie("some", "unknown")
	assert.True(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker()
	checker.Test(mockT, "GET", "/some").
		Check().
		HasCookie("unknown", "cookie")
	assert.True(t, mockT.Failed())
}

func TestHasJson(t *testing.T) {
	mockT := new(testing.T)
	person := &testPerson{
		Name: "Some",
		Age:  30,
	}
	checker := makeTestChecker()
	result := checker.Test(mockT, "GET", "/json").
		Check().
		HasJSON(person)
	assert.False(t, mockT.Failed())

	person = &testPerson{
		Name: "Unknown",
		Age:  30,
	}
	result.HasJSON(person)
	assert.True(t, mockT.Failed())
}

func TestHasXml(t *testing.T) {
	mockT := new(testing.T)
	person := &testPerson{
		Name: "Some",
		Age:  30,
	}
	checker := makeTestChecker()
	result := checker.Test(mockT, "GET", "/xml").
		Check().
		HasXML(person)
	assert.False(t, mockT.Failed())

	person = &testPerson{
		Name: "Unknown",
		Age:  30,
	}
	result.HasXML(person)
	assert.True(t, mockT.Failed())
}

func TestHasBody(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()
	checker.Test(mockT, "GET", "/byte").
		Check().
		HasBody([]byte("hello world"))
}

func TestHasString(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()
	checker.Test(mockT, "GET", "/byte").
		Check().
		HasString("hello world")
}

func TestContainsBody(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()
	checker.Test(mockT, "GET", "/byte").
		Check().
		ContainsBody([]byte("llo")).
		ContainsBody([]byte("llo"))
	assert.False(t, mockT.Failed())
}

func TestNotContainsBody(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()
	checker.Test(mockT, "GET", "/byte").
		Check().
		NotContainsBody([]byte("aloha"))
	assert.False(t, mockT.Failed())
}

func TestContainsString(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()
	checker.Test(mockT, "GET", "/byte").
		Check().
		ContainsString("llo").
		ContainsString("llo")
	assert.False(t, mockT.Failed())
}

func TestNotContainsString(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()
	checker.Test(mockT, "GET", "/byte").
		Check().
		NotContainsString("aloha")
	assert.False(t, mockT.Failed())
}

func TestCb(t *testing.T) {
	mockT := new(testing.T)
	called := false
	checker := makeTestChecker()
	checker.Test(mockT, "GET", "/json").
		Check().
		Cb(func(response *http.Response) {
			called = true
		})
	assert.True(t, called)
}

func TestStringBody(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()

	checker.Test(mockT, "GET", "/mirrorbody").
		WithString("Test123").
		Check().
		HasString("Test123")

	assert.False(t, mockT.Failed())
}

func TestBytesBody(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()

	checker.Test(mockT, "GET", "/mirrorbody").
		WithBody([]byte("Test123")).
		Check().
		HasBody([]byte("Test123"))

	assert.False(t, mockT.Failed())
}

func TestJsonBody(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()

	person := &testPerson{
		Name: "Some",
		Age:  30,
	}

	checker.Test(mockT, "GET", "/mirrorbody").
		WithJSON(person).
		Check().
		HasJSON(person)

	assert.False(t, mockT.Failed())
}

func TestXmlBody(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()

	person := &testPerson{
		Name: "Some",
		Age:  30,
	}

	checker.Test(mockT, "GET", "/mirrorbody").
		WithXML(person).
		Check().
		HasXML(person)

	assert.False(t, mockT.Failed())
}

func TestCookies(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker()
	checker.PersistCookie("some")

	checker.Test(mockT, "GET", "/cookies").
		Check().
		HasCookie("some", "cookie").
		HasCookie("other", "secondcookie")
	assert.False(t, mockT.Failed())

	result := checker.Test(mockT, "GET", "/nothing").
		Check().
		HasCookie("some", "cookie")
	assert.False(t, mockT.Failed())

	result.UnpersistCookie("some")
	result = checker.Test(mockT, "GET", "/nothing").
		Check().
		HasCookie("some", "cookie")
	assert.True(t, mockT.Failed())

	result.HasCookie("other", "secondcookie")
	assert.True(t, mockT.Failed())
}
