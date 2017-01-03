package httpcheck

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"
	"time"

	"github.com/ivpusic/golog"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
)

type (
	Checker struct {
		t        *testing.T
		client   *http.Client
		request  *http.Request
		response *http.Response
		pcookies map[string]bool
		server   *httptest.Server
		handler  http.Handler
	}

	Callback func(*http.Response)
)

var (
	logger = golog.GetLogger("github.com/ivpusic/httpcheck")
)

func New(t *testing.T, handler http.Handler) *Checker {
	logger.Level = golog.INFO

	jar, _ := cookiejar.New(nil)
	instance := &Checker{
		t: t,
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

// enables a cookie to be preserved between requests
func (c *Checker) PersistCookie(cookie string) {
	c.pcookies[cookie] = true
}

// the specified cookie will not be preserved during requests anymore
func (c *Checker) UnpersistCookie(cookie string) {
	delete(c.pcookies, cookie)
}

// Will run HTTP server
func (c *Checker) run() {
	logger.Debug("running server")
	c.server.Start()
}

// Will stop HTTP server
func (c *Checker) stop() {
	logger.Debug("stopping server")
	c.server.Close()
	c.server = createServer(c.handler)
}

func (c *Checker) SetTesting(t *testing.T) *Checker {
	if t == nil {
		panic("testing.T is nil")
	}

	c.t = t
	return c
}

// make request /////////////////////////////////////////////////

// If you want to provide you custom http.Request instance, you can do it using this method
// In this case internal http.Request instance won't be created, and passed instane will be used
// for making request
func (c *Checker) TestRequest(request *http.Request) *Checker {
	assert.NotNil(c.t, request, "Request nil")

	c.request = request
	return c
}

// Prepare for testing some part of code which lives on provided path and method.
func (c *Checker) Test(method, path string) *Checker {
	method = strings.ToUpper(method)
	request, err := http.NewRequest(method, c.GetUrl()+path, nil)

	assert.Nil(c.t, err, "Failed to make new request")

	c.request = request
	return c
}

func (c *Checker) GetUrl() string {
	return "http://" + c.server.Listener.Addr().String()
}

// headers ///////////////////////////////////////////////////////

// Will put header on request
func (c *Checker) WithHeader(key, value string) *Checker {
	c.request.Header.Set(key, value)
	return c
}

// Will put a map of headers on request
func (c *Checker) WithHeaders(headers map[string]string) *Checker {
	for key, value := range headers {
		c.request.Header.Set(key, value)
	}
	return c
}

// Will check if response contains header on provided key with provided value
func (c *Checker) HasHeader(key, expectedValue string) *Checker {
	value := c.response.Header.Get(key)
	assert.Exactly(c.t, expectedValue, value)

	return c
}

// Will check if response contains a provided headers map
func (c *Checker) HasHeaders(headers map[string]string) *Checker {

	for key, expectedValue := range headers {
		value := c.response.Header.Get(key)
		assert.Exactly(c.t, expectedValue, value)
	}

	return c
}

// cookies ///////////////////////////////////////////////////////

// Will put cookie on request
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

// Will ckeck if response contains cookie with provided key and value
func (c *Checker) WithCookie(key, value string) *Checker {
	c.request.AddCookie(&http.Cookie{
		Name:  key,
		Value: value,
	})

	return c
}

// status ////////////////////////////////////////////////////////

// Will ckeck if response status is equal to provided
func (c *Checker) HasStatus(status int) *Checker {
	assert.Exactly(c.t, status, c.response.StatusCode)
	return c
}

// json body /////////////////////////////////////////////////////

// Will add the json-encoded struct to the body
func (c *Checker) WithJson(value interface{}) *Checker {
	encoded, err := json.Marshal(value)
	assert.Nil(c.t, err)
	return c.WithBody(encoded)
}

// Will ckeck if body contains json with provided value
func (c *Checker) HasJson(value interface{}) *Checker {
	body, err := ioutil.ReadAll(c.response.Body)
	assert.Nil(c.t, err)

	valueBytes, err := json.Marshal(value)
	assert.Nil(c.t, err)
	assert.Equal(c.t, string(valueBytes), string(body))

	return c
}

// xml //////////////////////////////////////////////////////////

// Adds a XML encoded body to the request
func (c *Checker) WithXml(value interface{}) *Checker {
	encoded, err := xml.Marshal(value)
	assert.Nil(c.t, err)
	return c.WithBody(encoded)
}

// Will ckeck if body contains xml with provided value
func (c *Checker) HasXml(value interface{}) *Checker {
	body, err := ioutil.ReadAll(c.response.Body)
	assert.Nil(c.t, err)

	valueBytes, err := xml.Marshal(value)
	assert.Nil(c.t, err)
	assert.Equal(c.t, string(valueBytes), string(body))

	return c
}

// body //////////////////////////////////////////////////////////

// Adds the []byte data to the body
func (c *Checker) WithBody(body []byte) *Checker {
	c.request.Body = newClosingBuffer(body)
	c.request.ContentLength = int64(len(body))
	return c
}

// Will check if body contains provided []byte data
func (c *Checker) HasBody(body []byte) *Checker {
	responseBody, err := ioutil.ReadAll(c.response.Body)

	assert.Nil(c.t, err)
	assert.Equal(c.t, body, responseBody)

	return c
}

// Adds the string to the body
func (c *Checker) WithString(body string) *Checker {
	c.request.Body = newClosingBufferString(body)
	c.request.ContentLength = int64(len(body))
	return c
}

// Convenience wrapper for HasBody
// Checks if body is equal to the given string
func (c *Checker) HasString(body string) *Checker {
	return c.HasBody([]byte(body))
}

// Will make reqeust to built request object.
// After request is made, it will save response object for future assertions
// Responsibility of this method is also to start and stop HTTP server
func (c *Checker) Check() *Checker {
	// start server in new goroutine
	c.run()
	defer c.stop()

	newJar, _ := cookiejar.New(nil)

	for name, _ := range c.pcookies {
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

// Will call provided callback function with current response
func (c *Checker) Cb(cb Callback) {
	cb(c.response)
}
