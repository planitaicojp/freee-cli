package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/planitaicojp/freee-cli/internal/config"
)

// Debug levels
const (
	DebugOff     = 0
	DebugVerbose = 1
	DebugAPI     = 2
)

var debugLevel = DebugOff

func init() {
	switch os.Getenv(config.EnvDebug) {
	case "api":
		debugLevel = DebugAPI
	case "1", "true":
		debugLevel = DebugVerbose
	}
}

// SetDebugLevel sets the debug logging level.
func SetDebugLevel(level int) {
	if level > debugLevel {
		debugLevel = level
	}
}

func debugLogRequest(req *http.Request, body []byte) {
	if debugLevel < DebugVerbose {
		return
	}
	fmt.Fprintf(os.Stderr, "> %s %s\n", req.Method, req.URL)
	if debugLevel >= DebugAPI {
		for key, vals := range req.Header {
			for _, val := range vals {
				if isSensitiveHeader(key) {
					val = "****"
				}
				fmt.Fprintf(os.Stderr, "> %s: %s\n", key, val)
			}
		}
		if len(body) > 0 {
			fmt.Fprintf(os.Stderr, "> %s\n", maskSensitiveFields(string(body)))
		}
	}
}

func debugLogResponse(resp *http.Response, elapsed time.Duration, body []byte) {
	if debugLevel < DebugVerbose {
		return
	}
	fmt.Fprintf(os.Stderr, "< %d %s (%s)\n", resp.StatusCode, http.StatusText(resp.StatusCode), elapsed.Truncate(time.Millisecond))
	if debugLevel >= DebugAPI {
		for key, vals := range resp.Header {
			for _, val := range vals {
				if isSensitiveHeader(key) {
					val = "****"
				}
				fmt.Fprintf(os.Stderr, "< %s: %s\n", key, val)
			}
		}
		if len(body) > 0 {
			fmt.Fprintf(os.Stderr, "< %s\n", maskSensitiveFields(string(body)))
		}
	}
	fmt.Fprintln(os.Stderr)
}

func isSensitiveHeader(key string) bool {
	k := strings.ToLower(key)
	return k == "authorization" || k == "x-auth-token" || k == "x-subject-token"
}

func maskSensitiveFields(s string) string {
	// Simple masking for common sensitive fields in JSON
	for _, field := range []string{"password", "access_token", "refresh_token", "client_secret"} {
		s = strings.ReplaceAll(s, fmt.Sprintf(`"%s":"`, field), fmt.Sprintf(`"%s":"****`, field))
	}
	return s
}
