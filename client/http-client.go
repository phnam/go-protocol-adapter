// Package client provides API client implementations for different protocols.
package client

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	sdk "github.com/phnam/go-protocol-adapter"
	"github.com/phnam/go-protocol-adapter/common"
	"github.com/phnam/go-protocol-adapter/request"
)

// RestClient implements the APIClient interface for HTTP protocol communication.
type RestClient[T any] struct {
	// BaseURL is the base URL for all API requests
	BaseURL *url.URL
	// UserAgent is the user agent string sent with requests
	UserAgent string

	// private fields
	// httpClient is the underlying HTTP client used for requests
	httpClient *http.Client
	// maxRetryTime is the maximum number of retry attempts for failed requests
	maxRetryTime int
	// waitTime is the duration to wait between retry attempts (in milliseconds)
	waitTime time.Duration
	// timeOut is the request timeout duration (in milliseconds)
	timeOut time.Duration
	// errorLogOnly when true, only logs errors and not successful requests
	errorLogOnly bool
	// logExpiration defines how long logs should be kept
	logExpiration *time.Duration

	// debug enables debug logging when true
	debug bool
	// acceptHttpError when true, treats HTTP error codes as valid responses
	acceptHttpError bool
}

// RequestLogEntry represents a log entry for an API request with all relevant information.
type RequestLogEntry struct {
	// Status indicates the overall status of the request (SUCCESS/FAILED)
	Status string `json:"status,omitempty" bson:"status,omitempty"`
	// ReqURL is the full URL of the request
	ReqURL string `json:"reqUrl,omitempty" bson:"req_url,omitempty"`
	// ReqMethod is the HTTP method used for the request
	ReqMethod string `json:"reqMethod,omitempty" bson:"req_method,omitempty"`
	// Caller identifies the client making the request
	Caller string `json:"caller,omitempty" bson:"caller,omitempty"`
	// ReqHeader contains the request headers
	ReqHeader *map[string]string `json:"reqHeader,omitempty" bson:"req_header,omitempty"`
	// ReqFormData contains form data sent with the request
	ReqFormData *map[string]string `json:"reqFormData,omitempty" bson:"req_form_data,omitempty"`
	// ReqBody contains the request body
	ReqBody *interface{} `json:"reqBody,omitempty" bson:"req_body,omitempty"`

	// TotalTime is the total time taken for the request in milliseconds
	TotalTime int64 `json:"totalTime,omitempty" bson:"total_time,omitempty"`
	// RetryCount is the number of retry attempts made
	RetryCount int `json:"retryCount,omitempty" bson:"retry_count,omitempty"`
	// Results contains the results of each attempt
	Results []*CallResult `json:"results,omitempty" bson:"results,omitempty"`
	// ErrorLog contains error messages if any
	ErrorLog *string `json:"errorLog,omitempty" bson:"error_log,omitempty"`
	// Keys contains any associated keys for the request
	Keys *[]string `json:"keys,omitempty" bson:"keys,omitempty"`
	// Date is the timestamp when the request was made
	Date *time.Time `json:"date,omitempty" bson:"date,omitempty"`
}

// CallResult represents the result of a single API call attempt.
type CallResult struct {
	// RespCode is the HTTP response status code
	RespCode int `json:"respCode,omitempty" bson:"resp_code,omitempty"`
	// RespHeader contains the response headers
	RespHeader map[string][]string `json:"respHeader,omitempty" bson:"resp_header,omitempty"`
	// RespBody contains the response body as a string
	RespBody *string `json:"respBody,omitempty" bson:"resp_body,omitempty"`
	// ResponseTime is the time taken for this attempt in milliseconds
	ResponseTime int64 `json:"responseTime,omitempty" bson:"response_time,omitempty"`
	// ErrorLog contains error messages if any
	ErrorLog *string `json:"errorLog,omitempty" bson:"error_log,omitempty"`
}

// RestResult represents the result of a REST API call.
type RestResult struct {
	// Body is the response body as a string
	Body string `json:"body,omitempty" bson:"body,omitempty"`
	// Content is the raw response body as bytes
	Content []byte `json:"content,omitempty" bson:"content,omitempty"`
	// Code is the HTTP status code
	Code int `json:"code,omitempty" bson:"code,omitempty"`
}

// HTTPMethod is a type representing HTTP methods as strings.
type HTTPMethod string

// HTTPMethodEnum defines a struct containing all supported HTTP methods.
type HTTPMethodEnum struct {
	// Get represents the HTTP GET method
	Get HTTPMethod
	// Query represents a custom QUERY method
	Query HTTPMethod
	// Post represents the HTTP POST method
	Post HTTPMethod
	// Put represents the HTTP PUT method
	Put HTTPMethod
	// Patch represents the HTTP PATCH method
	Patch HTTPMethod
	// Head represents the HTTP HEAD method
	Head HTTPMethod
	// Delete represents the HTTP DELETE method
	Delete HTTPMethod
	// Option represents the HTTP OPTION method
	Option HTTPMethod
}

// HTTPMethods is a global variable containing all supported HTTP methods.
// It provides easy access to HTTP method constants throughout the application.
var HTTPMethods = &HTTPMethodEnum{
	Get:    "GET",
	Query:  "QUERY",
	Post:   "POST",
	Put:    "PUT",
	Patch:  "PATCH",
	Head:   "HEAD",
	Delete: "DELETE",
	Option: "OPTION",
}

// NewHTTPClient creates a new HTTP client based on the provided configuration.
// It implements the APIClient interface for HTTP protocol communication.
//
// Parameters:
//   - config: Configuration parameters for the HTTP client
//
// Returns:
//   - An implementation of the APIClient interface
func NewHTTPClient[T any](config *APIClientConfiguration) APIClient[T] {
	var restCl RestClient[T]

	// Ensure the base URL has the http prefix
	baseURL := config.Address
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "http://" + baseURL
	}

	// Parse the URL
	u, err := url.Parse(baseURL)
	if err == nil {
		restCl.BaseURL = u
	}

	// Create transport with TLS configuration that skips certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Initialize the HTTP client with the transport and timeout
	restCl.httpClient = &http.Client{
		Transport: tr,
		Timeout:   config.Timeout,
	}

	// Configure client settings from the provided configuration
	restCl.SetMaxRetryTime(config.MaxRetry)
	restCl.SetWaitTime(config.WaitToRetry)
	restCl.SetTimeout(config.Timeout)
	restCl.debug = false
	restCl.errorLogOnly = config.ErrorLogOnly
	return &restCl
}

// NewRESTClient creates a new instance of RestClient without proxy support.
//
// Parameters:
//   - baseURL: The base URL for API requests
//   - logName: The name for logging purposes
//   - timeout: The request timeout duration
//   - maxRetryTime: Maximum number of retry attempts
//   - waitTime: Duration to wait between retries
//
// Returns:
//   - A pointer to a new RestClient instance
func NewRESTClient[T any](baseURL string, logName string, timeout time.Duration, maxRetryTime int, waitTime time.Duration) *RestClient[T] {
	return NewRESTClientWithProxy[T](baseURL, logName, "", timeout, maxRetryTime, waitTime)
}

// NewRESTClientWithProxy creates a new instance of RestClient with proxy support.
//
// Parameters:
//   - baseURL: The base URL for API requests
//   - logName: The name for logging purposes
//   - proxyUrl: The URL of the proxy server (empty string for no proxy)
//   - timeout: The request timeout duration
//   - maxRetryTime: Maximum number of retry attempts
//   - waitTime: Duration to wait between retries
//
// Returns:
//   - A pointer to a new RestClient instance
func NewRESTClientWithProxy[T any](baseURL string, logName string, proxyUrl string, timeout time.Duration, maxRetryTime int, waitTime time.Duration) *RestClient[T] {

	var restCl RestClient[T]

	// Ensure the base URL has the http prefix
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "http://" + baseURL
	}

	// Parse the URL
	u, err := url.Parse(baseURL)
	if err == nil {
		restCl.BaseURL = u
	}

	// Create transport with TLS configuration that skips certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Configure proxy if provided
	if proxyUrl != "" {
		urlObj, err := url.Parse(proxyUrl)
		if err == nil {
			tr.Proxy = http.ProxyURL(urlObj)
		}
	}

	// Initialize the HTTP client with the transport and timeout
	restCl.httpClient = &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	// Configure client settings
	restCl.SetMaxRetryTime(maxRetryTime)
	restCl.SetWaitTime(waitTime)
	restCl.SetTimeout(timeout)
	restCl.debug = false

	return &restCl
}

// addParams appends query parameters to a base URL.
//
// Parameters:
//   - baseURL: The base URL to append parameters to
//   - params: A map of parameter names to values
//
// Returns:
//   - The URL with query parameters appended
func addParams(baseURL string, params map[string]string) string {
	baseURL += "?"
	p := url.Values{}
	for key, value := range params {
		p.Add(key, value)
	}
	return baseURL + p.Encode()
}

// addResult adds a CallResult to the RequestLogEntry's Results slice.
//
// Parameters:
//   - rs: The CallResult to add
func (entry *RequestLogEntry) addResult(rs *CallResult) {
	entry.Results = append(entry.Results, rs)
}

// SetDebug enables or disables debug logging for the RestClient.
//
// Parameters:
//   - val: true to enable debug logging, false to disable
func (c *RestClient[T]) SetDebug(val bool) {
	c.debug = val
}

// SetTimeout sets the timeout duration for HTTP requests.
//
// Parameters:
//   - timeout: The duration to wait before timing out a request
func (c *RestClient[T]) SetTimeout(timeout time.Duration) {
	c.timeOut = timeout
	c.httpClient.Timeout = timeout
}

// AcceptHTTPError configures whether HTTP error responses should be treated as valid responses.
//
// Parameters:
//   - accept: When true, HTTP error responses will be returned as valid responses
func (c *RestClient[T]) AcceptHTTPError(accept bool) {
	c.acceptHttpError = accept
}

// SetWaitTime sets the duration to wait between retry attempts.
//
// Parameters:
//   - waitTime: The duration to wait between retries
func (c *RestClient[T]) SetWaitTime(waitTime time.Duration) {
	c.waitTime = waitTime
}

// SetMaxRetryTime sets the maximum number of retry attempts for failed requests.
//
// Parameters:
//   - maxRetryTime: The maximum number of retry attempts
func (c *RestClient[T]) SetMaxRetryTime(maxRetryTime int) {
	c.maxRetryTime = maxRetryTime
}

// initRequest creates and initializes an HTTP request with the specified parameters.
//
// Parameters:
//   - method: The HTTP method to use
//   - headers: HTTP headers to include in the request
//   - params: Query parameters to include in the URL
//   - body: The request body (for POST, PUT, etc.)
//   - path: The path to append to the base URL
//   - userAgent: The User-Agent header value
//
// Returns:
//   - A pointer to an http.Request
//   - An error if request creation fails
func (c *RestClient[T]) initRequest(method HTTPMethod, headers map[string]string, params map[string]string, body interface{}, path string, userAgent string) (*http.Request, error) {

	// Construct the full URL by combining base URL and path
	urlStr := c.BaseURL.String()
	if path != "" {
		if strings.HasSuffix(urlStr, "/") || strings.HasPrefix(path, "/") {
			urlStr += path
		} else {
			urlStr += "/" + path
		}
	}

	// Prepare the request body if provided
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	var err error
	var req *http.Request

	// Handle form-encoded POST requests differently
	if method == HTTPMethods.Post && headers != nil && headers["Content-Type"] == "application/x-www-form-urlencoded" && params != nil && len(params) > 0 {
		data := url.Values{}
		for key, val := range params {
			data.Set(key, val)
		}
		req, err = http.NewRequest(string(method), urlStr, strings.NewReader(data.Encode()))
	} else {
		// For other requests, add params to the URL
		urlStr = addParams(urlStr, params)
		req, err = http.NewRequest(string(method), urlStr, buf)
	}

	if err != nil {
		return nil, err
	}

	// Set common headers
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", userAgent)

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

// MakeHTTPRequest makes an HTTP request with the specified parameters.
// This is a convenience wrapper around MakeHTTPRequestWithKey without keys.
//
// Parameters:
//   - method: The HTTP method to use
//   - headers: HTTP headers to include in the request
//   - params: Query parameters to include in the URL
//   - body: The request body (for POST, PUT, etc.)
//   - path: The path to append to the base URL
//
// Returns:
//   - A pointer to a RestResult containing the response
//   - An error if the request fails
func (c *RestClient[T]) MakeHTTPRequest(method HTTPMethod, headers map[string]string, params map[string]string, body interface{}, path string) (*RestResult, error) {
	return c.MakeHTTPRequestWithKey(method, headers, params, body, path, nil)
}

// writeLog writes a request log entry to the console.
// If errorLogOnly is true, it only logs entries with a status other than "SUCCESS".
//
// Parameters:
//   - logEntry: The RequestLogEntry to log
func (c *RestClient[T]) writeLog(logEntry *RequestLogEntry) {

	if c.debug {
		fmt.Println(" +++ Writing log ...")
	}

	// Only log errors if errorLogOnly is true
	if logEntry.Status != "SUCCESS" || !c.errorLogOnly {
		str, err := json.Marshal(logEntry)
		if err != nil {
			fmt.Println("Error when marshal log entry")
		} else {
			fmt.Println(string(str))
		}
	}
}

// MakeHTTPRequestWithKey makes an HTTP request with the specified parameters and associated keys.
// It handles retries, logging, and response processing.
//
// Parameters:
//   - method: The HTTP method to use
//   - headers: HTTP headers to include in the request
//   - params: Query parameters to include in the URL
//   - body: The request body (for POST, PUT, etc.)
//   - path: The path to append to the base URL
//   - keys: Optional keys associated with this request for tracking/logging
//
// Returns:
//   - A pointer to a RestResult containing the response
//   - An error if the request fails after all retry attempts
func (c *RestClient[T]) MakeHTTPRequestWithKey(method HTTPMethod, headers map[string]string, params map[string]string, body interface{}, path string, keys *[]string) (*RestResult, error) {

	date := time.Now()
	// init log
	userAgent := "go-protocol-adapter"
	hostname, err := os.Hostname()
	if err == nil {
		userAgent += " " + hostname + "/" + os.Getenv("env")
	}
	logEntry := &RequestLogEntry{
		ReqURL:      c.BaseURL.String() + path,
		ReqMethod:   string(method),
		ReqFormData: &params,
		ReqHeader:   &headers,
		ReqBody:     &body,
		Keys:        keys,
		Date:        &date,
		Caller:      userAgent,
	}

	if c.debug {
		fmt.Println(" +++ Try to init request ...")
	}

	canRetryCount := c.maxRetryTime

	tstart := time.Now().UnixNano() / 1e6

	for canRetryCount >= 0 {

		req, reqErr := c.initRequest(method, headers, params, body, path, userAgent)

		if c.debug {
			fmt.Println(" +++ Request inited.")
		}

		if reqErr != nil {
			msg := reqErr.Error()
			logEntry.ErrorLog = &msg
			if c.debug {
				fmt.Println("Error when init request: " + msg)
			}
			return nil, reqErr
		}
		// start time
		startCallTime := time.Now().UnixNano() / 1e6
		if c.debug {
			fmt.Println("+++ Let call: " + logEntry.ReqMethod + " " + logEntry.ReqURL)
		}

		// add call result
		callRs := &CallResult{}

		// do request
		resp, err := c.httpClient.Do(req)
		if c.debug {
			fmt.Println("+++ HTTP call ended!")
		}

		// make request successful
		if err == nil {
			restResult, err := c.readBody(resp, callRs, logEntry, canRetryCount, startCallTime, tstart)
			if restResult != nil {
				logEntry.Status = "SUCCESS"
				return restResult, err
			}

			if c.acceptHttpError {
				logEntry.Status = "FAILED"
				return restResult, err
			}
		} else {
			if c.debug {
				fmt.Println("HTTP Error: " + err.Error())
			}
			msg := err.Error()
			callRs.ErrorLog = &msg
		}

		tend := time.Now().UnixNano() / 1e6
		callRs.ResponseTime = tend - startCallTime

		canRetryCount--

		if canRetryCount >= 0 {
			time.Sleep(c.waitTime)
			if c.debug {
				fmt.Println("Comeback from sleep ...")
			}
		}

		if c.debug {
			fmt.Println("Count down ...")
		}
		if canRetryCount >= 0 {
			logEntry.RetryCount = c.maxRetryTime - canRetryCount
		}
		logEntry.addResult(callRs)
		if c.debug {
			fmt.Println("Try to exit loop ...")
		}
	}

	if c.debug {
		fmt.Println("Exit retry loop.")
	}

	tend := time.Now().UnixNano() / 1e6
	logEntry.TotalTime = tend - tstart
	logEntry.Status = "FAILED"
	return nil, errors.New("fail to call endpoint API " + logEntry.ReqURL)
}

// readBody reads and processes the HTTP response body.
// It handles gzip decompression and updates the call result and log entry.
//
// Parameters:
//   - resp: The HTTP response
//   - callRs: The call result to update
//   - logEntry: The request log entry to update
//   - canRetryCount: The number of remaining retry attempts
//   - startCallTime: The timestamp when the call started (in milliseconds)
//   - tstart: The timestamp when the entire request started (in milliseconds)
//
// Returns:
//   - A pointer to a RestResult containing the response
//   - An error if processing fails
func (c *RestClient[T]) readBody(resp *http.Response, callRs *CallResult, logEntry *RequestLogEntry, canRetryCount int, startCallTime int64, tstart int64) (*RestResult, error) {
	defer resp.Body.Close()
	v, err := io.ReadAll(resp.Body)
	if err != nil {
		msg := err.Error()
		callRs.ErrorLog = &msg
		return nil, err
	}

	if c.debug {
		fmt.Println("+++ IO read ended!")
	}
	restResult := RestResult{
		Code:    resp.StatusCode,
		Body:    string(v),
		Content: v,
	}

	encoding := resp.Header.Get("Content-Encoding")
	if encoding == "gzip" {
		if c.debug {
			fmt.Println("+++ Start to gunzip")
		}
		gr, _ := gzip.NewReader(bytes.NewBuffer(restResult.Content))
		data, err := io.ReadAll(gr)
		gr.Close()
		if err != nil {
			return nil, err
		}
		if c.debug {
			fmt.Println("+++ gunzip successfully")
		}
		restResult.Content = data
		restResult.Body = string(data)
	}

	// set call result
	callRs.RespCode = restResult.Code
	callRs.RespBody = &restResult.Body
	if resp.Header != nil {
		h := (map[string][]string)(resp.Header)
		if h != nil {
			callRs.RespHeader = map[string][]string{}
			for k, v := range h {
				if strings.HasPrefix(k, "X-") {
					callRs.RespHeader[k] = v
				}
			}
		}
	}

	if c.debug {
		fmt.Println("+++ Read data end, http code: " + string(resp.StatusCode))
	}
	if c.acceptHttpError || (resp.StatusCode >= 200 && resp.StatusCode < 300) || (resp.StatusCode >= 400 && resp.StatusCode < 500) {
		// add log
		tend := time.Now().UnixNano() / 1e6
		callRs.ResponseTime = tend - startCallTime
		logEntry.TotalTime = tend - tstart
		if canRetryCount >= 0 {
			logEntry.RetryCount = c.maxRetryTime - canRetryCount
		}
		//sample
		logEntry.addResult(callRs)
		//return
		return &restResult, err
	}
	return nil, nil
}

// MakeRequest implements the APIClient interface method for making API requests.
// It converts the generic APIRequest to an HTTP request and processes the response.
//
// Parameters:
//   - req: The API request to process
//
// Returns:
//   - A pointer to a common.APIResponse containing the response
func (c *RestClient[T]) MakeRequest(req request.APIRequest) *common.APIResponse[T] {
	var data interface{}
	var reqMethod = req.GetMethod()
	var method HTTPMethod

	switch reqMethod.Value {
	case "GET":
		method = HTTPMethods.Get
	case "PUT":
		method = HTTPMethods.Put
		req.ParseBody(&data)
	case "POST":
		method = HTTPMethods.Post
		req.ParseBody(&data)
	case "PATCH":
		method = HTTPMethods.Put
		req.ParseBody(&data)
	case "DELETE":
		method = HTTPMethods.Delete
	case "OPTIONS":
		method = HTTPMethods.Option
	}

	if c.debug {
		fmt.Println("Req info: " + reqMethod.Value + " / " + req.GetPath())
		if data != nil {
			fmt.Println("Data not null")
		}
	}

	result, err := c.MakeHTTPRequest(method, req.GetHeaders(), req.GetParams(), data, req.GetPath())

	if err != nil {
		return &common.APIResponse[T]{
			Status:  common.APIStatus.Error,
			Message: "HTTP Endpoint Error: " + err.Error(),
		}
	}

	var resp = &common.APIResponse[T]{}
	err = json.Unmarshal(result.Content, &resp)

	if resp.Data != nil {
		jsonStr, err := json.Marshal(resp.Data)
		if err == nil {
			resp.Data = sdk.ConvertToObjectSlice[T](string(jsonStr))
		}
	}

	if resp.Status == "" {
		if result.Code >= 500 {
			resp.Status = common.APIStatus.Error
		} else if result.Code >= 400 {
			if result.Code == 404 {
				resp.Status = common.APIStatus.NotFound
			} else if result.Code == 403 {
				resp.Status = common.APIStatus.Forbidden
			} else if result.Code == 401 {
				resp.Status = common.APIStatus.Unauthorized
			} else {
				resp.Status = common.APIStatus.Invalid
			}
		} else {
			resp.Status = common.APIStatus.Ok
		}
	}

	if err != nil {
		return &common.APIResponse[T]{
			Status:  common.APIStatus.Error,
			Message: "Response Data Error: " + err.Error() + " body=" + result.Body,
		}
	}
	return resp
}
