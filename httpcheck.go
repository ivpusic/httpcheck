package httpcheck

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type checker struct {
	client   *http.Client
	request  *http.Request
	response *http.Response
	pcookies map[string]bool
	server   *httptest.Server
	handler  http.Handler
}

// Checker represents the HTTP Checker.
type Checker struct {
	t *testing.T
	*checker
}

// Option represents the option for the HTTP Checker.
type Option func(*checker)

// ClientTimeout sets the client timeout.
func ClientTimeout(d time.Duration) Option {
	return func(c *checker) {
		c.client.Timeout = d
	}
}

// New creates a HTTP Checker.
func New(handler http.Handler, options ...Option) *checker {
	jar, _ := cookiejar.New(nil)
	checker := &checker{
		client: &http.Client{
			Timeout: time.Duration(5 * time.Second),
			Jar:     jar,
		},
		pcookies: map[string]bool{},
		server:   createServer(handler),
		handler:  handler,
	}
	for _, v := range options {
		v(checker)
	}
	return checker
}

func createServer(handler http.Handler) *httptest.Server {
	return httptest.NewUnstartedServer(handler)
}

// PersistCookie - enables a cookie to be preserved between requests
func (c *checker) PersistCookie(cookie string) {
	c.pcookies[cookie] = true
}

// UnpersistCookie - the specified cookie will not be preserved during requests anymore
func (c *checker) UnpersistCookie(cookie string) {
	delete(c.pcookies, cookie)
}

// Will run HTTP server
func (c *Checker) run() {
	//log.Println("running server")
	c.server.Start()
}

// Will stop HTTP server
func (c *Checker) stop() {
	//log.Println("stopping server")
	c.server.Close()
	c.server = createServer(c.handler)
}

// make request /////////////////////////////////////////////////

// TestRequest - If you want to provide you custom http.Request instance, you can do it using this method
// In this case internal http.Request instance won't be created, and passed instane will be used
// for making request
func (c *checker) TestRequest(t *testing.T, request *http.Request) *Checker {
	assert.NotNil(t, request, "Request nil")

	c.request = request
	return &Checker{
		t:       t,
		checker: c,
	}
}

// Test - Prepare for testing some part of code which lives on provided path and method.
func (c *checker) Test(t *testing.T, method, path string) *Checker {
	method = strings.ToUpper(method)
	request, err := http.NewRequest(method, c.GetURL()+path, nil)

	assert.Nil(t, err, "Failed to make new request")

	c.request = request
	return &Checker{
		t:       t,
		checker: c,
	}
}

// GetURL returns the server URL.
func (c *checker) GetURL() string {
	return "http://" + c.server.Listener.Addr().String()
}

// headers ///////////////////////////////////////////////////////

// WithHeader - Will put header on request
func (c *Checker) WithHeader(key, value string) *Checker {
	c.request.Header.Set(key, value)
	return c
}

// HasHeader - Will check if response contains header on provided key with provided value
func (c *Checker) HasHeader(key, expectedValue string) *Checker {
	value := c.response.Header.Get(key)
	assert.Exactly(c.t, expectedValue, value)

	return c
}

// WithBasicAuth - Alias for the basic auth request header.
func (c *Checker) WithBasicAuth(user, pass string) *Checker {
	var b bytes.Buffer
	b.WriteString(user)
	b.WriteString(":")
	b.WriteString(pass)
	return c.WithHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString(b.Bytes()))
}

// cookies ///////////////////////////////////////////////////////

// HasCookie - Will put cookie on request
func (c *Checker) HasCookie(key, expectedValue string) *Checker {
	found := false
	for _, cookie := range c.client.Jar.Cookies(c.request.URL) {
		if cookie.Name == key && cookie.Value == expectedValue {
			found = true
			break
		}
	}
	assert.True(c.t, found)

	return c
}

// WithCookie - Will check if response contains cookie with provided key and value
func (c *Checker) WithCookie(key, value string) *Checker {
	c.request.AddCookie(&http.Cookie{
		Name:  key,
		Value: value,
	})

	return c
}

// status ////////////////////////////////////////////////////////

// HasStatus - Will ckeck if response status is equal to provided
func (c *Checker) HasStatus(status int) *Checker {
	assert.Exactly(c.t, status, c.response.StatusCode)
	return c
}

// json body /////////////////////////////////////////////////////

// WithJSON - Will add the json-encoded struct to the body
func (c *Checker) WithJSON(value interface{}) *Checker {
	encoded, err := json.Marshal(value)
	assert.Nil(c.t, err)
	return c.WithBody(encoded)
}

// WithJson - deprecated
func (c *Checker) WithJson(value interface{}) *Checker {
	return c.WithJSON(value)
}

// HasJSON - Will check if body contains json with provided value
func (c *Checker) HasJSON(value interface{}) *Checker {
	body, err := ioutil.ReadAll(c.response.Body)
	assert.Nil(c.t, err)

	valueBytes, err := json.Marshal(value)
	assert.Nil(c.t, err)
	assert.Equal(c.t, string(valueBytes), string(body))

	return c
}

// HasJson - deprecated
func (c *Checker) HasJson(value interface{}) *Checker {
	return c.HasJSON(value)
}

// xml //////////////////////////////////////////////////////////

// WithXML - Adds a XML encoded body to the request
func (c *Checker) WithXML(value interface{}) *Checker {
	encoded, err := xml.Marshal(value)
	assert.Nil(c.t, err)
	return c.WithBody(encoded)
}

// WithXml - deprecated
func (c *Checker) WithXml(value interface{}) *Checker {
	return c.WithXML(value)
}

// HasXML - Will check if body contains xml with provided value
func (c *Checker) HasXML(value interface{}) *Checker {
	body, err := ioutil.ReadAll(c.response.Body)
	assert.Nil(c.t, err)

	valueBytes, err := xml.Marshal(value)
	assert.Nil(c.t, err)
	assert.Equal(c.t, string(valueBytes), string(body))

	return c
}

// HasXml - deprecated
func (c *Checker) HasXml(value interface{}) *Checker {
	return c.HasXML(value)
}

// body //////////////////////////////////////////////////////////

// WithBody - Adds the []byte data to the body
func (c *Checker) WithBody(body []byte) *Checker {
	c.request.Body = newClosingBuffer(body)
	c.request.ContentLength = int64(len(body))
	return c
}

// HasBody - Will check if body contains provided []byte data
func (c *Checker) HasBody(body []byte) *Checker {
	responseBody, err := ioutil.ReadAll(c.response.Body)

	assert.Nil(c.t, err)
	assert.Equal(c.t, body, responseBody)

	return c
}

// WithString - Adds the string to the body
func (c *Checker) WithString(body string) *Checker {
	c.request.Body = newClosingBufferString(body)
	c.request.ContentLength = int64(len(body))
	return c
}

// HasString - Convenience wrapper for HasBody
// Checks if body is equal to the given string
func (c *Checker) HasString(body string) *Checker {
	return c.HasBody([]byte(body))
}

// Check - Will make request to built request object.
// After request is made, it will save response object for future assertions
// Responsibility of this method is also to start and stop HTTP server
func (c *Checker) Check() *Checker {
	// start server in new goroutine
	c.run()
	defer c.stop()

	newJar, _ := cookiejar.New(nil)

	for name := range c.pcookies {
		for _, oldCookie := range c.client.Jar.Cookies(c.request.URL) {
			if name == oldCookie.Name {
				newJar.SetCookies(c.request.URL, []*http.Cookie{oldCookie})
				break
			}
		}
	}

	c.client.Jar = newJar
	response, err := c.client.Do(c.request)
	if err != nil {
		println(err.Error())
		c.t.FailNow()
	}

	// assert.Nil(c.t, err, "Failed while making new request.", err)

	// save response for assertion checks
	c.response = response

	return c
}

// Cb - Will call provided callback function with current response
func (c *Checker) Cb(cb func(*http.Response)) {
	cb(c.response)
}
