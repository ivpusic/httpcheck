package main

import (
	"fmt"
	"github.com/ivpusic/httpcheck"
	"net/http"
	"testing"
)

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
}

func main() {
	server := httpcheck.New(&testing.T{}, &handler{}, ":3000")

	server.Test("GET", "/some/path").
		WithHeader("key", "value").
		WithCookie("key", "value").
		WithCookie("key", "value").
		Check().
		HasCookie("key", "value").
		Cb(
		func(response *http.Response) {
			fmt.Println("aaa")
		})
}
