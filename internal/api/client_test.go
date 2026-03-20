package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
)

func TestNewClient(t *testing.T) {
	c := NewClient("test-token", 12345)
	if c.Token != "test-token" {
		t.Errorf("Token = %q, want %q", c.Token, "test-token")
	}
	if c.CompanyID != 12345 {
		t.Errorf("CompanyID = %d, want %d", c.CompanyID, 12345)
	}
	if c.HTTP == nil {
		t.Error("HTTP client is nil")
	}
}

func TestClient_Do_SetsHeaders(t *testing.T) {
	var gotHeaders http.Header
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header
		w.WriteHeader(200)
	}))
	defer ts.Close()

	c := NewClient("my-token", 1)
	c.HTTP = ts.Client()

	req, _ := http.NewRequest("GET", ts.URL+"/test", nil)
	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Do error: %v", err)
	}
	resp.Body.Close()

	if got := gotHeaders.Get("Authorization"); got != "Bearer my-token" {
		t.Errorf("Authorization = %q, want %q", got, "Bearer my-token")
	}
	if got := gotHeaders.Get("User-Agent"); got != UserAgent {
		t.Errorf("User-Agent = %q, want %q", got, UserAgent)
	}
	if got := gotHeaders.Get("Accept"); got != "application/json" {
		t.Errorf("Accept = %q, want %q", got, "application/json")
	}
}

func TestClient_Do_NoTokenSkipsAuth(t *testing.T) {
	var gotHeaders http.Header
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header
		w.WriteHeader(200)
	}))
	defer ts.Close()

	c := NewClient("", 1)
	c.HTTP = ts.Client()

	req, _ := http.NewRequest("GET", ts.URL+"/test", nil)
	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Do error: %v", err)
	}
	resp.Body.Close()

	if got := gotHeaders.Get("Authorization"); got != "" {
		t.Errorf("Authorization should be empty, got %q", got)
	}
}

func TestClient_Get_DecodesJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"name": "test"}) //nolint:errcheck
	}))
	defer ts.Close()

	c := NewClient("tok", 1)
	c.HTTP = ts.Client()

	var result map[string]string
	if err := c.Get(ts.URL+"/api", &result); err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result["name"] != "test" {
		t.Errorf("name = %q, want %q", result["name"], "test")
	}
}

func TestClient_Get_NilResult(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer ts.Close()

	c := NewClient("tok", 1)
	c.HTTP = ts.Client()

	if err := c.Get(ts.URL+"/api", nil); err != nil {
		t.Fatalf("Get error: %v", err)
	}
}

func TestClient_Post(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "received": body}) //nolint:errcheck
	}))
	defer ts.Close()

	c := NewClient("tok", 1)
	c.HTTP = ts.Client()

	var result map[string]any
	resp, err := c.Post(ts.URL+"/api", map[string]string{"key": "val"}, &result)
	if err != nil {
		t.Fatalf("Post error: %v", err)
	}
	resp.Body.Close()

	if result["id"] != float64(1) {
		t.Errorf("id = %v, want 1", result["id"])
	}
}

func TestClient_Put(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "updated"}) //nolint:errcheck
	}))
	defer ts.Close()

	c := NewClient("tok", 1)
	c.HTTP = ts.Client()

	var result map[string]string
	if err := c.Put(ts.URL+"/api", map[string]string{"a": "b"}, &result); err != nil {
		t.Fatalf("Put error: %v", err)
	}
	if result["status"] != "updated" {
		t.Errorf("status = %q, want %q", result["status"], "updated")
	}
}

func TestClient_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(200)
	}))
	defer ts.Close()

	c := NewClient("tok", 1)
	c.HTTP = ts.Client()

	if err := c.Delete(ts.URL + "/api"); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
}

func TestClient_Do_401ReturnsAuthError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		_, _ = w.Write([]byte(`{"status_code":401,"errors":[{"type":"unauthorized","messages":["invalid token"]}]}`)) //nolint:errcheck
	}))
	defer ts.Close()

	c := NewClient("bad-token", 1)
	c.HTTP = ts.Client()

	req, _ := http.NewRequest("GET", ts.URL+"/api", nil)
	_, err := c.Do(req)
	if err == nil {
		t.Fatal("expected error for 401")
	}
	if _, ok := err.(*cerrors.AuthError); !ok {
		t.Errorf("expected *AuthError, got %T: %v", err, err)
	}
}

func TestClient_Do_404ReturnsNotFoundError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(`{}`)) //nolint:errcheck
	}))
	defer ts.Close()

	c := NewClient("tok", 1)
	c.HTTP = ts.Client()

	req, _ := http.NewRequest("GET", ts.URL+"/api", nil)
	_, err := c.Do(req)
	if err == nil {
		t.Fatal("expected error for 404")
	}
	if _, ok := err.(*cerrors.NotFoundError); !ok {
		t.Errorf("expected *NotFoundError, got %T: %v", err, err)
	}
}

func TestClient_Do_400ReturnsAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"status_code":400,"errors":[{"type":"invalid_param","messages":["amount is required"]}]}`)) //nolint:errcheck
	}))
	defer ts.Close()

	c := NewClient("tok", 1)
	c.HTTP = ts.Client()

	req, _ := http.NewRequest("GET", ts.URL+"/api", nil)
	_, err := c.Do(req)
	if err == nil {
		t.Fatal("expected error for 400")
	}
	apiErr, ok := err.(*cerrors.APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", apiErr.StatusCode)
	}
	if apiErr.Code != "invalid_param" {
		t.Errorf("Code = %q, want %q", apiErr.Code, "invalid_param")
	}
}

func TestRetryAfterDuration(t *testing.T) {
	tests := []struct {
		name       string
		retryAfter string
		attempt    int
		wantSecs   int
	}{
		{"with header", "5", 0, 5},
		{"no header fallback", "", 0, 1},
		{"no header attempt 1", "", 1, 2},
		{"invalid header", "abc", 0, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{Header: http.Header{}}
			if tt.retryAfter != "" {
				resp.Header.Set("Retry-After", tt.retryAfter)
			}
			got := retryAfterDuration(resp, tt.attempt)
			want := tt.wantSecs
			if int(got.Seconds()) != want {
				t.Errorf("got %v, want %ds", got, want)
			}
		})
	}
}

func TestParseAPIError_FreeeFormat(t *testing.T) {
	body := `{"status_code":400,"errors":[{"type":"validation","messages":["field is required","invalid format"]}]}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(body)) //nolint:errcheck
	}))
	defer ts.Close()

	c := NewClient("tok", 1)
	c.HTTP = ts.Client()

	req, _ := http.NewRequest("GET", ts.URL+"/api", nil)
	_, err := c.Do(req)

	apiErr, ok := err.(*cerrors.APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Code != "validation" {
		t.Errorf("Code = %q, want %q", apiErr.Code, "validation")
	}
	if apiErr.Message != "[field is required invalid format]" {
		t.Errorf("Message = %q", apiErr.Message)
	}
}

func TestParseAPIError_NonFreeeFormat(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(`Internal Server Error`)) //nolint:errcheck
	}))
	defer ts.Close()

	c := NewClient("tok", 1)
	c.HTTP = ts.Client()

	req, _ := http.NewRequest("GET", ts.URL+"/api", nil)
	_, err := c.Do(req)

	apiErr, ok := err.(*cerrors.APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", apiErr.StatusCode)
	}
	if apiErr.Message != "Internal Server Error" {
		t.Errorf("Message = %q", apiErr.Message)
	}
}
