package curl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type HTTPMethod int

const (
	MethodGet = iota + 1
	MethodPost
	MethodPut
	MethodDelete
	MethodPatch
)

// Curl encapsulates the details required to perform HTTP requests.
type Curl struct {
	client    *http.Client
	transport http.RoundTripper
	ctx       context.Context
	url       string
	headers   http.Header
	method    string
	buff      *bytes.Buffer
	err       error

	// basic auth
	isBasicAuth bool
	username    string
	password    string
}

// NewCurl creates a new Curl instance with the specified URL and HTTP method.
//
// This function initializes a new Curl instance, by default sets the context to the background,
// and configures the URL, method, HTTP client, transport, and buffer. It returns the
// configured Curl instance, ready to be used for making HTTP requests.
func NewCurl(url string, method HTTPMethod) *Curl {
	c := new(Curl)

	return c.WithContext(context.Background()).
		SetUrl(url).
		SetMethod(method).
		SetHTTPClient(http.DefaultClient).
		SetHTTPTransport(http.DefaultTransport).
		SetupBuffer()
}

// SetHTTPClient sets the HTTP client for the Curl instance.
func (c *Curl) SetHTTPClient(client *http.Client) *Curl {
	c.client = client

	return c
}

// SetHTTPTransport sets the HTTP transport for the Curl instance.
func (c *Curl) SetHTTPTransport(t http.RoundTripper) *Curl {
	c.transport = t

	return c
}

// SetupBuffer initializes the internal buffer for the Curl instance if it is not already initialized.
func (c *Curl) SetupBuffer() *Curl {
	if c.buff == nil {
		c.buff = new(bytes.Buffer)
	}

	return c
}

// SetUrl sets the URL for the Curl instance.
func (c *Curl) SetUrl(url string) *Curl {
	c.url = url

	return c
}

// SetMethod sets the HTTP method for the Curl instance.
//
// This method sets the provided HTTP method to the Curl instance,
// converting it to its corresponding canonical form (e.g., "GET" to http.MethodGet).
// If the provided method is not recognized, it sets an error indicating an invalid HTTP method.
func (c *Curl) SetMethod(method HTTPMethod) *Curl {
	switch method {
	case MethodGet:
		c.method = http.MethodGet
	case MethodPost:
		c.method = http.MethodPost
	case MethodPut:
		c.method = http.MethodPut
	case MethodDelete:
		c.method = http.MethodDelete
	case MethodPatch:
		c.method = http.MethodPatch
	default:
		c.err = errors.New("invalid http method")
	}

	return c
}

// WithContext sets the context for the Curl instance.
func (c *Curl) WithContext(ctx context.Context) *Curl {
	c.ctx = ctx
	return c
}

// SetHeader sets a header key-value pair for the Curl instance.
//
// This method adds or updates the header with the provided key and value
// in the headers map of the Curl instance. If the headers map is nil,
// it initializes it before setting the header.
func (c *Curl) SetHeader(k, v string) *Curl {
	if c.headers == nil {
		c.headers = make(http.Header)
	}

	c.headers.Set(k, v)

	return c
}

func (c *Curl) SetBasicAuth(username, password string) *Curl {
	c.username = username
	c.password = password
	c.isBasicAuth = true

	return c
}

// Body sets the request body of the Curl instance by encoding the provided payload to JSON.
func (c *Curl) Body(payload any) *Curl {
	c.err = writeJSON(c.buff, payload)
	return c
}

// writeJSON encodes a value v into JSON format and writes it to the provided io.Writer.
func writeJSON(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

// Exec executes the HTTP request configured in the Curl instance.
//
// This method constructs an HTTP request based on the configured method, URL,
// context, headers, and request body (if any), then sends the request using
// the configured HTTP client and returns the HTTP response received. If any errors
// occur during request construction, sending, or receiving the response, they are
// returned. If there was an error previously set in the Curl instance, that error
// is returned immediately.
func (c *Curl) Exec() (*http.Response, error) {
	if c.err != nil {
		return nil, c.err
	}

	req, err := http.NewRequestWithContext(c.ctx, c.method, c.url, c.buff)
	if err != nil {
		return nil, err
	}

	req.Header = c.headers

	if c.isBasicAuth {
		req.SetBasicAuth(c.username, c.password)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
