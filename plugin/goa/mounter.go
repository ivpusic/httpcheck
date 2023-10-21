package goa

import (
	"context"
	"log"
	"net/http"

	"github.com/ikawaha/httpcheck"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

var _ http.Handler = (*Mounter)(nil)

type (
	// type aliases
	decoder      = func(*http.Request) goahttp.Decoder
	encoder      = func(context.Context, http.ResponseWriter) goahttp.Encoder
	errorHandler = func(context.Context, http.ResponseWriter, error)
	formatter    = func(context.Context, error) goahttp.Statuser
	middleware   = func(http.Handler) http.Handler
)

type (
	// HandlerBuilder represents the goa http handler builder.
	HandlerBuilder func(goa.Endpoint, goahttp.Muxer, decoder, encoder, errorHandler, formatter) http.Handler
	// HandlerMounter represents the goa http handler mounter.
	HandlerMounter func(goahttp.Muxer, http.Handler)
)

// Mounter represents Goa v3 handler mounter.
type Mounter struct {
	Mux           goahttp.Muxer
	Middleware    []middleware
	Decoder       decoder
	Encoder       encoder
	ErrorHandler  errorHandler
	Formatter     formatter
	ClientOptions []httpcheck.Option
}

// Option represents options for API checker.
type Option func(m *Mounter)

// Muxer sets the goa http multiplexer.
func Muxer(mux goahttp.Muxer) Option {
	return func(m *Mounter) {
		m.Mux = mux
	}
}

// Decoder sets the decoder.
func Decoder(dec decoder) Option {
	return func(m *Mounter) {
		m.Decoder = dec
	}
}

// Encoder sets the encoder.
func Encoder(enc encoder) Option {
	return func(m *Mounter) {
		m.Encoder = enc
	}
}

// ErrorHandler sets the error handler.
func ErrorHandler(eh errorHandler) Option {
	return func(m *Mounter) {
		m.ErrorHandler = eh
	}
}

// Formatter sets the error handler.
func Formatter(fm formatter) Option {
	return func(m *Mounter) {
		m.Formatter = fm
	}
}

// NewMounter constructs a mounter of the goa endpoints.
func NewMounter(opts ...Option) *Mounter {
	ret := &Mounter{
		Mux:        goahttp.NewMuxer(),
		Middleware: []middleware{},
		Decoder:    goahttp.RequestDecoder,
		Encoder:    goahttp.ResponseEncoder,
		ErrorHandler: func(ctx context.Context, w http.ResponseWriter, err error) {
			log.Printf("ERROR: %v", err)
		},
	}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

// MountEndpoint mounts an endpoint handler and it's middleware.
func (m *Mounter) MountEndpoint(builder HandlerBuilder, mounter HandlerMounter, endpoint goa.Endpoint, middlewares ...middleware) {
	handler := builder(endpoint, m.Mux, m.Decoder, m.Encoder, m.ErrorHandler, m.Formatter)
	for _, v := range middlewares {
		handler = v(handler)
	}
	mounter(m.Mux, handler)
}

type EndpointModular struct {
	Builder            HandlerBuilder
	Mounter            HandlerMounter
	Endpoint           goa.Endpoint
	EndpointMiddleware []middleware
}

func (m *Mounter) Mount(e EndpointModular) {
	m.MountEndpoint(e.Builder, e.Mounter, e.Endpoint, e.EndpointMiddleware...)
}

// Use sets the middleware.
func (m *Mounter) Use(middleware func(http.Handler) http.Handler) {
	m.Middleware = append(m.Middleware, middleware)
}

func (m *Mounter) Handler() http.Handler {
	var h http.Handler = m.Mux
	for _, v := range m.Middleware {
		h = v(h)
	}
	return h
}

func (m *Mounter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler = m.Mux
	for _, v := range m.Middleware {
		h = v(h)
	}
	h.ServeHTTP(w, r)
}
