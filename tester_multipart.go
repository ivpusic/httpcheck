package httpcheck

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/stretchr/testify/require"
)

// Parter is the interface that create a multipart part.
type Parter interface {
	Part(mw *multipart.Writer) error
}

// FieldPart represents a multipart part of the field type part.
type FieldPart struct {
	FieldName string
	Value     string
}

// Part returns a multipart part.
func (p FieldPart) Part(mw *multipart.Writer) error {
	w, err := mw.CreateFormField(p.FieldName)
	if err != nil {
		return fmt.Errorf("failed to creat form field: %w", err)
	}
	if _, err := w.Write([]byte(p.Value)); err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	return nil
}

// FilePart represents a multipart part of the file type part.
type FilePart struct {
	FieldName string
	FileName  string
}

// Part returns a multipart part.
func (p FilePart) Part(mw *multipart.Writer) error {
	w, err := mw.CreateFormFile(p.FieldName, p.FileName)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	b, err := os.ReadFile(p.FileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}

// WithMultipart add a multipart data to the body.
func (tt *Tester) WithMultipart(part Parter, parts ...Parter) *Tester {
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	require.NoError(tt.t, part.Part(mw))
	for _, v := range parts {
		require.NoError(tt.t, v.Part(mw))
	}
	require.NoError(tt.t, mw.Close())
	return tt.WithHeader("Content-Type", mw.FormDataContentType()).WithBody(b.Bytes())
}
