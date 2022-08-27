package httpcheck

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTester_WithBasicAuth(t *testing.T) {
	checker := newTestChecker()
	checker.Test(t, "GET", "/some").
		WithBasicAuth("alice", "secret")

	user, pass, ok := checker.request.BasicAuth()
	assert.True(t, ok)
	assert.Equal(t, "alice", user)
	assert.Equal(t, "secret", pass)

	h := base64.StdEncoding.EncodeToString([]byte("alice:secret"))
	assert.Equal(t, checker.request.Header.Get("Authorization"), "Basic "+h)
}

func TestTester_WithBearerAuth(t *testing.T) {
	checker := newTestChecker()
	checker.Test(t, "GET", "/some").
		WithBearerAuth("token")

	v := checker.request.Header.Get("Authorization")
	assert.Equal(t, "Bearer: token", v)
}
