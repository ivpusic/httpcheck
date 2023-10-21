package calc

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/ikawaha/httpcheck"
	"github.com/ikawaha/httpcheck/plugin/goa"
	gen "github.com/ikawaha/httpcheck/plugin/goa/calc/gen/calc"
	"github.com/ikawaha/httpcheck/plugin/goa/calc/gen/http/calc/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMounter(t *testing.T) {
	m := goa.NewMounter()
	m.Mount(goa.EndpointModular{
		Builder:  server.NewMultiplyHandler,
		Mounter:  server.MountMultiplyHandler,
		Endpoint: gen.NewMultiplyEndpoint(NewCalc()),
	})
	httpcheck.New(m).Test(t, "GET", "/multiply/3/5").
		Check().
		MustHasStatus(http.StatusOK).
		Cb(func(r *http.Response) {
			b, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			i, err := strconv.Atoi(strings.TrimSpace(string(b)))
			assert.NoError(t, err, string(b))
			assert.Equal(t, 3*5, i)
		})
}
