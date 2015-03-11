package httpcheck

import (
	"bytes"
)

// Internal helper struct for submitting a body to the http.Request
type closingBuffer struct {
	*bytes.Buffer
}

func (cb *closingBuffer) Close() (err error) {
	return
}

func newClosingBuffer(body []byte) *closingBuffer {
	return &closingBuffer{
		bytes.NewBuffer(body),
	}
}

func newClosingBufferString(body string) *closingBuffer {
	return &closingBuffer{
		bytes.NewBufferString(body),
	}
}
