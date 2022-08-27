package httpcheck

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// JSON body /////////////////////////////////////////////////////

// WithJSON adds a json encoded struct to the body.
func (tt *Tester) WithJSON(expected any) *Tester {
	encoded, err := json.Marshal(expected)
	require.NoError(tt.t, err)
	return tt.WithBody(encoded)
}

// Deprecated: WithJson adds a json encoded struct to the body.
//
//nolint:golint
func (tt *Tester) WithJson(expected any) *Tester {
	return tt.WithJSON(expected)
}

// HasJSON checks if the response body contains json with provided value.
func (tt *Tester) HasJSON(expected any) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	var b []byte
	switch v := expected.(type) {
	case string:
		b = []byte(v)
	case []byte:
		b = v
	default:
		b, err = json.Marshal(expected)
		require.NoError(tt.t, err)
	}
	if len(b) > 0 && len(body) == 0 {
		assert.Fail(tt.t, "response body is empty")
		return tt
	}
	assert.JSONEq(tt.t, string(b), string(body))

	return tt
}

// MustHasJSON checks if the response body contains json with provided value.
func (tt *Tester) MustHasJSON(expected any) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	var b []byte
	switch v := expected.(type) {
	case string:
		b = []byte(v)
	case []byte:
		b = v
	default:
		b, err = json.Marshal(expected)
		require.NoError(tt.t, err)
	}
	if len(b) > 0 && len(body) == 0 {
		require.Fail(tt.t, "response body is empty")
		return tt
	}
	require.JSONEq(tt.t, string(b), string(body))
	return tt
}

// Deprecated: HasJson checks if the response body contains json with provided value.
func (tt *Tester) HasJson(expected any) *Tester {
	return tt.HasJSON(expected)
}
