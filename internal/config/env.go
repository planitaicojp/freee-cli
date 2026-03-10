package config

import "os"

// Environment variable names
const (
	EnvProfile   = "FREEE_PROFILE"
	EnvCompanyID = "FREEE_COMPANY_ID"
	EnvToken     = "FREEE_TOKEN"
	EnvFormat    = "FREEE_FORMAT"
	EnvConfigDir = "FREEE_CONFIG_DIR"
	EnvNoInput   = "FREEE_NO_INPUT"
	EnvDebug     = "FREEE_DEBUG"
)

// EnvOr returns the environment variable value if set, otherwise the fallback.
func EnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// IsNoInput returns true if non-interactive mode is requested.
func IsNoInput() bool {
	return os.Getenv(EnvNoInput) == "1" || os.Getenv(EnvNoInput) == "true"
}
