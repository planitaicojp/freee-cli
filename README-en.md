# freee - freee API CLI

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[日本語](README.md) | [한국어](README-ko.md)

A command-line interface for the freee public API. Written in Go as a single binary with an agent-friendly design.

> **Note**: This is an unofficial tool and is not affiliated with or endorsed by freee K.K.

## Features

- Single binary, cross-platform (Linux / macOS / Windows)
- OAuth2 Authorization Code + PKCE browser-based login
- Multiple profile support (`gh auth` style)
- Structured output (`--format json/yaml/csv/table`)
- Agent-friendly design (`--no-input`, deterministic exit codes, stderr/stdout separation)
- Automatic token refresh (access token 6h, refresh token 90 days)
- No SDK dependency — direct HTTP calls based on OpenAPI spec

## Installation

### Build from source

```bash
go install github.com/planitaicojp/freee-cli@latest
```

### Build from Git

```bash
git clone https://github.com/planitaicojp/freee-cli.git
cd freee-cli
make build
sudo mv freee /usr/local/bin/
```

### Release binaries

Download from the [Releases](https://github.com/planitaicojp/freee-cli/releases) page, or use the following commands:

**Linux (amd64)**

```bash
curl -Lo freee https://github.com/planitaicojp/freee-cli/releases/latest/download/freee-linux-amd64
chmod +x freee
sudo mv freee /usr/local/bin/
```

**macOS (Apple Silicon)**

```bash
curl -Lo freee https://github.com/planitaicojp/freee-cli/releases/latest/download/freee-darwin-arm64
chmod +x freee
sudo mv freee /usr/local/bin/
```

**Windows (amd64)**

```powershell
Invoke-WebRequest -Uri https://github.com/planitaicojp/freee-cli/releases/latest/download/freee-windows-amd64.exe -OutFile freee.exe
```

## Prerequisites

1. Create an app at the [freee Developer Console](https://app.secure.freee.co.jp/developers/applications/new)
   - **App Type**: `プライベート` (Private)
   - **Callback URL**: `http://localhost:8080/callback`
2. Note down the **Client ID** and **Client Secret**

> The app can be used in `draft` status. No need to publish it.

See the [Quick Start Guide](docs/getting-started-en.md) for detailed setup instructions.

## Quick Start

```bash
# Login (opens browser)
freee auth login

# Check authentication status
freee auth status

# List companies
freee company list

# List deals (transactions)
freee deal list

# JSON output
freee deal list --format json

# List partners
freee partner list
```

## Commands

| Command | Description |
|---------|-------------|
| `freee auth` | Authentication (login / logout / status / list / switch / token / remove) |
| `freee company` | Company management (list / show / switch) |
| `freee deal` | Deal management (list / show / create / update / delete) |
| `freee invoice` | Invoice management (list / show / create / update / delete) |
| `freee partner` | Partner management (list / show / create / update / delete) |
| `freee account` | Account items (list / show) |
| `freee section` | Section management (list / create / update / delete) |
| `freee tag` | Tag management (list / create / update / delete) |
| `freee item` | Item management (list / create / update / delete) |
| `freee journal` | Journals (list) |
| `freee expense` | Expense applications (list / show / create / update / delete) |
| `freee walletable` | Walletable management (list / show) |
| `freee config` | CLI configuration (show / set / path) |

## Configuration

Configuration files are stored in `~/.config/freee/`:

| File | Description | Permissions |
|------|-------------|-------------|
| `config.yaml` | Profile settings | 0600 |
| `credentials.yaml` | OAuth tokens | 0600 |

### Environment Variables

| Variable | Description |
|----------|-------------|
| `FREEE_PROFILE` | Profile name to use |
| `FREEE_COMPANY_ID` | Company ID |
| `FREEE_TOKEN` | Access token (direct, for CI) |
| `FREEE_FORMAT` | Output format |
| `FREEE_CONFIG_DIR` | Config directory |
| `FREEE_NO_INPUT` | Non-interactive mode (`1` or `true`) |
| `FREEE_DEBUG` | Debug logging (`1` or `api`) |

Priority: Environment variables > Flags > Profile settings > Defaults

### Global Flags

```
--profile      Profile to use
--format       Output format (table / json / yaml / csv)
--company-id   Company ID (override)
--no-input     Disable interactive prompts
--quiet        Suppress non-essential output
--verbose      Verbose output
--no-color     Disable color output
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Authentication failure |
| 3 | Resource not found |
| 4 | Validation error |
| 5 | API error |
| 6 | Network error |
| 10 | User cancelled |

## Agent Integration

This CLI is designed for use with scripts and AI agents:

```bash
# Non-interactive mode with JSON output
freee deal list --format json --no-input

# Get token for scripting
TOKEN=$(freee auth token)

# Error handling with exit codes
freee deal show 12345 || echo "Exit code: $?"

# CI/CD usage
export FREEE_TOKEN="<access-token>"
export FREEE_COMPANY_ID="12345678"
freee deal list --format json
```

## Development

```bash
make build     # Build binary
make test      # Run tests
make lint      # Run linter
make clean     # Clean artifacts
```

## Supported APIs

| API | Status |
|-----|--------|
| [Accounting API](https://developer.freee.co.jp/reference/accounting/reference) | Phase 1 in progress |
| [HR API](https://developer.freee.co.jp/reference/hr) | Phase 2 planned |
| [Invoice API](https://developer.freee.co.jp/reference/iv) | Phase 2 planned |
| [Time Tracking API](https://developer.freee.co.jp/reference/time-tracking) | Phase 3 planned |
| [Sales API](https://developer.freee.co.jp/reference/sales) | Phase 3 planned |

## References

- [freee API Reference](https://developer.freee.co.jp)
- [freee OpenAPI Schema](https://github.com/freee/freee-api-schema)

## License

MIT License - See [LICENSE](LICENSE) for details.
