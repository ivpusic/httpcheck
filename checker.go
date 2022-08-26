package httpcheck

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	// DefaultClientTimeout is the default timeout for requests made by checker.
	DefaultClientTimeout = 5 * time.Second
)

// Option represents the option for the checker.
type Option func(*Checker)

// ClientTimeout sets the client timeout.
func ClientTimeout(d time.Duration) Option {
	return func(c *Checker) {
		c.client.Timeout = d
	}
}

// CheckRedirect sets the policy of redirection to the HTTP client.
func CheckRedirect(policy func(req *http.Request, via []*http.Request) error) Option {
	return func(c *Checker) {
		c.client.CheckRedirect = policy
	}
}

// NoRedirect is the alias of the following:
//
//	CheckRedirect(func(req *http.Request, via []*http.Request) error {
//	    return http.ErrUseLastResponse
//	})
//
// Client returns ErrUseLastResponse, the next request is not sent and the most recent
// response is returned with its body unclosed.
func NoRedirect() Option {
	return CheckRedirect(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	})
}

// Checker represents the HTTP checker without testing.T.
type Checker struct {
	client   *http.Client
	request  *http.Request
	response *http.Response
	pcookies map[string]bool
	server   *httptest.Server
	handler  http.Handler
}

// New creates an HTTP Checker.
func New(h http.Handler, options ...Option) *Checker {
	jar, _ := cookiejar.New(nil)
	ret := &Checker{
		client: &http.Client{
			Timeout: DefaultClientTimeout,
			Jar:     jar,
		},
		pcookies: map[string]bool{},
		server:   createServer(h),
		handler:  h,
	}
	for _, v := range options {
		v(ret)
	}
	return ret
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

// make request /////////////////////////////////////////////////

// TestRequest - If you want to provide you custom http.Request instance, you can do it using this method
// In this case internal http.Request instance won't be created, and passed instance will be used
// for making request
func (c *Checker) TestRequest(t *testing.T, request *http.Request) *Tester {
	require.NotNil(t, request, "request is nil")

	c.request = request
	return &Tester{
		t:       t,
		Checker: c,
	}
}

// Test - Prepare for testing some part of code which lives on provided path and method.
func (c *Checker) Test(t *testing.T, method, path string) *Tester {
	method = strings.ToUpper(method)
	request, err := http.NewRequest(method, c.GetURL()+path, nil)

	require.NoError(t, err, "failed to make new request")

	c.request = request
	return &Tester{
		t:       t,
		Checker: c,
	}
}

// GetURL returns the server URL.
func (c *Checker) GetURL() string {
	return "http://" + c.server.Listener.Addr().String()
}
