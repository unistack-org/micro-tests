// +build ignore

package rpc

import (
	"bytes"
	"net/http"
	"testing"
)

func TestRequestPayloadFromRequest(t *testing.T) {

	jsonUrlBytes := []byte(`{"key1":"val1","key2":"val2","name":"Test"}`)

	t.Run("extracting a json from a POST request with url params", func(t *testing.T) {
		r, err := http.NewRequest("POST", "http://localhost/my/path?key1=val1&key2=val2", bytes.NewReader(jsonBytes))
		if err != nil {
			t.Fatalf("Failed to created http.Request: %v", err)
		}

		extByte, err := requestPayload(r)
		if err != nil {
			t.Fatalf("Failed to extract payload from request: %v", err)
		}
		if string(extByte) != string(jsonUrlBytes) {
			t.Fatalf("Expected %v and %v to match", string(extByte), jsonUrlBytes)
		}
	})

	t.Run("extracting a proto from a POST request", func(t *testing.T) {
		r, err := http.NewRequest("POST", "http://localhost/my/path", bytes.NewReader(protoBytes))
		if err != nil {
			t.Fatalf("Failed to created http.Request: %v", err)
		}

		extByte, err := requestPayload(r)
		if err != nil {
			t.Fatalf("Failed to extract payload from request: %v", err)
		}
		if string(extByte) != string(protoBytes) {
			t.Fatalf("Expected %v and %v to match", string(extByte), string(protoBytes))
		}
	})

	t.Run("extracting JSON from a POST request", func(t *testing.T) {
		r, err := http.NewRequest("POST", "http://localhost/my/path", bytes.NewReader(jsonBytes))
		if err != nil {
			t.Fatalf("Failed to created http.Request: %v", err)
		}

		extByte, err := requestPayload(r)
		if err != nil {
			t.Fatalf("Failed to extract payload from request: %v", err)
		}
		if string(extByte) != string(jsonBytes) {
			t.Fatalf("Expected %v and %v to match", string(extByte), string(jsonBytes))
		}
	})

	t.Run("extracting params from a GET request", func(t *testing.T) {

		r, err := http.NewRequest("GET", "http://localhost/my/path", nil)
		if err != nil {
			t.Fatalf("Failed to created http.Request: %v", err)
		}

		q := r.URL.Query()
		q.Add("name", "Test")
		r.URL.RawQuery = q.Encode()

		extByte, err := requestPayload(r)
		if err != nil {
			t.Fatalf("Failed to extract payload from request: %v", err)
		}
		if string(extByte) != string(jsonBytes) {
			t.Fatalf("Expected %v and %v to match", string(extByte), string(jsonBytes))
		}
	})

	t.Run("GET request with no params", func(t *testing.T) {

		r, err := http.NewRequest("GET", "http://localhost/my/path", nil)
		if err != nil {
			t.Fatalf("Failed to created http.Request: %v", err)
		}

		extByte, err := requestPayload(r)
		if err != nil {
			t.Fatalf("Failed to extract payload from request: %v", err)
		}
		if string(extByte) != "" {
			t.Fatalf("Expected %v and %v to match", string(extByte), "")
		}
	})
}
