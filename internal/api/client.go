package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
)

// UserAgent is the User-Agent header value sent with all requests.
var UserAgent = "planitaicojp/freee-cli/dev"

const (
	baseURL        = "https://api.freee.co.jp"
	defaultTimeout = 30 * time.Second
	maxRetries     = 3
)

// Client is the HTTP client for freee API.
type Client struct {
	HTTP            *http.Client
	Token           string
	CompanyID       int64
	baseURLOverride string // for testing only
}

// NewClient creates a new API client.
func NewClient(token string, companyID int64) *Client {
	return &Client{
		HTTP:      &http.Client{Timeout: defaultTimeout},
		Token:     token,
		CompanyID: companyID,
	}
}

// BaseURL returns the freee API base URL.
func (c *Client) BaseURL() string {
	if c.baseURLOverride != "" {
		return c.baseURLOverride
	}
	return baseURL
}

// SetBaseURL overrides the base URL (for testing).
func (c *Client) SetBaseURL(url string) {
	c.baseURLOverride = url
}

// Do executes an HTTP request with auth headers and error handling.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", UserAgent)
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	if req.Header.Get("Content-Type") == "" && req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	var reqBody []byte
	if debugLevel >= DebugAPI && req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(reqBody))
	}
	debugLogRequest(req, reqBody)

	start := time.Now()
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err = c.HTTP.Do(req)
		if err != nil {
			if attempt == maxRetries {
				return nil, &cerrors.NetworkError{Err: err}
			}
			if req.Body != nil {
				return nil, &cerrors.NetworkError{Err: err}
			}
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			if attempt < maxRetries && req.Body == nil {
				wait := retryAfterDuration(resp, attempt)
				resp.Body.Close()
				time.Sleep(wait)
				continue
			}
		}
		break
	}
	elapsed := time.Since(start)

	if resp == nil {
		return nil, &cerrors.NetworkError{Err: fmt.Errorf("no response after retries")}
	}

	if debugLevel >= DebugAPI {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
		debugLogResponse(resp, elapsed, respBody)
	} else {
		debugLogResponse(resp, elapsed, nil)
	}

	if resp.StatusCode >= 400 {
		return resp, parseAPIError(resp)
	}

	return resp, nil
}

// Request creates and executes a request.
func (c *Client) Request(method, url string, body any) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	return c.Do(req)
}

// Get performs a GET request and decodes the response.
func (c *Client) Get(url string, result any) error {
	resp, err := c.Request(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

// Post performs a POST request and decodes the response.
func (c *Client) Post(url string, body, result any) (*http.Response, error) {
	resp, err := c.Request(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

// Put performs a PUT request.
func (c *Client) Put(url string, body, result any) error {
	resp, err := c.Request(http.MethodPut, url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

// Delete performs a DELETE request.
func (c *Client) Delete(url string) error {
	resp, err := c.Request(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// retryAfterDuration parses the Retry-After header and returns the wait duration.
// Falls back to exponential backoff if the header is missing or unparseable.
func retryAfterDuration(resp *http.Response, attempt int) time.Duration {
	if ra := resp.Header.Get("Retry-After"); ra != "" {
		if secs, err := strconv.Atoi(ra); err == nil && secs > 0 {
			return time.Duration(secs) * time.Second
		}
	}
	return time.Duration(attempt+1) * time.Second
}

// parseAPIError reads the response body and returns an appropriate error.
func parseAPIError(resp *http.Response) error {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	apiErr := &cerrors.APIError{
		StatusCode: resp.StatusCode,
		Message:    string(body),
	}

	// freee API error format: {"status_code":400,"errors":[{"type":"...","messages":["..."]}]}
	var errResp struct {
		StatusCode int `json:"status_code"`
		Errors     []struct {
			Type     string   `json:"type"`
			Messages []string `json:"messages"`
		} `json:"errors"`
	}
	if json.Unmarshal(body, &errResp) == nil && len(errResp.Errors) > 0 {
		var msgs []string
		for _, e := range errResp.Errors {
			msgs = append(msgs, e.Messages...)
		}
		if len(msgs) > 0 {
			apiErr.Message = fmt.Sprintf("%s", msgs)
		}
		if len(errResp.Errors) > 0 {
			apiErr.Code = errResp.Errors[0].Type
		}
	}

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return &cerrors.AuthError{Message: apiErr.Message}
	}
	if resp.StatusCode == 404 {
		return &cerrors.NotFoundError{Resource: "resource", ID: ""}
	}

	return apiErr
}
