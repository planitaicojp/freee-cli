# freee CLI Quick Start Guide

## Prerequisites

- A freee account (register at [freee.co.jp](https://www.freee.co.jp))
- A freee developer app (created below)

## 1. Create a freee Developer App

1. Go to the [freee Developer Console](https://app.secure.freee.co.jp/developers/applications/new)
2. Create an app with the following settings:
   - **App Type (アプリタイプ)**: `プライベート` (Private — for personal use)
   - **App Name**: anything (e.g., `freee-cli`)
   - **Callback URL (コールバックURL)**: `http://localhost:8080/callback`
3. Note down the **Client ID** and **Client Secret**

> **Note**: The app can be used in `draft` status. You do NOT need to make it `active` (public).

## 2. Login

```bash
$ freee auth login
Client ID (from https://app.secure.freee.co.jp/developers): <your-client-id>
Client Secret: ****
Opening browser for authorization...
Waiting for authorization...
Logged in as taro@example.com
Default company: 株式会社サンプル (ID: 12345678)
Token expires: 2026-03-11 02:53 JST
```

A browser window will open automatically, showing the freee login and authorization page.
After granting access, the token is saved to the CLI automatically.

### If the browser doesn't open

Copy the URL shown in the terminal and paste it into your browser manually.

## 3. Check Authentication Status

```bash
$ freee auth status
Profile:   default
Email:     taro@example.com
Company:   株式会社サンプル (ID: 12345678)
Token:     valid (expires in 5h45m, 2026-03-11 02:53 JST)
```

## 4. Basic Usage

```bash
# List companies
$ freee company list

# List deals (transactions)
$ freee deal list

# JSON output (agent/script-friendly)
$ freee deal list --format json

# List partners
$ freee partner list

# List account items
$ freee account list
```

## 5. Multiple Account Management

```bash
# Login with a different profile
$ freee auth login --profile sub-account

# List profiles
$ freee auth list

# Switch profile
$ freee auth switch sub-account
```

## 6. Switch Companies

If you have access to multiple companies:

```bash
# List companies
$ freee company list

# Change default company
$ freee company switch <company-id>
```

## 7. CI/CD Usage

```bash
# Provide token directly via environment variable
export FREEE_TOKEN="<access-token>"
export FREEE_COMPANY_ID="12345678"

# Run in non-interactive mode
freee deal list --no-input --format json
```

## Troubleshooting

### "必須パラメータが不足" (Required Parameters Missing) Error

- Verify the app's **Callback URL** is set to exactly `http://localhost:8080/callback`
- Verify the app's **App Type** is `プライベート` (Private)

### Port 8080 Already in Use

If another application is using port 8080, stop it before attempting to login.

### Token Expiration

Access tokens are valid for 6 hours. When expired, they are automatically refreshed using the refresh token on the next API call. Refresh tokens are valid for 90 days.

## Configuration Files

Configuration is stored in `~/.config/freee/`:

```
~/.config/freee/
├── config.yaml        # Profile settings
└── credentials.yaml   # OAuth tokens (0600 permissions)
```
