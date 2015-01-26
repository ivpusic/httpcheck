package httpcheck

import (
	"github.com/braintree/manners"
	"github.com/ivpusic/golog"
	"net/http"
	"strings"
	"testing"
	"time"
)

type Checker struct {
	t        *testing.T
	handler  http.Handler
	port     string
	server   *manners.GracefulServer
	request  *http.Request
	response *http.Response
	assert   struct {
		status  int
		cookies map[string]string
		headers map[string]string
	}
}

type Callback func(*http.Response)

var logger = golog.GetLogger("github.com/ivpusic/golog")

func New(t *testing.T, handler http.Handler, port string) *Checker {
	instance := &Checker{
		t:       t,
		handler: handler,
		port:    port,
	}
	instance.server = manners.NewServer()

	return instance
}

func (c *Checker) failExpected(expected, got interface{}) {
	c.t.Fatal("Expected %v but got %v", expected, got)
}

func (c *Checker) failMsg(args ...interface{}) {
	c.t.Fatal(args...)
}

func (c *Checker) run() {
	logger.Info("running")
	c.server.ListenAndServe(c.port, c.handler)
}

func (c *Checker) stop() {
	logger.Info("stopping")
	c.server.Shutdown <- true
}

func (c *Checker) TestRequest(request *http.Request) *Checker {
	if request == nil {
		c.failMsg("Request nil")
	}

	c.request = request
	return c
}

func (c *Checker) Test(method, path string) *Checker {
	method = strings.ToUpper(method)
	request, err := http.NewRequest(method, path, nil)

	if err != nil {
		c.failMsg("Failed to make new request")
	}

	c.request = request
	return c
}

func (c *Checker) WithCookie(key, value string) *Checker {
	c.request.AddCookie(&http.Cookie{
		Name:  key,
		Value: value,
	})

	return c
}

func (c *Checker) WithHeader(key, value string) *Checker {
	c.request.Header.Set(key, value)
	return c
}

func (c *Checker) HasCookie(key, value string) *Checker {
	c.assert.cookies[key] = value
	return c
}

func (c *Checker) HasHeader(key, value string) *Checker {
	c.assert.headers[key] = value
	return c
}

func (c *Checker) HasStatus(status int) *Checker {
	c.assert.status = status
	return c
}

func (c *Checker) HasJson(content string) *Checker {
	return c
}

func (c *Checker) Check() *Checker {
	// start server in new goroutine
	go c.run()
	time.Sleep(1000 * time.Millisecond)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.Do(c.request)

	if err != nil {
		c.failMsg("Failed while making new request.", err)
	}

	// check status
	if response.StatusCode != c.assert.status {
		c.failExpected(c.assert.status, response.StatusCode)
	}

	// check headers
	for k, v := range c.assert.headers {
		value := response.Header.Get(k)

		if value != v {
			c.failExpected(v, value)
		}
	}

	// check cookies
	responseCookiesMap := cookiesToMap(response.Cookies())
	for k, v := range c.assert.cookies {
		value, ok := responseCookiesMap[k]
		if !ok {
			c.failMsg("Cookie %s not found!", k)
		}

		if value != v {
			c.failExpected(v, value)
		}
	}

	// save response in case of callback
	c.response = response

	// stop server
	c.stop()

	return c
}

func (c *Checker) Cb(cb Callback) {
	cb(c.response)
}
