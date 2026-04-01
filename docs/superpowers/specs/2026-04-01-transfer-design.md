# Transfer Command Design Spec

## Overview

freee CLI에 transfer (口座振替 / 이체 전표) CRUD 커맨드 추가.
freee API `/api/1/transfers` 엔드포인트를 래핑하여 list/show/create/update/delete 5개 서브커맨드 제공.

## Design Decisions

| 결정 항목 | 선택 | 이유 |
|-----------|------|------|
| create 입력 모드 | 플래그 전용 | 필드 6~7개로 단순, JSON 모드 불필요 (YAGNI) |
| walletable 지정 | type + id 분리 플래그 | API 1:1 대응, `freee walletable list` 결과 복붙 가능 |
| list 필터 | API 제공 필터만 | 클라이언트 사이드 필터 불필요, `jq`로 대체 가능 |
| 파일 구조 | 단일 파일 (manual-journal 패턴) | ~300줄 예상, 분리하면 오히려 파편화 |

## Reference Pattern

`cmd/manualjournal/manualjournal.go` — 동일 CRUD 패턴을 답습하되, JSON 입력 모드와 복잡한 details 관련 로직은 제외.

## File Changes

| 파일 | 변경 | 설명 |
|------|------|------|
| `cmd/transfer/transfer.go` | 신규 | 5개 서브커맨드 (~300줄) |
| `internal/model/transfer.go` | 신규 | Transfer/TransferResponse/TransferRow (~50줄) |
| `internal/api/freee.go` | 수정 | 5개 API 메서드 추가 |
| `cmd/root.go` | 수정 | `transfer.Cmd` 등록 |

## Command Tree

```
freee transfer list     [--from DATE] [--to DATE] [--limit N] [--offset N] [--all] [--format json|yaml|csv|table]
freee transfer show     <id> [--format json|yaml|csv|table]
freee transfer create   --date DATE --amount N --from-type TYPE --from-id ID --to-type TYPE --to-id ID [--description TEXT] [--dry-run]
freee transfer update   <id> [--date] [--amount] [--from-type] [--from-id] [--to-type] [--to-id] [--description] [--dry-run]
freee transfer delete   <id>
```

## Model

```go
// internal/model/transfer.go

type TransfersResponse struct {
    Transfers []Transfer `json:"transfers"`
}

type TransferResponse struct {
    Transfer Transfer `json:"transfer"`
}

type Transfer struct {
    ID                 int64  `json:"id"`
    CompanyID          int64  `json:"company_id"`
    Date               string `json:"date"`
    Amount             int64  `json:"amount"`
    FromWalletableType string `json:"from_walletable_type"`
    FromWalletableID   int64  `json:"from_walletable_id"`
    ToWalletableType   string `json:"to_walletable_type"`
    ToWalletableID     int64  `json:"to_walletable_id"`
    Description        string `json:"description"`
}

type TransferRow struct {
    ID          int64  `json:"id"`
    Date        string `json:"date"`
    Amount      int64  `json:"amount"`
    From        string `json:"from"`        // "bank_account:12345" format
    To          string `json:"to"`          // "wallet:67890" format
    Description string `json:"description"`
}

func (t Transfer) ToRow() TransferRow {
    return TransferRow{
        ID:          t.ID,
        Date:        t.Date,
        Amount:      t.Amount,
        From:        fmt.Sprintf("%s:%d", t.FromWalletableType, t.FromWalletableID),
        To:          fmt.Sprintf("%s:%d", t.ToWalletableType, t.ToWalletableID),
        Description: t.Description,
    }
}
```

## API Methods

`internal/api/freee.go` に追加:

```go
func (a *FreeeAPI) ListTransfers(companyID int64, params string, result any) error
// GET /api/1/transfers?company_id={companyID}&{params}

func (a *FreeeAPI) GetTransfer(companyID, id int64, result any) error
// GET /api/1/transfers/{id}?company_id={companyID}

func (a *FreeeAPI) CreateTransfer(body, result any) error
// POST /api/1/transfers

func (a *FreeeAPI) UpdateTransfer(id int64, body, result any) error
// PUT /api/1/transfers/{id}

func (a *FreeeAPI) DeleteTransfer(companyID, id int64) error
// DELETE /api/1/transfers/{id}?company_id={companyID}
```

## Command Behavior

### list

- `cmdutil.ResolveFiscalYear()` で基本日付範囲を自動設定
- `--all` 指定時: offset を 100 ずつループし全件取得
- テーブル出力: ID | Date | Amount | From | To | Description
- JSON/YAML/CSV: raw API レスポンスをそのまま出力

### show

- 引数: transfer ID (必須)
- テーブル: key-value フォーマット
- JSON/YAML/CSV: raw レスポンス

### create

- 必須フラグ: `--date`, `--amount`, `--from-type`, `--from-id`, `--to-type`, `--to-id`
- 任意フラグ: `--description`
- `--from-type` / `--to-type` のバリデーション: `bank_account`, `credit_card`, `wallet` のみ許可
- `--dry-run`: リクエストボディのみ出力、API 呼び出しなし
- 成功時: 作成された transfer の詳細を出力

### update (GET-then-merge)

1. GET で現在の transfer を取得
2. 指定されたフラグの値のみ上書き
3. PUT で更新
- 全フラグ任意 (変更するもののみ指定)
- `--dry-run` サポート
- 成功時: 更新後の transfer 詳細を出力

### delete

- 引数: transfer ID (必須)
- 成功時: `Deleted transfer {id}` 出力

## Flags Summary

| Flag | list | show | create | update | delete |
|------|------|------|--------|--------|--------|
| `--from` (date) | o | | | | |
| `--to` (date) | o | | | | |
| `--limit` | o | | | | |
| `--offset` | o | | | | |
| `--all` | o | | | | |
| `--format` | o | o | | | |
| `--date` | | | required | optional | |
| `--amount` | | | required | optional | |
| `--from-type` | | | required | optional | |
| `--from-id` | | | required | optional | |
| `--to-type` | | | required | optional | |
| `--to-id` | | | required | optional | |
| `--description` | | | optional | optional | |
| `--dry-run` | | | o | o | |

## Walletable Type Validation

`--from-type` と `--to-type` は以下のいずれかのみ許可:
- `bank_account` — 銀行口座
- `credit_card` — クレジットカード
- `wallet` — 現金・その他

無効な値が指定された場合、CLI レベルでエラーを返す (API 呼び出し前)。

## Error Handling

既存パターンに準拠:
- 認証エラー → exit code 2, `hint: run 'freee auth login'`
- Not Found → exit code 3
- バリデーションエラー → exit code 4
- API エラー → exit code 5

## Out of Scope

- Walletable 名前解決 (名前 → type+id) — 必要になった時に追加
- JSON 入力モード — フィールドが少ないため不要
- クライアントサイドフィルタリング — `jq` で代替可能
