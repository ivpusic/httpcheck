package httpcheck

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
	case "/nothing":

	}
}

func makeTestChecker(t *testing.T) *Checker {
	handler := &testHandler{}
	return New(t, handler)
}

func TestNew(t *testing.T) {
	handler := &testHandler{}
	checker := New(t, handler)

	assert.NotNil(t, checker)
	assert.Exactly(t, t, checker.t)
}

func TestSetTesting(t *testing.T) {
	checker := New(nil, &testHandler{})
	checker.SetTesting(t)

	assert.NotNil(t, checker.t)
	assert.Exactly(t, t, checker.t)
}

func TestTest(t *testing.T) {
	checker := makeTestChecker(t)
	checker.Test("GET", "/some")

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
	checker.Test("GET", "/some")

	checker.WithHeader("key", "value")

	assert.Equal(t, checker.request.Header.Get("key"), "value")
	assert.Equal(t, "", checker.request.Header.Get("unknown"))
}

func TestWithHeaders(t *testing.T) {
	checker := makeTestChecker(t)
	checker.Test("GET", "/some")

	headers := map[string]string{
		"key":             "value",
		"Authorization":   "Token abce-1234",
		"X-Custom-Header": "custom_value_000",
	}

	checker.WithHeaders(headers)

	for k, v := range headers {
		assert.Equal(t, checker.request.Header.Get(k), v)
	}

	assert.Equal(t, "", checker.request.Header.Get("unknown"))
}

func TestWithCookie(t *testing.T) {
	checker := makeTestChecker(t)
	checker.Test("GET", "/some")

	checker.WithCookie("key", "value")

	cookie, err := checker.request.Cookie("key")
	assert.Nil(t, err)
	assert.Equal(t, cookie.Value, "value")

	cookie, err = checker.request.Cookie("unknown")
	assert.NotNil(t, err)
}

func TestCheck(t *testing.T) {
	checker := makeTestChecker(t)
	checker.Test("GET", "/some")
	checker.Check()

	assert.NotNil(t, checker.response)
	assert.Exactly(t, 204, checker.response.StatusCode)
}

func TestHasStatus(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "/some")
	checker.Check()

	checker.HasStatus(202)
	assert.True(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker(mockT)
	checker.Test("GET", "/some")
	checker.Check()

	checker.HasStatus(204)
	assert.False(t, mockT.Failed())
}

func TestHasHeader(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "/some")
	checker.Check()

	checker.HasHeader("some", "header")
	assert.False(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker(mockT)
	checker.Test("GET", "/some")
	checker.Check()

	checker.HasHeader("some", "unknown")
	assert.True(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker(mockT)
	checker.Test("GET", "/some")
	checker.Check()

	checker.HasHeader("unknown", "header")
	assert.True(t, mockT.Failed())
}

func TestHasHeaders(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "/some")
	checker.Check()

	checker.HasHeaders(map[string]string{
		"some": "header",
	})
	assert.False(t, mockT.Failed())

	checker.HasHeaders(map[string]string{
		"unknown":   "header",
		"X-Unknown": "abc",
	})
	assert.True(t, mockT.Failed())
}

func TestHasNilHeaders(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "/some")

	// nil is the zero value for maps, so isn't a problem passing nil as parameter to WithHeaders and HasHeaders

	checker.WithHeaders(nil)
	checker.Check()
	checker.HasHeaders(nil)

	assert.False(t, mockT.Failed())
}

func TestHasCookie(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "/some")
	checker.Check()

	checker.HasCookie("some", "cookie")
	assert.False(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker(mockT)
	checker.Test("GET", "/some")
	checker.Check()

	checker.HasCookie("some", "unknown")
	assert.True(t, mockT.Failed())

	mockT = new(testing.T)
	checker = makeTestChecker(mockT)
	checker.Test("GET", "/some")
	checker.Check()

	checker.HasCookie("unknown", "cookie")
	assert.True(t, mockT.Failed())
}

func TestHasJson(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "/json")
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

func TestHasXml(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "/xml")
	checker.Check()

	person := &testPerson{
		Name: "Some",
		Age:  30,
	}
	checker.HasXml(person)
	assert.False(t, mockT.Failed())

	person = &testPerson{
		Name: "Unknown",
		Age:  30,
	}
	checker.HasXml(person)
	assert.True(t, mockT.Failed())
}

func TestHasBody(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "/byte")
	checker.Check()

	checker.HasBody([]byte("hello world"))
}

func TestCb(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.Test("GET", "/json")
	checker.Check()

	called := false
	checker.Cb(func(response *http.Response) {
		called = true
	})

	assert.True(t, called)
}

func TestStringBody(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)

	checker.Test("GET", "/mirrorbody").
		WithString("Test123").
		Check().
		HasString("Test123")

	assert.False(t, mockT.Failed())
}

func TestBytesBody(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)

	checker.Test("GET", "/mirrorbody").
		WithBody([]byte("Test123")).
		Check().
		HasBody([]byte("Test123"))

	assert.False(t, mockT.Failed())
}

func TestJsonBody(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)

	person := &testPerson{
		Name: "Some",
		Age:  30,
	}

	checker.Test("GET", "/mirrorbody").
		WithJson(person).
		Check().
		HasJson(person)

	assert.False(t, mockT.Failed())
}

func TestXmlBody(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)

	person := &testPerson{
		Name: "Some",
		Age:  30,
	}

	checker.Test("GET", "/mirrorbody").
		WithXml(person).
		Check().
		HasXml(person)

	assert.False(t, mockT.Failed())
}

func TestCookies(t *testing.T) {
	mockT := new(testing.T)
	checker := makeTestChecker(mockT)
	checker.PersistCookie("some")

	checker.Test("GET", "/cookies")
	checker.Check()

	checker.HasCookie("some", "cookie")
	checker.HasCookie("other", "secondcookie")
	assert.False(t, mockT.Failed())

	checker.Test("GET", "/nothing")
	checker.Check()

	checker.HasCookie("some", "cookie")
	assert.False(t, mockT.Failed())

	checker.HasCookie("other", "secondcookie")
	assert.True(t, mockT.Failed())
}
