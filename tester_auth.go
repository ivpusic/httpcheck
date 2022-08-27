package httpcheck

import (
	"bytes"
	"encoding/base64"
)

// authorization headers ///////////////////////////////////////////////////////

// WithBasicAuth is an alias to set basic auth in the request header.
func (tt *Tester) WithBasicAuth(user, pass string) *Tester {
	var b bytes.Buffer
	b.WriteString(user)
	b.WriteString(":")
	b.WriteString(pass)
	return tt.WithHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString(b.Bytes()))
}

// WithBearerAuth is an alias to set bearer auth in the request header.
func (tt *Tester) WithBearerAuth(token string) *Tester {
	return tt.WithHeader("Authorization", "Bearer: "+token)
}
