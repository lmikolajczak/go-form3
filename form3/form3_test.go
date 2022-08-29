package form3_test

import (
	"github.com/lmikolajczak/go-form3/form3"
	"io"
	"net/http"
	"testing"
)

func TestClient_NewRequest(t *testing.T) {
	f3, _, teardown := form3.TestClientWithServer(t)
	defer teardown()

	type payload struct {
		Example string `json:"example"`
	}

	urlIn, urlOut := "/v1/api/endpoint", f3.BaseURL()+"/v1/api/endpoint"
	payloadIn, payloadOut := &payload{Example: "test example"}, `{"example":"test example"}`

	request, err := f3.NewRequest(http.MethodPost, urlIn, payloadIn)
	if err != nil {
		t.Fatalf("err = %s; want nil", err)
	}

	// Check if proper method, endpoint and payload are set on the request.
	testMethod(t, request, http.MethodPost)
	testEndpoint(t, request, urlOut)
	testBody(t, request, payloadOut)
}

func TestClient_Request(t *testing.T) {
	testcases := []struct {
		name        string
		wantHeaders map[string]string
	}{
		{
			name:        "without headers",
			wantHeaders: map[string]string{},
		},
		{
			name: "with headers",
			wantHeaders: map[string]string{
				"Accept":       "application/vnd.api+json",
				"Content-Type": "application/vnd.api+json",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			f3, mux, teardown := form3.TestClientWithServer(t)
			defer teardown()

			mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
				testHeaders(t, r, tc.wantHeaders)
			})

			request, err := f3.NewRequest(http.MethodPost, "/test", nil)
			if err != nil {
				t.Fatalf("err = %v; want: nil", err)
			}

			err = f3.Request(nil, request, tc.wantHeaders)
			if err != nil {
				t.Fatalf("err = %v; want: nil", err)
			}
		})
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("method = %s; want: %s", got, want)
	}
}

func testEndpoint(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.URL.String(); got != want {
		t.Errorf("url = %s; want: %s", got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("error reading request body: %v", err)
	}
	if got := string(body); got != want {
		t.Errorf("body: %s; want %s", got, want)
	}
}

func testHeaders(t *testing.T, r *http.Request, want map[string]string) {
	t.Helper()
	for key, value := range want {
		if got := r.Header.Get(key); got != value {
			t.Errorf("headers[%s] = %s; want: %s", key, got, value)
		}
	}
}
