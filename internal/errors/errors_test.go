package errors

import (
	"fmt"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want string
	}{
		{
			name: "with code",
			err:  APIError{StatusCode: 400, Code: "invalid_param", Message: "bad request"},
			want: "API error (HTTP 400, invalid_param): bad request",
		},
		{
			name: "without code",
			err:  APIError{StatusCode: 500, Message: "internal"},
			want: "API error (HTTP 500): internal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAPIError_ExitCode(t *testing.T) {
	err := &APIError{StatusCode: 400}
	if got := err.ExitCode(); got != ExitAPI {
		t.Errorf("got %d, want %d", got, ExitAPI)
	}
}

func TestAuthError(t *testing.T) {
	err := &AuthError{Message: "token expired"}
	if got := err.Error(); got != "auth error: token expired\nhint: run 'freee auth login' to authenticate" {
		t.Errorf("unexpected error message: %q", got)
	}
	if got := err.ExitCode(); got != ExitAuth {
		t.Errorf("got exit code %d, want %d", got, ExitAuth)
	}
}

func TestConfigError(t *testing.T) {
	err := &ConfigError{Message: "no profile"}
	if got := err.Error(); got != "config error: no profile\nhint: run 'freee auth login' to set up a profile" {
		t.Errorf("unexpected error message: %q", got)
	}
	if got := err.ExitCode(); got != ExitGeneral {
		t.Errorf("got exit code %d, want %d", got, ExitGeneral)
	}
}

func TestNotFoundError(t *testing.T) {
	err := &NotFoundError{Resource: "deal", ID: "123"}
	if got := err.Error(); got != "deal not found: 123\nhint: check the ID and ensure the resource exists in the current company" {
		t.Errorf("unexpected error message: %q", got)
	}
	if got := err.ExitCode(); got != ExitNotFound {
		t.Errorf("got exit code %d, want %d", got, ExitNotFound)
	}
}

func TestValidationError(t *testing.T) {
	tests := []struct {
		name string
		err  ValidationError
		want string
	}{
		{
			name: "with field",
			err:  ValidationError{Field: "amount", Message: "must be positive"},
			want: "validation error on amount: must be positive",
		},
		{
			name: "without field",
			err:  ValidationError{Message: "invalid input"},
			want: "validation error: invalid input",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
	if got := (&ValidationError{}).ExitCode(); got != ExitValidation {
		t.Errorf("got exit code %d, want %d", got, ExitValidation)
	}
}

func TestNetworkError(t *testing.T) {
	inner := fmt.Errorf("connection refused")
	err := &NetworkError{Err: inner}
	if got := err.Error(); got != "network error: connection refused\nhint: check your internet connection and try again" {
		t.Errorf("unexpected error message: %q", got)
	}
	if got := err.Unwrap(); got != inner {
		t.Errorf("Unwrap returned %v, want %v", got, inner)
	}
	if got := err.ExitCode(); got != ExitNetwork {
		t.Errorf("got exit code %d, want %d", got, ExitNetwork)
	}
}

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, ExitOK},
		{"api error", &APIError{StatusCode: 400}, ExitAPI},
		{"auth error", &AuthError{}, ExitAuth},
		{"not found", &NotFoundError{}, ExitNotFound},
		{"validation", &ValidationError{}, ExitValidation},
		{"network", &NetworkError{Err: fmt.Errorf("x")}, ExitNetwork},
		{"config", &ConfigError{}, ExitGeneral},
		{"generic error", fmt.Errorf("unknown"), ExitGeneral},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetExitCode(tt.err); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
