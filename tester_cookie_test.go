package httpcheck

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTester_WithCookie(t *testing.T) {
	checker := newTestChecker()
	checker.Test(t, http.MethodGet, "/some").
		WithCookie("key", "value")

	cookie, err := checker.request.Cookie("key")
	require.NoError(t, err)
	assert.Equal(t, cookie.Value, "value")

	_, err = checker.request.Cookie("unknown")
	assert.Error(t, err)
}

func TestTester_Cookies(t *testing.T) {
	mockT := new(testing.T)
	checker := newTestChecker()
	checker.PersistCookie("some")

	checker.Test(mockT, "GET", "/cookies").
		Check().
		HasCookie("some", "cookie").
		HasCookie("other", "secondcookie")
	assert.False(t, mockT.Failed())

	result := checker.Test(mockT, "GET", "/nothing").
		Check().
		HasCookie("some", "cookie")
	assert.False(t, mockT.Failed())

	result.UnpersistCookie("some")
	result = checker.Test(mockT, "GET", "/nothing").
		Check().
		HasCookie("some", "cookie")
	assert.True(t, mockT.Failed())

	result.HasCookie("other", "secondcookie")
	assert.True(t, mockT.Failed())
}

func TestTester_HasCookie(t *testing.T) {
	type pair struct {
		key   string
		value string
	}
	testdata := []struct {
		name     string
		expected pair
		want     bool
	}{
		{
			name: "OK: response has some cookie.",
			expected: pair{
				key:   "some",
				value: "cookie",
			},
			want: false,
		},
		{
			name: "NG: response does not have the specified cookie value.",
			expected: pair{
				key:   "some",
				value: "unknown",
			},
			want: true,
		},
		{
			name: "NG: response does not have the specified cookie key.",
			expected: pair{
				key:   "unknown",
				value: "cookie",
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
				HasCookie(tt.expected.key, tt.expected.value)
			assert.Equal(t, tt.want, mockT.Failed())
		})
	}
}

func TestTester_MustHasCookie(t *testing.T) {
	t.Skip("skip this test because it expects failure.")
	checker := newTestChecker()
	checker.Test(t, "GET", "/some").
		Check().
		MustHasCookie("some", "unknown").
		Cb(func(response *http.Response) {
			t.Fatal("it is expected that this assertion will not be executed.")
		})
}
