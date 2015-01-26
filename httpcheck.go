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
		port     string
		server   *manners.GracefulServer
		request  *http.Request
		response *http.Response
	}

	Callback func(*http.Response)
)

var (
	logger = golog.GetLogger("github.com/ivpusic/golog")
)

func New(t *testing.T, handler http.Handler, port string) *Checker {
	instance := &Checker{
		t:       t,
		handler: handler,
		port:    port,
	}
	instance.server = manners.NewServer()

	return instance
}

func (c *Checker) run() {
	logger.Info("running")
	c.server.ListenAndServe(c.port, c.handler)
}

func (c *Checker) stop() {
	logger.Info("stopping")
	c.server.Shutdown <- true
}

// make request /////////////////////////////////////////////////
func (c *Checker) TestRequest(request *http.Request) *Checker {
	assert.NotNil(c.t, request, "Request nil")

	c.request = request
	return c
}

func (c *Checker) Test(method, path string) *Checker {
	method = strings.ToUpper(method)
	request, err := http.NewRequest(method, path, nil)

	assert.Nil(c.t, err, "Failed to make new request")

	c.request = request
	return c
}

// headers ///////////////////////////////////////////////////////
func (c *Checker) WithHeader(key, value string) *Checker {
	c.request.Header.Set(key, value)
	return c
}

func (c *Checker) HasHeader(key, expectedValue string) *Checker {
	value := c.response.Header.Get(key)
	assert.Exactly(c.t, expectedValue, value)

	return c
}

// cookies ///////////////////////////////////////////////////////
func (c *Checker) HasCookie(key, expectedValue string) *Checker {
	responseCookiesMap := cookiesToMap(c.response.Cookies())
	cookieValue, ok := responseCookiesMap[key]

	assert.True(c.t, ok)
	assert.Exactly(c.t, expectedValue, cookieValue)

	return c
}

func (c *Checker) WithCookie(key, value string) *Checker {
	c.request.AddCookie(&http.Cookie{
		Name:  key,
		Value: value,
	})

	return c
}

// status ////////////////////////////////////////////////////////
func (c *Checker) HasStatus(status int) *Checker {
	assert.Exactly(c.t, status, c.response.StatusCode)
	return c
}

// json body /////////////////////////////////////////////////////
func (c *Checker) HasJson(value interface{}) *Checker {
	body, err := ioutil.ReadAll(c.response.Body)
	assert.Nil(c.t, err)

	valueBytes, err := json.Marshal(value)
	assert.Nil(c.t, err)
	assert.Equal(c.t, string(valueBytes), string(body))

	return c
}

// body //////////////////////////////////////////////////////////
func (c *Checker) HasBody(body []byte) *Checker {
	responseBody, err := ioutil.ReadAll(c.response.Body)

	assert.Nil(c.t, err)
	assert.Equal(c.t, body, responseBody)

	return c
}

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

func (c *Checker) Cb(cb Callback) {
	cb(c.response)
}
