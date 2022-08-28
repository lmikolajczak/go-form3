package form3

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestClient(t *testing.T) (*Client, func()) {
	t.Helper()
	baseURL := os.Getenv("FORM3_API_BASE_URL")
	if len(baseURL) == 0 {
		baseURL = "http://localhost:8080"
	}
	return NewClient(baseURL), func() {}
}

func TestClientWithServer(t *testing.T) (*Client, *http.ServeMux, func()) {
	t.Helper()
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	c := NewClient(server.URL)
	return c, mux, func() {
		server.Close()
	}
}
