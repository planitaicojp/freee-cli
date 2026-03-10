# freee CLI クイックスタートガイド

## 前提条件

- freee アカウント（[freee.co.jp](https://www.freee.co.jp) で登録）
- freee 開発者アプリ（下記で作成）

## 1. freee 開発者アプリの作成

1. [freee 開発者コンソール](https://app.secure.freee.co.jp/developers/applications/new) にアクセス
2. 以下の設定でアプリを作成：
   - **アプリタイプ**: `プライベート`（自分だけが使う場合）
   - **アプリ名**: 任意（例: `freee-cli`）
   - **コールバックURL**: `http://localhost:8080/callback`
3. 作成後、**Client ID** と **Client Secret** をメモ

> **注意**: アプリは `draft` 状態のままで使用できます。`active`（公開）にする必要はありません。

## 2. ログイン

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

ブラウザが自動的に開き、freee のログイン・認可画面が表示されます。
認可を許可すると、CLI に自動的にトークンが保存されます。

### ブラウザが開かない場合

ターミナルに表示される URL を手動でブラウザに貼り付けてください。

## 3. 認証状態の確認

```bash
$ freee auth status
Profile:   default
Email:     taro@example.com
Company:   株式会社サンプル (ID: 12345678)
Token:     valid (expires in 5h45m, 2026-03-11 02:53 JST)
```

## 4. 基本的な使い方

```bash
# 事業所一覧
$ freee company list

# 取引一覧
$ freee deal list

# JSON 出力（エージェント/スクリプト向け）
$ freee deal list --format json

# 取引先一覧
$ freee partner list

# 勘定科目一覧
$ freee account list
```

## 5. 複数アカウント管理

```bash
# 別のプロファイルでログイン
$ freee auth login --profile sub-account

# プロファイル一覧
$ freee auth list

# プロファイル切り替え
$ freee auth switch sub-account
```

## 6. 事業所の切り替え

複数の事業所にアクセスできる場合：

```bash
# 事業所一覧を確認
$ freee company list

# デフォルト事業所を変更
$ freee company switch <company-id>
```

## 7. CI/CD での利用

```bash
# 環境変数でトークンを直接指定
export FREEE_TOKEN="<access-token>"
export FREEE_COMPANY_ID="12345678"

# 非対話モードで実行
freee deal list --no-input --format json
```

## トラブルシューティング

### 「必須パラメータが不足」エラー

- アプリの **コールバックURL** が `http://localhost:8080/callback` に設定されているか確認
- アプリの **アプリタイプ** が `プライベート` であることを確認

### ポート 8080 が使用中

別のアプリケーションがポート 8080 を使用している場合、そのアプリケーションを停止してから再度ログインしてください。

### トークン期限切れ

アクセストークンは 6 時間有効です。期限切れの場合、次の API 呼び出し時にリフレッシュトークンで自動更新されます。リフレッシュトークンの有効期限は 90 日間です。

## 設定ファイル

設定は `~/.config/freee/` に保存されます：

```
~/.config/freee/
├── config.yaml        # プロファイル設定
└── credentials.yaml   # OAuth トークン（0600 パーミッション）
```
