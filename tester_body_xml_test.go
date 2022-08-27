package httpcheck

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTester_WithXML(t *testing.T) {
	mockT := new(testing.T)
	checker := newTestChecker()

	person := &testPerson{
		Name: "Some",
		Age:  30,
	}

	checker.Test(mockT, "GET", "/mirrorbody").
		WithXML(person).
		Check().
		HasXML(person)

	assert.False(t, mockT.Failed())
}

func TestTester_HasXML(t *testing.T) {
	mockT := new(testing.T)
	person := &testPerson{
		Name: "Some",
		Age:  30,
	}
	checker := newTestChecker()
	result := checker.Test(mockT, "GET", "/xml").
		Check().
		HasXML(person)
	assert.False(t, mockT.Failed())

	person = &testPerson{
		Name: "Unknown",
		Age:  30,
	}
	result.HasXML(person)
	assert.True(t, mockT.Failed())
}

func TestTester_MustHasXml(t *testing.T) {
	t.Skip("skip this test because it expects failure.")
	checker := newTestChecker()
	checker.Test(t, "GET", "/xml").
		Check().
		MustHasXML(&testPerson{
			Name: "Unknown",
			Age:  30,
		}).
		Cb(func(response *http.Response) {
			t.Fatal("it is expected that this assertion will not be executed.")
		})
}
