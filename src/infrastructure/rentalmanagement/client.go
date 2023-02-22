// Package rentalManagement provides primitives to interact with the openapi HTTP API.
package rentalManagement

//go:generate mockgen -source=client.go -package=rentalmanagementmocks -destination=../../mocks/rentalmanagementmocks/client_mock.go

import (
	"PFleetManagement/logic/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// HttpRequestDoer performs HTTP requests.
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// NewClient creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// ClientInterface specifies an interface for the client above.
type ClientInterface interface {
	// GetNextRental request
	GetNextRental(ctx context.Context, vin model.VinParam) (*http.Response, error)
}

func (c *Client) GetNextRental(ctx context.Context, vin model.VinParam) (*http.Response, error) {
	req, err := NewGetNextRentalRequest(c.Server, vin)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return c.Client.Do(req)
}

// NewGetNextRentalRequest generates requests for GetNextRental
func NewGetNextRentalRequest(server string, vin model.VinParam) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "vin", runtime.ParamLocationPath, vin)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/cars/%s/rentalStatus", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetNextRentalWithResponse request
	GetNextRentalWithResponse(ctx context.Context, vin model.VinParam) (*GetNextRentalResponse, error)
}

type GetNextRentalResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *model.Rental
}

// Status returns HTTPResponse.Status
func (r GetNextRentalResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetNextRentalResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetNextRentalWithResponse request returning *GetNextRentalResponse
func (c *ClientWithResponses) GetNextRentalWithResponse(ctx context.Context, vin model.VinParam) (*GetNextRentalResponse, error) {
	rsp, err := c.GetNextRental(ctx, vin)
	if err != nil {
		return nil, err
	}
	return ParseGetNextRentalResponse(rsp)
}

// ParseGetNextRentalResponse parses an HTTP response from a GetNextRentalWithResponse call
func ParseGetNextRentalResponse(rsp *http.Response) (*GetNextRentalResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetNextRentalResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest model.Rental
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}
