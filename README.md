# freee - freee API CLI

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[English](README-en.md) | [한국어](README-ko.md)

freee 公開 API 用のコマンドラインインターフェースです。Go で書かれたシングルバイナリで、エージェントフレンドリーな設計を採用しています。

> **注意**: 本ツールは非公式であり、freee 株式会社とは提携・推奨関係にありません。

## 特徴

- シングルバイナリ、クロスプラットフォーム対応（Linux / macOS / Windows）
- OAuth2 Authorization Code + PKCE によるブラウザベースログイン
- 複数プロファイル対応（`gh auth` スタイル）
- 構造化出力（`--format json/yaml/csv/table`）
- エージェントフレンドリー設計（`--no-input`、決定的な終了コード、stderr/stdout 分離）
- トークン自動リフレッシュ（アクセストークン 6 時間、リフレッシュトークン 90 日）
- SDK 不使用 — OpenAPI スペック準拠の直接 HTTP 呼び出し

## インストール

### ソースからビルド

```bash
go install github.com/planitaicojp/freee-cli@latest
```

### Git からビルド

```bash
git clone https://github.com/planitaicojp/freee-cli.git
cd freee-cli
make build
sudo mv freee /usr/local/bin/
```

### リリースバイナリ

[Releases](https://github.com/planitaicojp/freee-cli/releases) ページからダウンロード、または以下のコマンドを使用してください：

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

## 事前準備

1. [freee 開発者コンソール](https://app.secure.freee.co.jp/developers/applications/new) でアプリを作成
   - **アプリタイプ**: `プライベート`
   - **コールバック URL**: `http://localhost:8080/callback`
2. **Client ID** と **Client Secret** をメモ

> アプリは `draft` 状態のまま使用できます。

詳細は [クイックスタートガイド](docs/getting-started.md) を参照してください。

## クイックスタート

```bash
# ログイン（ブラウザが開きます）
freee auth login

# 認証状態を確認
freee auth status

# 事業所一覧
freee company list

# 取引一覧
freee deal list

# JSON 形式で出力
freee deal list --format json

# 取引先一覧
freee partner list
```

## コマンド一覧

| コマンド | 説明 |
|---------|------|
| `freee auth` | 認証管理（login / logout / status / list / switch / token / remove） |
| `freee company` | 事業所管理（list / show / switch） |
| `freee deal` | 取引管理（list / show / create / update / delete） |
| `freee invoice` | 請求書管理（list / show / create / update / delete） |
| `freee partner` | 取引先管理（list / show / create / update / delete） |
| `freee account` | 勘定科目（list / show） |
| `freee section` | 部門管理（list / create / update / delete） |
| `freee tag` | メモタグ管理（list / create / update / delete） |
| `freee item` | 品目管理（list / create / update / delete） |
| `freee journal` | 仕訳帳（list） |
| `freee expense` | 経費申請管理（list / show / create / update / delete） |
| `freee walletable` | 口座管理（list / show） |
| `freee config` | CLI 設定管理（show / set / path） |

## 設定

設定ファイルは `~/.config/freee/` に保存されます：

| ファイル | 説明 | パーミッション |
|---------|------|------------|
| `config.yaml` | プロファイル設定 | 0600 |
| `credentials.yaml` | OAuth トークン | 0600 |

### 環境変数

| 変数 | 説明 |
|-----|------|
| `FREEE_PROFILE` | 使用するプロファイル名 |
| `FREEE_COMPANY_ID` | 事業所 ID |
| `FREEE_TOKEN` | アクセストークン（直接指定、CI 用） |
| `FREEE_FORMAT` | 出力形式 |
| `FREEE_CONFIG_DIR` | 設定ディレクトリ |
| `FREEE_NO_INPUT` | 非対話モード（`1` or `true`） |
| `FREEE_DEBUG` | デバッグログ（`1` or `api`） |

優先順位: 環境変数 > フラグ > プロファイル設定 > デフォルト値

### グローバルフラグ

```
--profile      使用するプロファイル
--format       出力形式（table / json / yaml / csv）
--company-id   事業所 ID（上書き指定）
--no-input     対話プロンプトを無効化
--quiet        不要な出力を抑制
--verbose      詳細出力
--no-color     カラー出力を無効化
```

## 終了コード

| コード | 意味 |
|-------|------|
| 0 | 成功 |
| 1 | 一般エラー |
| 2 | 認証失敗 |
| 3 | リソース未検出 |
| 4 | バリデーションエラー |
| 5 | API エラー |
| 6 | ネットワークエラー |
| 10 | ユーザーキャンセル |

## エージェント連携

本 CLI はスクリプトや AI エージェントからの利用を想定して設計されています：

```bash
# 非対話モードで JSON 出力
freee deal list --format json --no-input

# トークンを取得してスクリプトで利用
TOKEN=$(freee auth token)

# 終了コードでエラーハンドリング
freee deal show 12345 || echo "Exit code: $?"

# CI/CD 環境での利用
export FREEE_TOKEN="<access-token>"
export FREEE_COMPANY_ID="12345678"
freee deal list --format json
```

## 開発

```bash
make build     # バイナリをビルド
make test      # テストを実行
make lint      # リンターを実行
make clean     # 成果物を削除
```

## 対応 API

| API | 状態 |
|-----|------|
| [会計 API](https://developer.freee.co.jp/reference/accounting/reference) | Phase 1 対応中 |
| [人事労務 API](https://developer.freee.co.jp/reference/hr) | Phase 2 予定 |
| [請求書 API](https://developer.freee.co.jp/reference/iv) | Phase 2 予定 |
| [工数管理 API](https://developer.freee.co.jp/reference/time-tracking) | Phase 3 予定 |
| [販売 API](https://developer.freee.co.jp/reference/sales) | Phase 3 予定 |

## 参考

- [freee API リファレンス](https://developer.freee.co.jp)
- [freee OpenAPI スキーマ](https://github.com/freee/freee-api-schema)

## ライセンス

MIT License - 詳細は [LICENSE](LICENSE) をご覧ください。
