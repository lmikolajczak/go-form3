package form3

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient(t *testing.T, baseURL string) (*Client, func()) {
	t.Helper()
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
