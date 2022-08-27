package httpcheck

import (
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// cookies ///////////////////////////////////////////////////////

// WithCookie puts cookie on the request.
func (tt *Tester) WithCookie(key, value string) *Tester {
	tt.request.AddCookie(&http.Cookie{
		Name:  key,
		Value: value,
	})
	return tt
}

// HasCookie checks if the response contains cookie with provided key and value.
func (tt *Tester) HasCookie(key, expectedValue string) *Tester {
	found := false
	for _, cookie := range tt.client.Jar.Cookies(tt.request.URL) {
		if cookie.Name == key && cookie.Value == expectedValue {
			found = true
			break
		}
	}
	assert.True(tt.t, found, "not found, expected key:"+key+", value:"+expectedValue)
	return tt
}

// MustHasCookie checks if the response contains cookie with provided key and value.
func (tt *Tester) MustHasCookie(key, expectedValue string) *Tester {
	found := false
	for _, cookie := range tt.client.Jar.Cookies(tt.request.URL) {
		if cookie.Name == key && cookie.Value == expectedValue {
			found = true
			break
		}
	}
	require.True(tt.t, found, "not found, expected key:"+key+", value:"+expectedValue)
	return tt
}
