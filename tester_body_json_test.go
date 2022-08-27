package httpcheck

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTester_WithJSON(t *testing.T) {
	mockT := new(testing.T)
	checker := newTestChecker()

	person := &testPerson{
		Name: "Some",
		Age:  30,
	}

	checker.Test(mockT, "GET", "/mirrorbody").
		WithJSON(person).
		Check().
		HasJSON(person).
		HasJSON(`{"Age":30, "Name":"Some"}`)

	assert.False(t, mockT.Failed())
}

func TestHasJSON(t *testing.T) {
	testdata := []struct {
		name     string
		expected testPerson
		want     bool
	}{
		{
			name: "OK",
			expected: testPerson{
				Name: "Some",
				Age:  30,
			},
			want: false,
		},
		{
			name: "NG",
			expected: testPerson{
				Name: "Unknown",
				Age:  30,
			},
			want: true,
		},
	}
	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)
			checker := newTestChecker()
			checker.Test(mockT, "GET", "/json").
				Check().
				HasJSON(tt.expected)
			assert.Equal(t, tt.want, mockT.Failed())
		})
	}

	t.Run("OK: json string", func(t *testing.T) {
		checker := newTestChecker()
		checker.Test(t, "GET", "/json").
			Check().
			HasJSON(`{"Age": 30, "Name": "Some"}`)
	})

	t.Run("NG: json string", func(t *testing.T) {
		mockT := new(testing.T)
		checker := newTestChecker()
		checker.Test(mockT, "GET", "/json").
			Check().
			HasJSON(`{"Age": 20, "Name": "Some"}`)
		assert.True(t, mockT.Failed())
	})
}

func TestTester_MustHasJSON(t *testing.T) {
	t.Skip("skip this test because it expects failure.")
	checker := newTestChecker()
	checker.Test(t, "GET", "/json").
		Check().
		MustHasJSON(`{"Age": 20, "Name": "Some"}`).
		Cb(func(response *http.Response) {
			t.Fatal("it is expected that this assertion will not be executed.")
		})
}
