package httpcheck

import (
	"bytes"
	"encoding/xml"
	"io"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// XML //////////////////////////////////////////////////////////

// WithXML adds a xml encoded body to the request.
func (tt *Tester) WithXML(body any) *Tester {
	encoded, err := xml.Marshal(body)
	require.NoError(tt.t, err)
	return tt.WithBody(encoded)
}

// Deprecated: WithXml adds a xml encoded body to the request.
func (tt *Tester) WithXml(body any) *Tester {
	return tt.WithXML(body)
}

// HasXML checks if body contains xml with provided value.
func (tt *Tester) HasXML(expected any) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	b, err := xml.Marshal(expected)
	require.NoError(tt.t, err)
	assert.Equal(tt.t, string(b), string(body))
	return tt
}

// MustHasXML checks if body contains xml with provided value.
func (tt *Tester) MustHasXML(expected any) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	b, err := xml.Marshal(expected)
	require.NoError(tt.t, err)
	require.Equal(tt.t, string(b), string(body))
	return tt
}

// Deprecated: HasXml checks if body contains xml with provided value.
//
//nolint:golint
func (tt *Tester) HasXml(expected any) *Tester {
	return tt.HasXML(expected)
}
