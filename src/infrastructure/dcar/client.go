// Package dcar provides primitives to interact with the openapi HTTP API.
package dcar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// HttpRequestDoer performs HTTP requests.
//
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

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// NewClient Creates a new Client, with reasonable defaults
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

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// ClientInterface The interface specification for the client above.
type ClientInterface interface {
	// GetCars request
	GetCars(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AddVehicleWithBody AddVehicle request with any body
	AddVehicleWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	AddVehicle(ctx context.Context, body AddVehicleJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteCar request
	DeleteCar(ctx context.Context, vin VinParam, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetCar request
	GetCar(ctx context.Context, vin VinParam, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetCars(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetCarsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AddVehicleWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAddVehicleRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AddVehicle(ctx context.Context, body AddVehicleJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAddVehicleRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteCar(ctx context.Context, vin VinParam, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteCarRequest(c.Server, vin)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetCar(ctx context.Context, vin VinParam, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetCarRequest(c.Server, vin)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetCarsRequest generates requests for GetCars
func NewGetCarsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/cars")
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

// NewAddVehicleRequest calls the generic AddVehicle builder with application/json body
func NewAddVehicleRequest(server string, body AddVehicleJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewAddVehicleRequestWithBody(server, "application/json", bodyReader)
}

// NewAddVehicleRequestWithBody generates requests for AddVehicle with any type of body
func NewAddVehicleRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/cars")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteCarRequest generates requests for DeleteCar
func NewDeleteCarRequest(server string, vin VinParam) (*http.Request, error) {
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

	operationPath := fmt.Sprintf("/cars/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetCarRequest generates requests for GetCar
func NewGetCarRequest(server string, vin VinParam) (*http.Request, error) {
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

	operationPath := fmt.Sprintf("/cars/%s", pathParam0)
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

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
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

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetCarsWithResponse GetCars request
	GetCarsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCarsResponse, error)

	// AddVehicleWithBodyWithResponse AddVehicle request with any body
	AddVehicleWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AddVehicleResponse, error)

	AddVehicleWithResponse(ctx context.Context, body AddVehicleJSONRequestBody, reqEditors ...RequestEditorFn) (*AddVehicleResponse, error)

	// DeleteCarWithResponse DeleteCar request
	DeleteCarWithResponse(ctx context.Context, vin VinParam, reqEditors ...RequestEditorFn) (*DeleteCarResponse, error)

	// GetCarWithResponse GetCar request
	GetCarWithResponse(ctx context.Context, vin VinParam, reqEditors ...RequestEditorFn) (*GetCarResponse, error)
}

type GetCarsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Vin
}

// Status returns HTTPResponse.Status
func (r GetCarsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetCarsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AddVehicleResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Vin
}

// Status returns HTTPResponse.Status
func (r AddVehicleResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AddVehicleResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteCarResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeleteCarResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteCarResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetCarResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Car
}

// Status returns HTTPResponse.Status
func (r GetCarResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetCarResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetCarsWithResponse request returning *GetCarsResponse
func (c *ClientWithResponses) GetCarsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCarsResponse, error) {
	rsp, err := c.GetCars(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetCarsResponse(rsp)
}

// AddVehicleWithBodyWithResponse request with arbitrary body returning *AddVehicleResponse
func (c *ClientWithResponses) AddVehicleWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AddVehicleResponse, error) {
	rsp, err := c.AddVehicleWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAddVehicleResponse(rsp)
}

func (c *ClientWithResponses) AddVehicleWithResponse(ctx context.Context, body AddVehicleJSONRequestBody, reqEditors ...RequestEditorFn) (*AddVehicleResponse, error) {
	rsp, err := c.AddVehicle(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAddVehicleResponse(rsp)
}

// DeleteCarWithResponse request returning *DeleteCarResponse
func (c *ClientWithResponses) DeleteCarWithResponse(ctx context.Context, vin VinParam, reqEditors ...RequestEditorFn) (*DeleteCarResponse, error) {
	rsp, err := c.DeleteCar(ctx, vin, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteCarResponse(rsp)
}

// GetCarWithResponse request returning *GetCarResponse
func (c *ClientWithResponses) GetCarWithResponse(ctx context.Context, vin VinParam, reqEditors ...RequestEditorFn) (*GetCarResponse, error) {
	rsp, err := c.GetCar(ctx, vin, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetCarResponse(rsp)
}

// ParseGetCarsResponse parses an HTTP response from a GetCarsWithResponse call
func ParseGetCarsResponse(rsp *http.Response) (*GetCarsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetCarsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Vin
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseAddVehicleResponse parses an HTTP response from a AddVehicleWithResponse call
func ParseAddVehicleResponse(rsp *http.Response) (*AddVehicleResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AddVehicleResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest Vin
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseDeleteCarResponse parses an HTTP response from a DeleteCarWithResponse call
func ParseDeleteCarResponse(rsp *http.Response) (*DeleteCarResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteCarResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetCarResponse parses an HTTP response from a GetCarWithResponse call
func ParseGetCarResponse(rsp *http.Response) (*GetCarResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetCarResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Car
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}
