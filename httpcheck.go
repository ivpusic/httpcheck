package httpcheck

import (
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

// Checker represents the HTTP checker.
type Checker struct {
	client   *http.Client
	request  *http.Request
	response *http.Response
	pcookies map[string]bool
	server   *httptest.Server
	handler  http.Handler
}

type checker struct {
	t *testing.T
	*Checker
}

// Option represents the option for the HTTP checker.
type Option func(*Checker)

// ClientTimeout sets the client timeout.
func ClientTimeout(d time.Duration) Option {
	return func(c *Checker) {
		c.client.Timeout = d
	}
}

// New creates a HTTP checker.
func New(handler http.Handler, options ...Option) *Checker {
	jar, _ := cookiejar.New(nil)
	instance := &Checker{
		//t: t,
		client: &http.Client{
			Timeout: time.Duration(5 * time.Second),
			Jar:     jar,
		},
		pcookies: map[string]bool{},
		server:   createServer(handler),
		handler:  handler,
	}

	return instance
}

func createServer(handler http.Handler) *httptest.Server {
	return httptest.NewUnstartedServer(handler)
}

// PersistCookie - enables a cookie to be preserved between requests
func (c *Checker) PersistCookie(cookie string) {
	c.pcookies[cookie] = true
}

// UnpersistCookie - the specified cookie will not be preserved during requests anymore
func (c *Checker) UnpersistCookie(cookie string) {
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
func (c *Checker) TestRequest(t *testing.T, request *http.Request) *checker {
	assert.NotNil(t, request, "Request nil")

	c.request = request
	return &checker{
		t:       t,
		Checker: c,
	}
}

// Test - Prepare for testing some part of code which lives on provided path and method.
func (c *Checker) Test(t *testing.T, method, path string) *checker {
	method = strings.ToUpper(method)
	request, err := http.NewRequest(method, c.GetURL()+path, nil)

	assert.Nil(t, err, "Failed to make new request")

	c.request = request
	return &checker{
		t:       t,
		Checker: c,
	}
}

// GetURL returns the server URL.
func (c *Checker) GetURL() string {
	return "http://" + c.server.Listener.Addr().String()
}

// headers ///////////////////////////////////////////////////////

// WithHeader - Will put header on request
func (c *checker) WithHeader(key, value string) *checker {
	c.request.Header.Set(key, value)
	return c
}

// HasHeader - Will check if response contains header on provided key with provided value
func (c *checker) HasHeader(key, expectedValue string) *checker {
	value := c.response.Header.Get(key)
	assert.Exactly(c.t, expectedValue, value)

	return c
}

// cookies ///////////////////////////////////////////////////////

// HasCookie - Will put cookie on request
func (c *checker) HasCookie(key, expectedValue string) *checker {
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
func (c *checker) WithCookie(key, value string) *checker {
	c.request.AddCookie(&http.Cookie{
		Name:  key,
		Value: value,
	})

	return c
}

// status ////////////////////////////////////////////////////////

// HasStatus - Will ckeck if response status is equal to provided
func (c *checker) HasStatus(status int) *checker {
	assert.Exactly(c.t, status, c.response.StatusCode)
	return c
}

// json body /////////////////////////////////////////////////////

// WithJSON - Will add the json-encoded struct to the body
func (c *checker) WithJSON(value interface{}) *checker {
	encoded, err := json.Marshal(value)
	assert.Nil(c.t, err)
	return c.WithBody(encoded)
}

// HasJSON - Will check if body contains json with provided value
func (c *checker) HasJSON(value interface{}) *checker {
	body, err := ioutil.ReadAll(c.response.Body)
	assert.Nil(c.t, err)

	valueBytes, err := json.Marshal(value)
	assert.Nil(c.t, err)
	assert.Equal(c.t, string(valueBytes), string(body))

	return c
}

// xml //////////////////////////////////////////////////////////

// WithXML - Adds a XML encoded body to the request
func (c *checker) WithXML(value interface{}) *checker {
	encoded, err := xml.Marshal(value)
	assert.Nil(c.t, err)
	return c.WithBody(encoded)
}

// HasXML - Will check if body contains xml with provided value
func (c *checker) HasXML(value interface{}) *checker {
	body, err := ioutil.ReadAll(c.response.Body)
	assert.Nil(c.t, err)

	valueBytes, err := xml.Marshal(value)
	assert.Nil(c.t, err)
	assert.Equal(c.t, string(valueBytes), string(body))

	return c
}

// body //////////////////////////////////////////////////////////

// WithBody - Adds the []byte data to the body
func (c *checker) WithBody(body []byte) *checker {
	c.request.Body = newClosingBuffer(body)
	c.request.ContentLength = int64(len(body))
	return c
}

// HasBody - Will check if body contains provided []byte data
func (c *checker) HasBody(body []byte) *checker {
	responseBody, err := ioutil.ReadAll(c.response.Body)

	assert.Nil(c.t, err)
	assert.Equal(c.t, body, responseBody)

	return c
}

// WithString - Adds the string to the body
func (c *checker) WithString(body string) *checker {
	c.request.Body = newClosingBufferString(body)
	c.request.ContentLength = int64(len(body))
	return c
}

// HasString - Convenience wrapper for HasBody
// Checks if body is equal to the given string
func (c *checker) HasString(body string) *checker {
	return c.HasBody([]byte(body))
}

// Check - Will make request to built request object.
// After request is made, it will save response object for future assertions
// Responsibility of this method is also to start and stop HTTP server
func (c *checker) Check() *checker {
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
func (c *checker) Cb(cb func(*http.Response)) {
	cb(c.response)
}
