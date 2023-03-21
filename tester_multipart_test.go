package httpcheck

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestTester_WithMultipart(t *testing.T) {
	testdata := []struct {
		name     string
		method   string
		expected int
		want     bool
	}{
		{
			name:     "OK: expected status 200",
			method:   http.MethodPost,
			expected: http.StatusOK,
			want:     true,
		},
	}
	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)
			checker := newTestChecker()
			checker.Test(mockT, "POST", "/some").
				WithMultipart([]FormData{
					{Key: "param1", Value: "value1"},
					{Key: "param2", Value: "value2", FileName: "items.csv"},
					{
						Key:   "param3",
						Value: `{"key1": "value1", "key2": "value2"}`,
					},
				}...).
				Check().
				HasStatus(tt.expected)
			assert.Equal(t, tt.want, mockT.Failed())
		})
	}
}
