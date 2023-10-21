package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("calc", func() {
	Description("The calc service performs operations on numbers")

	// Method describes a service method (endpoint)
	Method("multiply", func() {
		// Payload describes the method payload.
		// Here the payload is an object that consists of two fields.
		Payload(func() {
			// Attribute describes an object field
			Attribute("a", Int, "Left operand")
			Attribute("b", Int, "Right operand")
			Required("a", "b")
		})

		// Result describes the method result.
		// Here the result is a simple integer value.
		Result(Int)

		// HTTP describes the HTTP transport mapping.
		HTTP(func() {
			// Requests to the service consist of HTTP GET requests.
			// The payload fields are encoded as path parameters.
			GET("/multiply/{a}/{b}")
			// Responses use a "200 OK" HTTP status.
			// The result is encoded in the response body.
			Response(StatusOK)
		})
	})

	// Serve the file gen/http/openapi3.json for requests sent to
	// /openapi.json.
	Files("/openapi.json", "openapi3.json")
})
