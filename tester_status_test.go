package httpcheck

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTester_HasStatus(t *testing.T) {
	testdata := []struct {
		name     string
		method   string
		expected int
		want     bool
	}{
		{
			name:     "OK: expected status 202",
			method:   http.MethodGet,
			expected: http.StatusAccepted,
			want:     false,
		},
		{
			name:     "NG: expected status 200 but got 202",
			method:   http.MethodGet,
			expected: http.StatusOK,
			want:     true,
		},
	}
	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)
			checker := newTestChecker()
			checker.Test(mockT, "GET", "/some").
				Check().
				HasStatus(tt.expected)
			assert.Equal(t, tt.want, mockT.Failed())
		})
	}
}

func TestTester_MustHasStatus(t *testing.T) {
	t.Skip("skip this test because it expects failure.")
	checker := newTestChecker()
	checker.Test(t, "GET", "/some").
		Check().
		MustHasStatus(111). // ← it requires to be true, but it fails.
		Cb(func(r *http.Response) {
			t.Fatal("it is expected that this assertion will not be executed.")
		})
}
