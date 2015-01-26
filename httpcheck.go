package httpcheck

import (
	"encoding/json"
	"github.com/braintree/manners"
	"github.com/ivpusic/golog"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

type (
	Checker struct {
		t        *testing.T
		handler  http.Handler
		addr     string
		server   *manners.GracefulServer
		request  *http.Request
		response *http.Response
	}

	Callback func(*http.Response)
)

var (
	logger = golog.GetLogger("github.com/ivpusic/golog")
)

func New(t *testing.T, handler http.Handler, addr string) *Checker {
	logger.Level = golog.INFO

	instance := &Checker{
		t:       t,
		handler: handler,
		addr:    addr,
	}
	instance.server = manners.NewServer()

	return instance
}

// Will run HTTP server
func (c *Checker) run() {
	logger.Debug("running server")
	c.server.ListenAndServe(c.addr, c.handler)
}

// Will stop HTTP server
func (c *Checker) stop() {
	logger.Debug("stopping server")
	c.server.Shutdown <- true
	// todo: solve race condition
	time.Sleep(1 * time.Millisecond)
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
	request, err := http.NewRequest(method, path, nil)

	assert.Nil(c.t, err, "Failed to make new request")

	c.request = request
	return c
}

// headers ///////////////////////////////////////////////////////

// Will put header on request
func (c *Checker) WithHeader(key, value string) *Checker {
	c.request.Header.Set(key, value)
	return c
}

// Will check if response contains header on provided key with provided value
func (c *Checker) HasHeader(key, expectedValue string) *Checker {
	value := c.response.Header.Get(key)
	assert.Exactly(c.t, expectedValue, value)

	return c
}

// cookies ///////////////////////////////////////////////////////

// Will put cookie on request
func (c *Checker) HasCookie(key, expectedValue string) *Checker {
	responseCookiesMap := cookiesToMap(c.response.Cookies())
	cookieValue, ok := responseCookiesMap[key]

	assert.True(c.t, ok)
	assert.Exactly(c.t, expectedValue, cookieValue)

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

// Will ckeck if body contains json with provided value
func (c *Checker) HasJson(value interface{}) *Checker {
	body, err := ioutil.ReadAll(c.response.Body)
	assert.Nil(c.t, err)

	valueBytes, err := json.Marshal(value)
	assert.Nil(c.t, err)
	assert.Equal(c.t, string(valueBytes), string(body))

	return c
}

// body //////////////////////////////////////////////////////////

// Will check if body contains provided []byte data
func (c *Checker) HasBody(body []byte) *Checker {
	responseBody, err := ioutil.ReadAll(c.response.Body)

	assert.Nil(c.t, err)
	assert.Equal(c.t, body, responseBody)

	return c
}

// Will make reqeust to built request object.
// After request is made, it will save response object for future assertions
// Responsibility of this method is also to start and stop HTTP server
func (c *Checker) Check() *Checker {
	// start server in new goroutine
	go c.run()

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.Do(c.request)
	assert.Nil(c.t, err, "Failed while making new request.", err)

	// save response for assertion checks
	c.response = response

	// stop server
	c.stop()

	return c
}

// Will call provided callback function with current response
func (c *Checker) Cb(cb Callback) {
	cb(c.response)
}
