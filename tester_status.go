package httpcheck

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// status ////////////////////////////////////////////////////////

// HasStatus checks if the response status is equal to that provided.
func (tt *Tester) HasStatus(expected int) *Tester {
	assert.Exactly(tt.t, expected, tt.response.StatusCode)
	return tt
}

// MustHasStatus checks if the response status is equal to that provided.
func (tt *Tester) MustHasStatus(expected int) *Tester {
	require.Exactly(tt.t, expected, tt.response.StatusCode)
	return tt
}
