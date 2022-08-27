package httpcheck

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

// MustHasHeader checks if the response contains header on provided key with provided value.
func (tt *Tester) MustHasHeader(key, expectedValue string) *Tester {
	value := tt.response.Header.Get(key)
	require.Exactly(tt.t, expectedValue, value)
	return tt
}

// HasHeaders checks if the response contains a provided headers map.
func (tt *Tester) HasHeaders(headers map[string]string) *Tester {
	for key, expectedValue := range headers {
		value := tt.response.Header.Get(key)
		assert.Exactly(tt.t, expectedValue, value)
	}
	return tt
}

// MustHasHeaders checks if the response contains a provided headers map.
func (tt *Tester) MustHasHeaders(headers map[string]string) *Tester {
	for key, expectedValue := range headers {
		value := tt.response.Header.Get(key)
		require.Exactly(tt.t, expectedValue, value)
	}
	return tt
}
