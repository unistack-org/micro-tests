package http_handler

import (
	"net/http"
	"testing"

	swaggerui "go.unistack.org/micro-server-http/v3/handler/swagger-ui"
)

func TestTemplate(t *testing.T) {
	// t.Skip()
	h := http.NewServeMux()
	h.HandleFunc("/", swaggerui.Handler(""))
	if err := http.ListenAndServe(":8080", h); err != nil {
		t.Fatal(err)
	}
}
