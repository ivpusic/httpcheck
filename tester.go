package httpcheck

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	b, err := ioutil.ReadAll(response.Body)
	assert.NoError(tt.t, err)

	tt.response = response
	tt.response.Body = ioutil.NopCloser(bytes.NewReader(b))
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

// headers ///////////////////////////////////////////////////////

// WithHostHeader puts "Host" header on the request.
func (tt *Tester) WithHostHeader(value string) *Tester {
	tt.request.Host = value
	return tt
}

// WithHeader puts header on the request.
func (tt *Tester) WithHeader(key, value string) *Tester {
	tt.request.Header.Set(key, value)
	return tt
}

// WithHeaders puts a map of headers on the request.
func (tt *Tester) WithHeaders(headers map[string]string) *Tester {
	for key, value := range headers {
		tt.request.Header.Set(key, value)
	}
	return tt
}

// HasHeader checks if the response contains header on provided key with provided value.
func (tt *Tester) HasHeader(key, expectedValue string) *Tester {
	value := tt.response.Header.Get(key)
	assert.Exactly(tt.t, expectedValue, value)

	return tt
}

// HasHeaders checks if the response contains a provided headers map
func (tt *Tester) HasHeaders(headers map[string]string) *Tester {
	for key, expectedValue := range headers {
		value := tt.response.Header.Get(key)
		assert.Exactly(tt.t, expectedValue, value)
	}

	return tt
}

// authorization headers ///////////////////////////////////////////////////////

// WithBasicAuth is an alias to set basic auth in the request header.
func (tt *Tester) WithBasicAuth(user, pass string) *Tester {
	var b bytes.Buffer
	b.WriteString(user)
	b.WriteString(":")
	b.WriteString(pass)
	return tt.WithHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString(b.Bytes()))
}

// WithBearerAuth is an alias to set bearer auth in the request header.
func (tt *Tester) WithBearerAuth(token string) *Tester {
	return tt.WithHeader("Authorization", "Bearer: "+token)
}

// cookies ///////////////////////////////////////////////////////

// HasCookie checks if the response contains cookie with provided key and value.
func (tt *Tester) HasCookie(key, expectedValue string) *Tester {
	found := false
	for _, cookie := range tt.client.Jar.Cookies(tt.request.URL) {
		if cookie.Name == key && cookie.Value == expectedValue {
			found = true
			break
		}
	}
	assert.True(tt.t, found)

	return tt
}

// WithCookie puts cookie on the request.
func (tt *Tester) WithCookie(key, value string) *Tester {
	tt.request.AddCookie(&http.Cookie{
		Name:  key,
		Value: value,
	})

	return tt
}

// status ////////////////////////////////////////////////////////

// HasStatus checks if the response status is equal to provided.
func (tt *Tester) HasStatus(status int) *Tester {
	assert.Exactly(tt.t, status, tt.response.StatusCode)
	return tt
}

// JSON body /////////////////////////////////////////////////////

// WithJSON adds a json encoded struct to the body.
func (tt *Tester) WithJSON(value interface{}) *Tester {
	encoded, err := json.Marshal(value)
	assert.NoError(tt.t, err)
	return tt.WithBody(encoded)
}

// WithJson adds a json encoded struct to the body. (deprecated)
// nolint:golint
func (tt *Tester) WithJson(value interface{}) *Tester {
	return tt.WithJSON(value)
}

// HasJSON checks if the response body contains json with provided value.
func (tt *Tester) HasJSON(value interface{}) *Tester {
	b, err := ioutil.ReadAll(tt.response.Body)
	assert.NoError(tt.t, err)
	tt.response.Body.Close()

	valueBytes, err := json.Marshal(value)
	assert.NoError(tt.t, err)
	assert.Equal(tt.t, string(valueBytes), string(b))

	tt.response.Body = ioutil.NopCloser(bytes.NewReader(b))
	return tt
}

// HasJson checks if the response body contains json with provided value. (deprecated)
//nolint:golint
func (tt *Tester) HasJson(value interface{}) *Tester {
	return tt.HasJSON(value)
}

// XML //////////////////////////////////////////////////////////

// WithXML adds a xml encoded body to the request.
func (tt *Tester) WithXML(value interface{}) *Tester {
	encoded, err := xml.Marshal(value)
	assert.NoError(tt.t, err)
	return tt.WithBody(encoded)
}

// WithXml  adds a xml encoded body to the request. (deprecated)
//nolint:golint
func (tt *Tester) WithXml(value interface{}) *Tester {
	return tt.WithXML(value)
}

// HasXML checks if body contains xml with provided value.
func (tt *Tester) HasXML(value interface{}) *Tester {
	b, err := ioutil.ReadAll(tt.response.Body)
	assert.NoError(tt.t, err)
	tt.response.Body.Close()

	valueBytes, err := xml.Marshal(value)
	assert.NoError(tt.t, err)
	assert.Equal(tt.t, string(valueBytes), string(b))

	tt.response.Body = ioutil.NopCloser(bytes.NewReader(b))
	return tt
}

// HasXml checks if body contains xml with provided value. (deprecated)
//nolint:golint
func (tt *Tester) HasXml(value interface{}) *Tester {
	return tt.HasXML(value)
}

// body //////////////////////////////////////////////////////////

// WithBody adds the []byte data to the body.
func (tt *Tester) WithBody(body []byte) *Tester {
	tt.request.Body = ioutil.NopCloser(bytes.NewReader(body))
	tt.request.ContentLength = int64(len(body))
	return tt
}

// HasBody checks if the body is equal to provided []byte data.
func (tt *Tester) HasBody(body []byte) *Tester {
	b, err := ioutil.ReadAll(tt.response.Body)
	assert.NoError(tt.t, err)
	tt.response.Body.Close()
	assert.Equal(tt.t, body, b)

	tt.response.Body = ioutil.NopCloser(bytes.NewReader(b))
	return tt
}

// ContainsBody checks if the body contains provided [] byte data.
func (tt *Tester) ContainsBody(segment []byte) *Tester {
	b, err := ioutil.ReadAll(tt.response.Body)
	assert.NoError(tt.t, err)
	tt.response.Body.Close()

	if !bytes.Contains(b, segment) {
		assert.Fail(tt.t, fmt.Sprintf("%#v does not contain %#v", b, segment))
	}

	tt.response.Body = ioutil.NopCloser(bytes.NewReader(b))
	return tt
}

// NotContainsBody checks if the body does not contain provided [] byte data.
func (tt *Tester) NotContainsBody(segment []byte) *Tester {
	b, err := ioutil.ReadAll(tt.response.Body)
	assert.NoError(tt.t, err)
	tt.response.Body.Close()

	if bytes.Contains(b, segment) {
		assert.Fail(tt.t, fmt.Sprintf("%#v contains %#v", b, segment))
	}

	tt.response.Body = ioutil.NopCloser(bytes.NewReader(b))
	return tt
}

// WithString adds the string to the body.
func (tt *Tester) WithString(body string) *Tester {
	tt.request.Body = ioutil.NopCloser(strings.NewReader(body))
	tt.request.ContentLength = int64(len(body))
	return tt
}

// HasString converts the response to a string type and then compares it with the given string.
func (tt *Tester) HasString(body string) *Tester {
	b, err := ioutil.ReadAll(tt.response.Body)
	assert.NoError(tt.t, err)
	tt.response.Body.Close()
	assert.Equal(tt.t, body, string(b))

	tt.response.Body = ioutil.NopCloser(bytes.NewReader(b))
	return tt
}

// ContainsString converts the response to a string type and then checks it containing the given string.
func (tt *Tester) ContainsString(substr string) *Tester {
	b, err := ioutil.ReadAll(tt.response.Body)
	assert.NoError(tt.t, err)
	tt.response.Body.Close()

	assert.Contains(tt.t, string(b), substr)

	tt.response.Body = ioutil.NopCloser(bytes.NewReader(b))
	return tt
}

// NotContainsString converts the response to a string type and then checks if it does not
// contain the given string.
func (tt *Tester) NotContainsString(substr string) *Tester {
	b, err := ioutil.ReadAll(tt.response.Body)
	assert.NoError(tt.t, err)

	assert.NotContains(tt.t, string(b), substr)

	tt.response.Body = ioutil.NopCloser(bytes.NewReader(b))
	return tt
}
