package httpcheck

import (
	"bytes"
	"io"
	"net/http"
	"net/http/cookiejar"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tester represents the HTTP tester having testing.T.
type Tester struct {
	t *testing.T
	*Checker
}

// Check - Will make request to built request object.
// After request is made, it will save response object for future assertions
// Responsibility of this method is also to start and stop HTTP server
func (tt *Tester) Check() *Tester {
	// start server in new goroutine
	tt.run()
	defer tt.stop()

	newJar, _ := cookiejar.New(nil)
	for name := range tt.pcookies {
		for _, oldCookie := range tt.client.Jar.Cookies(tt.request.URL) {
			if name == oldCookie.Name {
				newJar.SetCookies(tt.request.URL, []*http.Cookie{oldCookie})
				break
			}
		}
	}

	tt.client.Jar = newJar
	response, err := tt.client.Do(tt.request)
	if err != nil {
		assert.FailNow(tt.t, err.Error())
	}
	defer response.Body.Close()

	// save response for assertion checks
	b, err := io.ReadAll(response.Body)
	require.NoError(tt.t, err)

	tt.response = response
	tt.response.Body = io.NopCloser(bytes.NewReader(b))
	return tt
}

// Cb - Will call provided callback function with current response
func (tt *Tester) Cb(callback func(*http.Response)) {
	callback(tt.response)
}

// Will run HTTP server
func (tt *Tester) run() {
	// log.Println("running server")
	tt.server.Start()
}

// Will stop HTTP server
func (tt *Tester) stop() {
	// log.Println("stopping server")
	tt.server.Close()
	tt.server = createServer(tt.handler)
}
