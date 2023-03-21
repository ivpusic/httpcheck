package httpcheck

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTester_WithMultipart(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/multipart", func(w http.ResponseWriter, r *http.Request) {
		mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		require.NoError(t, err)
		assert.Equal(t, "multipart/form-data", mediaType)
		if strings.HasPrefix(mediaType, "multipart/") {
			mr := multipart.NewReader(r.Body, params["boundary"])
			for {
				p, err := mr.NextPart()
				if err != nil {
					assert.ErrorIs(t, err, io.EOF)
					break
				}
				got, err := io.ReadAll(p)
				require.NoError(t, err)
				switch n := p.FormName(); n {
				case "part_1":
					assert.Equal(t, "part_1", p.FormName())
					assert.Equal(t, "value_1", string(got))
				case "part_2":
					assert.Equal(t, "neko_small.png", p.FileName())
					b, err := os.ReadFile("./testdata/neko_small.png")
					require.NoError(t, err)
					assert.Equal(t, b, got)
				}
				p.Close()
			}
		}
		w.WriteHeader(http.StatusOK)
	})
	checker := New(mux)
	checker.Test(t, "POST", "/multipart").
		WithMultipart(
			FieldPart{FieldName: "part_1", Value: "value_1"},
			FilePart{FieldName: "part_2", FileName: "./testdata/neko_small.png"},
		).
		Check().HasStatus(http.StatusOK)
}
