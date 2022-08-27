package httpcheck

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTester_Check(t *testing.T) {
	checker := newTestChecker()
	checker.Test(t, http.MethodGet, "/some").
		Check()

	assert.NotNil(t, checker.response)
	assert.Exactly(t, http.StatusAccepted, checker.response.StatusCode)
}

func TestTester_Cb(t *testing.T) {
	mockT := new(testing.T)
	called := false
	checker := newTestChecker()
	checker.Test(mockT, http.MethodGet, "/json").
		Check().
		Cb(func(response *http.Response) {
			called = true
		})
	assert.True(t, called)
}
