package httpcheck

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTester_WithHeader(t *testing.T) {
	checker := newTestChecker()
	checker.Test(t, http.MethodGet, "/some").
		WithHeader("key", "value")

	assert.Equal(t, checker.request.Header.Get("key"), "value")
	assert.Equal(t, "", checker.request.Header.Get("unknown"))
}

func TestTester_WithHeaders(t *testing.T) {
	headers := map[string]string{
		"key":             "value",
		"Authorization":   "Token abce-1234",
		"X-Custom-Header": "custom_value_000",
	}
	checker := newTestChecker()
	checker.Test(t, http.MethodGet, "/some").
		WithHeaders(headers)

	for k, v := range headers {
		assert.Equal(t, checker.request.Header.Get(k), v)
	}
	assert.Equal(t, "", checker.request.Header.Get("unknown"))
}

func TestHasHeader(t *testing.T) {
	type pair struct {
		key   string
		value string
	}
	testdata := []struct {
		name     string
		method   string
		expected pair
		want     bool
	}{
		{
			name:   "OK: response has the specified header",
			method: http.MethodGet,
			expected: pair{
				key:   "some",
				value: "header",
			},
			want: false,
		},
		{
			name:   "OK: response has the specified header",
			method: http.MethodGet,
			expected: pair{
				key:   "hello",
				value: "goodbye",
			},
			want: false,
		},
		{
			name:   "NG: response does not have the specified header value",
			method: http.MethodGet,
			expected: pair{
				key:   "some",
				value: "aloha",
			},
			want: true,
		},
		{
			name:   "NG: response does not have the specified header key",
			method: http.MethodGet,
			expected: pair{
				key:   "hello",
				value: "header",
			},
			want: true,
		},
	}
	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)
			checker := newTestChecker()
			checker.Test(mockT, http.MethodGet, "/some").
				Check().
				HasHeader(tt.expected.key, tt.expected.value)
			assert.Equal(t, tt.want, mockT.Failed())
		})
	}
}

func TestTester_MustHasHeader(t *testing.T) {
	t.Skip("skip this test because it expects failure.")
	checker := newTestChecker()
	checker.Test(t, http.MethodGet, "/some").
		Check().
		MustHasHeader("hello", "goodbye"). // ← it requires to be true, but it fails.
		Cb(func(response *http.Response) {
			t.Fatal("it is expected that this assertion will not be executed.")
		})
}

func TestTester_HasHeaders(t *testing.T) {
	testdata := []struct {
		name     string
		method   string
		expected map[string]string
		want     bool
	}{
		{
			name:   "OK: response has the specified header",
			method: http.MethodGet,
			expected: map[string]string{
				"some":  "header",
				"hello": "goodbye",
			},
			want: false,
		},
		{
			name:   "OK: response includes the specified headers",
			method: http.MethodGet,
			expected: map[string]string{
				"hello": "goodbye",
			},
			want: false,
		},
		{
			name:   "NG: response does not have the specified headers",
			method: http.MethodGet,
			expected: map[string]string{
				"some":      "header",
				"X-unknown": "aloha",
			},
			want: true,
		},
	}
	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)
			checker := newTestChecker()
			checker.Test(mockT, "GET", "/some").
				Check().
				HasHeaders(tt.expected)
			assert.Equal(t, tt.want, mockT.Failed())
		})
	}
}

func TestTester_HasNilHeaders(t *testing.T) {
	mockT := new(testing.T)
	checker := newTestChecker()

	// nil is the zero value for maps, so isn't a problem passing nil as parameter to WithHeaders and HasHeaders
	checker.Test(t, "GET", "/some").
		WithHeaders(nil).
		Check().
		HasHeaders(nil)
	assert.False(t, mockT.Failed())
}

func TestTester_MustHasHeaders(t *testing.T) {
	t.Skip("skip this test because it expects failure.")
	checker := newTestChecker()
	checker.Test(t, http.MethodGet, "/some").
		Check().
		MustHasHeaders(map[string]string{ // ← it requires to be true, but it fails.
			"hello": "unknown",
		}).
		Cb(func(response *http.Response) {
			t.Fatal("it is expected that this assertion will not be executed.")
		})
}
