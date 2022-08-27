// Package form3 is a simple Form3 REST API client.
//
// Currently, it implements Fetch, Create and Delete actions on the Account resource.
// For more info check: https://api-docs.form3.tech/api.html.
package form3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient interface allows to plug in and use custom HTTP clients.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ClientOption represents an option that can be used to configure client.
type ClientOption func(*Client)

// WithHTTPClient allows to set a custom HTTP client.
func WithHTTPClient(httpClient HTTPClient) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// Client represents Form3 REST API client
type Client struct {
	baseURL    string
	httpClient HTTPClient
}

// NewClient returns a new Form3 REST API client.
func NewClient(baseURL string, options ...ClientOption) *Client {
	c := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, option := range options {
		option(c)
	}

	return c
}

// BaseURL returns base URL configured on the Form3 client.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// NewRequest returns new http request with given method, endpoint and payload.
func (c *Client) NewRequest(method, endpoint string, payload interface{}) (*http.Request, error) {
	switch payload.(type) {
	case nil:
		return http.NewRequest(method, c.BaseURL()+endpoint, http.NoBody)
	default:
		body, err := c.marshal(payload)
		if err != nil {
			return nil, err
		}
		return http.NewRequest(method, c.BaseURL()+endpoint, bytes.NewBuffer(body))
	}
}

// F3Error represents an error returned from Form3 REST API.
type F3Error struct {
	StatusCode   int    `json:"-"`
	ErrorCode    int    `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// Error returns a string representation of the F3Error.
func (e F3Error) Error() string {
	return fmt.Sprintf("http %d: code: %d, message=%s", e.StatusCode, e.ErrorCode, e.ErrorMessage)
}

// Request makes a http request to the Form3 REST API.
func (c *Client) Request(v interface{}, request *http.Request, headers map[string]string) error {
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated:
		if err = c.unmarshal(body, &v); err != nil {
			return err
		}
		return nil
	case http.StatusNoContent:
		return nil
	default:
		f3Error := F3Error{StatusCode: response.StatusCode}
		if err = c.unmarshal(body, &f3Error); err != nil {
			return err
		}
		return &f3Error
	}
}

// unmarshal parses JSON-encoded data and stores the result in the value pointed by v.
func (c *Client) unmarshal(body []byte, v interface{}) error {
	if len(body) > 0 {
		if err := json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("json %s, error: %v", string(body), err)
		}
	}
	return nil
}

// marshal returns the JSON-encoding of v.
func (c *Client) marshal(v interface{}) ([]byte, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("json %s, error: %v", string(body), err)
	}
	return body, nil
}

// String returns a pointer to the given string v.
func String(v string) *string { return &v }
