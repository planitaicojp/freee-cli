# Wallet-Txn Command Design Spec

## Overview

freee CLI에 wallet-txn (口座明細 / 구좌 명세) 커맨드 추가.
freee API `/api/1/wallet_txns` 엔드포인트를 래핑하여 list/show/create/delete 4개 서브커맨드 제공.

API에 update 엔드포인트가 없으므로 CRUD가 아닌 CRD + show 구성.

## Design Decisions

| 결정 항목 | 선택 | 이유 |
|-----------|------|------|
| walletable 지정 | type + id 분리 플래그 | transfer 패턴과 일관, `freee walletable list` 결과 복붙 가능 |
| entry-side 필터 | list에 `--entry-side` 제공 | API 네이티브 필터, 입금/출금 구분은 기본 유스케이스 |
| status 표시 | table은 영어 라벨, JSON은 raw 숫자 | agent 파싱 용이, API 충실성 보장 |
| balance 플래그 | create에 `--balance` optional 제공 | API 지원 필드, 수동 입력 시 잔고 기록 유스케이스 |
| 파일 구조 | 단일 파일 (transfer 패턴) | ~250줄 예상, 기존 패턴 일관성 |
| SPEC.md 커맨드 수 | 3개로 수정 (update 제거) | API 실태에 맞게 조정, 총 271개 |

## Reference Pattern

`cmd/transfer/transfer.go` -- 동일 CRUD 패턴을 답습하되, update 관련 로직은 제외.

## File Changes

| 파일 | 변경 | 설명 |
|------|------|------|
| `cmd/wallettxn/wallettxn.go` | 신규 | 4개 서브커맨드 (~250줄) |
| `internal/model/wallet_txn.go` | 신규 | WalletTxn/Response/Row + status 매핑 (~60줄) |
| `internal/api/freee.go` | 수정 | 4개 API 메서드 추가 |
| `cmd/root.go` | 수정 | `wallettxn.Cmd` 등록 |
| `SPEC.md` | 수정 | #63-66 -> #63-65+#66, 커맨드 수 271로 조정 |

## Command Tree

```
freee wallet-txn list     [--walletable-type TYPE --walletable-id ID] [--entry-side income|expense] [--from DATE] [--to DATE] [--limit N] [--offset N] [--all] [--format json|yaml|csv|table]
freee wallet-txn show     <id> [--format json|yaml|csv|table]
freee wallet-txn create   --walletable-type TYPE --walletable-id ID --entry-side income|expense --amount N --date DATE [--description TEXT] [--balance N] [--dry-run]
freee wallet-txn delete   <id> [--dry-run]
```

## Model

```go
// internal/model/wallet_txn.go

type WalletTxnsResponse struct {
    WalletTxns []WalletTxn `json:"wallet_txns"`
}

type WalletTxnResponse struct {
    WalletTxn WalletTxn `json:"wallet_txn"`
}

type WalletTxn struct {
    ID             int64  `json:"id"`
    CompanyID      int64  `json:"company_id"`
    Date           string `json:"date"`
    Amount         int64  `json:"amount"`
    DueAmount      int64  `json:"due_amount"`
    Balance        *int64 `json:"balance"`         // nullable
    EntrySide      string `json:"entry_side"`      // income | expense
    WalletableType string `json:"walletable_type"` // bank_account | credit_card | wallet
    WalletableID   int64  `json:"walletable_id"`
    Description    string `json:"description"`
    Status         int    `json:"status"`           // 1,2,3,4,6
    RuleMatched    bool   `json:"rule_matched"`
}

type WalletTxnRow struct {
    ID          int64  `json:"id"`
    Date        string `json:"date"`
    EntrySide   string `json:"entry_side"`
    Amount      int64  `json:"amount"`
    Walletable  string `json:"walletable"`    // "bank_account:123"
    Status      string `json:"status"`        // "waiting" etc.
    Description string `json:"description"`
}
```

### Status Label Mapping

| Code | Label |
|------|-------|
| 1 | waiting |
| 2 | settled |
| 3 | ignored |
| 4 | settling |
| 6 | excluded |

Table 출력 시 영어 라벨로 변환. JSON/YAML/CSV 출력은 raw API 응답 (숫자 코드) 유지.

## API Methods

`internal/api/freee.go`에 추가:

```go
func (a *FreeeAPI) ListWalletTxns(companyID int64, params string, result any) error
// GET /api/1/wallet_txns?company_id={companyID}&{params}

func (a *FreeeAPI) GetWalletTxn(companyID, id int64, result any) error
// GET /api/1/wallet_txns/{id}?company_id={companyID}

func (a *FreeeAPI) CreateWalletTxn(body, result any) error
// POST /api/1/wallet_txns

func (a *FreeeAPI) DeleteWalletTxn(companyID, id int64) error
// DELETE /api/1/wallet_txns/{id}?company_id={companyID}
```

## Command Behavior

### list

- `cmdutil.ResolveFiscalYear()`로 기본 일자 범위 자동 설정
- `--walletable-type`과 `--walletable-id`는 반드시 동시 지정 (한쪽만 지정 시 CLI 에러)
- `--entry-side`는 `income` 또는 `expense`만 허용
- `--all` 지정 시: offset을 100씩 루프하여 전건 취득
- 테이블 출력: ID | Date | EntrySide | Amount | Walletable | Status | Description
- JSON/YAML/CSV: raw API 응답 그대로

### show

- 인수: wallet-txn ID (필수)
- 테이블: key-value 포맷
  ```
  ID:           1
  Date:         2019-12-17
  Entry Side:   income
  Amount:       5000
  Due Amount:   0
  Balance:      10000
  Walletable:   bank_account:1
  Status:       waiting
  Description:  振込 カ）ABC
  Rule Matched: true
  ```
- Balance가 null이면 해당 행 생략
- JSON/YAML/CSV: raw 응답

### create

- 필수 플래그: `--walletable-type`, `--walletable-id`, `--entry-side`, `--amount`, `--date`
- 임의 플래그: `--description`, `--balance`
- `--walletable-type` 검증: `bank_account`, `credit_card`, `wallet`만 허용
- `--entry-side` 검증: `income`, `expense`만 허용
- `--dry-run`: 리퀘스트 바디만 출력, API 호출 없음
- `--format` 지원 (transfer create와 동일)
- 성공 시: 생성된 wallet-txn 상세 출력

### delete

- 인수: wallet-txn ID (필수)
- `--dry-run` 지원
- 성공 시: `Deleted wallet transaction {id}` 출력 (stderr)
- 주의: 동기화로 취득한 데이터는 삭제 불가 (API 에러 반환)

## Flags Summary

| Flag | list | show | create | delete |
|------|------|------|--------|--------|
| `--walletable-type` | optional (pair) | | required | |
| `--walletable-id` | optional (pair) | | required | |
| `--entry-side` | optional filter | | required | |
| `--from` (date) | o | | | |
| `--to` (date) | o | | | |
| `--limit` | o | | | |
| `--offset` | o | | | |
| `--all` | o | | | |
| `--format` | o | o | o | |
| `--date` | | | required | |
| `--amount` | | | required | |
| `--description` | | | optional | |
| `--balance` | | | optional | |
| `--dry-run` | | | o | o |

## Validation

- `--walletable-type`: `bank_account`, `credit_card`, `wallet`만 허용 (transfer의 `validateWalletableType`을 `cmdutil`로 추출하여 공유)
- `--entry-side`: `income`, `expense`만 허용
- list에서 `--walletable-type`과 `--walletable-id` 중 하나만 지정 시 에러: `both --walletable-type and --walletable-id must be specified together`

## Error Handling

기존 패턴 준수:
- 인증 에러 -> exit code 2, `hint: run 'freee auth login'`
- Not Found -> exit code 3
- 검증 에러 -> exit code 4
- API 에러 -> exit code 5

## Out of Scope

- Walletable 이름 해결 (이름 -> type+id) -- 필요 시 추후 추가
- Status 필터 -- API 미지원
- Update 커맨드 -- API 미제공
- 클라이언트 사이드 필터링 -- `jq`로 대체 가능
