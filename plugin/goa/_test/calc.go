package calc

import (
	"context"

	calcsvc "github.com/ikawaha/httpcheck/plugin/goa/calc/gen/calc"
)

// calc service example implementation.
// The example methods log the requests and return zero values.
type calcSvc struct{}

// NewCalc returns the calc service implementation.
func NewCalc() calcsvc.Service {
	return &calcSvc{}
}

// Multiply implements multiply.
func (s *calcSvc) Multiply(_ context.Context, p *calcsvc.MultiplyPayload) (int, error) {
	return p.A * p.B, nil
}
