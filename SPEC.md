# freee CLI Specification

**Version**: 0.1.0
**Status**: 설계 중 (Design Phase)
**Binary**: `freee`
**Repository**: `planitai-freee-cli`

## 배경

freee는 일본 최대 클라우드 회계/HR SaaS이며, 5종의 공개 API를 제공한다.
공식 SDK는 대부분 archived 상태이고, CLI 도구는 존재하지 않는다.
Agent-first CLI 시대에 맞춰, 단일 실행파일로 freee API를 조작할 수 있는 CLI를 제공한다.

### SDK 미사용 방침

- 공식 SDK: Java, JS, C# 모두 **archived**. PHP만 유지
- 커뮤니티 Go SDK (`LayerXcom/freee-go`): 업데이트가 API 변경보다 느릴 수 있음
- **결정**: SDK를 사용하지 않고, OpenAPI 스펙 기반으로 `net/http` 직접 호출
- **장점**: API 변경에 즉시 대응 가능, 의존성 최소화, 빌드 크기 감소

### 참고: freee-mcp

freee는 공식 MCP 서버 (`freee/freee-mcp`, 303 stars)를 제공하지만,
이는 Claude Code 등 AI 에이전트용이며 단독 CLI로는 사용할 수 없다.
본 프로젝트는 전통적 CLI 도구로서 쉘 스크립트, CI/CD, 다양한 AI 에이전트에서 활용 가능하다.

---

## OAuth2 인증 플로우

### App 등록 (사전 준비)

1. https://app.secure.freee.co.jp/developers 에서 앱 등록
2. Redirect URI: `http://localhost:8080/callback` (CLI 로컬 서버)
3. Client ID / Client Secret 획득

### 로그인 플로우

```
freee auth login
```

1. CLI가 랜덤 `state` + PKCE `code_verifier`/`code_challenge` 생성
2. 로컬 HTTP 서버 시작 (`localhost:8080`)
3. 브라우저 오픈 → `https://accounts.secure.freee.co.jp/public_api/authorize`
4. 사용자가 freee에서 인가 → redirect → CLI가 `code` 수신
5. `code` → token exchange → access_token + refresh_token 저장
6. 사업소(company) 목록 조회 → 기본 사업소 설정

### 토큰 관리

- Access token: JWT, 유효기간 있음 → 자동 갱신
- Refresh token: `~/.config/freee/credentials.yaml`에 저장 (0600)
- 매 API 호출 전 토큰 유효성 확인, 만료 시 자동 refresh

---

## 명령어 목록

### Phase 1: 기반 + Accounting API 핵심

| # | 명령어 | 설명 | 진행 |
|---|--------|------|------|
| **인증/설정** | | | |
| 1 | `freee auth login` | OAuth2 브라우저 로그인 | ☐ |
| 2 | `freee auth logout` | 토큰 삭제 | ☐ |
| 3 | `freee auth status` | 현재 인증 상태 표시 | ☐ |
| 4 | `freee auth list` | 등록된 모든 프로필 목록 | ☐ |
| 5 | `freee auth switch <profile>` | 프로필 전환 | ☐ |
| 6 | `freee auth remove <profile>` | 프로필 삭제 | ☐ |
| 7 | `freee auth token` | 현재 access token 출력 (파이프용) | ☐ |
| 8 | `freee config show` | 현재 설정 표시 | ☐ |
| 9 | `freee config set <key> <value>` | 설정 변경 | ☐ |
| 10 | `freee config path` | 설정 파일 경로 출력 | ☐ |
| **사업소** | | | |
| 11 | `freee company list` | 사업소 목록 | ☐ |
| 12 | `freee company show [id]` | 사업소 상세 | ☐ |
| 13 | `freee company switch <id>` | 기본 사업소 전환 | ☐ |
| **거래 (Deals)** | | | |
| 14 | `freee deal list` | 거래 목록 | ☐ |
| 15 | `freee deal show <id>` | 거래 상세 | ☐ |
| 16 | `freee deal create` | 거래 등록 | ☐ |
| 17 | `freee deal update <id>` | 거래 수정 | ☐ |
| 18 | `freee deal delete <id>` | 거래 삭제 | ☐ |
| **청구서 (Invoices)** | | | |
| 19 | `freee invoice list` | 청구서 목록 | ☐ |
| 20 | `freee invoice show <id>` | 청구서 상세 | ☐ |
| 21 | `freee invoice create` | 청구서 작성 | ☐ |
| 22 | `freee invoice update <id>` | 청구서 수정 | ☐ |
| 23 | `freee invoice delete <id>` | 청구서 삭제 | ☐ |
| **거래처 (Partners)** | | | |
| 24 | `freee partner list` | 거래처 목록 | ☐ |
| 25 | `freee partner show <id>` | 거래처 상세 | ☐ |
| 26 | `freee partner create` | 거래처 등록 | ☐ |
| 27 | `freee partner update <id>` | 거래처 수정 | ☐ |
| 28 | `freee partner delete <id>` | 거래처 삭제 | ☐ |
| **勘定科目 (Account Items)** | | | |
| 29 | `freee account list` | 勘定科目 목록 | ☐ |
| 30 | `freee account show <id>` | 勘定科目 상세 | ☐ |
| **부문 (Sections)** | | | |
| 31 | `freee section list` | 부문 목록 | ☐ |
| 32 | `freee section create` | 부문 등록 | ☐ |
| 33 | `freee section update <id>` | 부문 수정 | ☐ |
| 34 | `freee section delete <id>` | 부문 삭제 | ☐ |
| **메모태그 (Tags)** | | | |
| 35 | `freee tag list` | 태그 목록 | ☐ |
| 36 | `freee tag create` | 태그 등록 | ☐ |
| 37 | `freee tag update <id>` | 태그 수정 | ☐ |
| 38 | `freee tag delete <id>` | 태그 삭제 | ☐ |
| **품목 (Items)** | | | |
| 39 | `freee item list` | 품목 목록 | ☐ |
| 40 | `freee item create` | 품목 등록 | ☐ |
| 41 | `freee item update <id>` | 품목 수정 | ☐ |
| 42 | `freee item delete <id>` | 품목 삭제 | ☐ |
| **분개 (Journals)** | | | |
| 43 | `freee journal list` | 분개 목록 다운로드 | ☐ |
| **경비신청 (Expenses)** | | | |
| 44 | `freee expense list` | 경비신청 목록 | ☐ |
| 45 | `freee expense show <id>` | 경비신청 상세 | ☐ |
| 46 | `freee expense create` | 경비신청 작성 | ☐ |
| 47 | `freee expense update <id>` | 경비신청 수정 | ☐ |
| 48 | `freee expense delete <id>` | 경비신청 삭제 | ☐ |
| **구좌 (Walletables)** | | | |
| 49 | `freee walletable list` | 구좌 목록 | ☐ |
| 50 | `freee walletable show <id>` | 구좌 상세 | ☐ |
| **유틸리티** | | | |
| 51 | `freee version` | 버전 표시 | ☐ |
| 52 | `freee completion` | 쉘 자동완성 생성 | ☐ |

### Phase 2: HR + Invoice API

| # | 명령어 | 설명 | 진행 |
|---|--------|------|------|
| **직원 (Employees)** | | | |
| 53 | `freee employee list` | 직원 목록 | ☐ |
| 54 | `freee employee show <id>` | 직원 상세 | ☐ |
| **근태 (Attendance)** | | | |
| 55 | `freee attendance list` | 근태 목록 | ☐ |
| 56 | `freee attendance update` | 근태 수정 | ☐ |
| **급여명세 (Payslips)** | | | |
| 57 | `freee payslip list` | 급여명세 목록 | ☐ |
| 58 | `freee payslip show <id>` | 급여명세 상세 | ☐ |

### Phase 3: Time Tracking + Sales API

| # | 명령어 | 설명 | 진행 |
|---|--------|------|------|
| **공수 (Time Tracking)** | | | |
| 59 | `freee timetrack project list` | 프로젝트 목록 | ☐ |
| 60 | `freee timetrack log list` | 작업시간 목록 | ☐ |
| 61 | `freee timetrack log create` | 작업시간 등록 | ☐ |
| **판매 (Sales)** | | | |
| 62 | `freee sales order list` | 수주 목록 | ☐ |
| 63 | `freee sales order show <id>` | 수주 상세 | ☐ |

---

## Usage 예시

### 인증

```bash
# 로그인 (브라우저가 열림)
$ freee auth login
✓ Logged in as taro@example.com
✓ Default company: 株式会社サンプル (ID: 1234567)

# 다른 계정으로 추가 로그인
$ freee auth login --profile sub-account

# 인증 목록
$ freee auth list
  PROFILE        USER                    COMPANY                STATUS
✓ default        taro@example.com        株式会社サンプル        Active
  sub-account    hanako@example.com      合同会社テスト          Active

# 프로필 전환
$ freee auth switch sub-account

# 현재 상태
$ freee auth status
Profile:     default
User:        taro@example.com
Company:     株式会社サンプル (ID: 1234567)
Token:       Valid (expires in 45m)

# 토큰 출력 (파이프/스크립트용)
$ freee auth token
eyJhbGci...
```

### 사업소 (Company)

```bash
# 사업소 목록
$ freee company list
ID        NAME                    ROLE
1234567   株式会社サンプル         admin
2345678   合同会社テスト           member

# 사업소 전환
$ freee company switch 2345678
✓ Switched to 合同会社テスト
```

### 거래 (Deals)

```bash
# 거래 목록 (기본: 이번 달)
$ freee deal list
ID       DATE        TYPE     PARTNER          AMOUNT      STATUS
123456   2026-03-01  income   株式会社A        100,000     settled
123457   2026-03-05  expense  B商事            50,000      unsettled

# 필터
$ freee deal list --type expense --partner "B商事" --from 2026-01-01 --to 2026-03-31

# 상세
$ freee deal show 123456

# JSON 출력 (agent-friendly)
$ freee deal list --format json

# 거래 등록
$ freee deal create \
  --type expense \
  --partner "B商事" \
  --date 2026-03-10 \
  --account "旅費交通費" \
  --amount 5000 \
  --tax-code 1
```

### 청구서 (Invoices)

```bash
$ freee invoice list --status draft
$ freee invoice create \
  --partner "株式会社A" \
  --title "3月分請求書" \
  --item "コンサルティング費" --amount 500000 \
  --due-date 2026-04-30
```

### 글로벌 플래그

```bash
# 출력 형식 지정
$ freee deal list --format json    # JSON (기본 agent 모드)
$ freee deal list --format yaml    # YAML
$ freee deal list --format csv     # CSV
$ freee deal list --format table   # 테이블 (기본 human 모드)

# 프로필 지정
$ freee deal list --profile sub-account

# 사업소 ID 직접 지정
$ freee deal list --company-id 1234567

# 비대화형 모드
$ freee deal list --no-input

# 디버그
$ freee deal list --verbose
$ FREEE_DEBUG=api freee deal list

# 조용한 모드
$ freee deal list --quiet
```

### CI/CD 활용 예

```bash
# 환경변수로 인증 (CI용)
export FREEE_TOKEN="eyJhbGci..."
export FREEE_COMPANY_ID="1234567"

# 이번 달 경비 합계를 JSON으로 추출
freee deal list --type expense --format json | jq '[.[].amount] | add'

# 청구서 일괄 생성 (스크립트)
cat partners.json | jq -r '.[] | .id' | while read id; do
  freee invoice create --partner-id "$id" --title "月次請求" --amount 100000
done
```

---

## Output Contract

- `--format json`: stdout = JSON only, stderr = progress/warnings
- `--format table`: stdout = table, stderr = warnings/notices
- 모든 list 명령은 배열 반환 (빈 결과 = `[]`)
- 모든 show 명령은 객체 반환
- 에러 메시지는 항상 stderr

---

## 진행 현황 요약

| Phase | 범위 | 명령어 수 | 완료 | 진행률 |
|-------|------|-----------|------|--------|
| Phase 1 | 기반 + Accounting | 52 | 0 | 0% |
| Phase 2 | HR + Invoice | 6 | 0 | 0% |
| Phase 3 | Time Tracking + Sales | 5 | 0 | 0% |
| **합계** | | **63** | **0** | **0%** |

---

## 구현 순서 (Phase 1 내)

1. **프로젝트 초기화**: `go mod init`, Makefile, .goreleaser.yaml, cobra root
2. **internal/config**: 프로필, 설정 파일 관리
3. **internal/api**: HTTP 클라이언트, OAuth2 인증 플로우
4. **internal/output**: JSON/YAML/CSV/Table 포매터
5. **internal/errors**: 에러 타입, exit code
6. **cmd/auth**: login, logout, status, list, switch, remove, token
7. **cmd/company**: list, show, switch
8. **cmd/deal**: list, show, create, update, delete
9. **cmd/partner**: list, show, create, update, delete
10. **cmd/invoice**: list, show, create, update, delete
11. **cmd/account, section, tag, item**: CRUD
12. **cmd/journal, expense, walletable**: 나머지 Accounting
13. **cmd/version, completion**: 유틸리티
14. **테스트, 문서, 릴리스**

---

## Version History

| Version | Date | Description |
|---------|------|-------------|
| 0.1.1 | TBD | 出力改善、型付きモデル導入 |
| 0.1.0 | 2026-03-10 | Phase 1 — 기반 + Accounting API |

### 0.1.1 Changes

#### 問題点

freee API のレスポンスはラッパーオブジェクトで包まれている（例: `{"company": {...}}`、`{"account_items": [...]}`）。
v0.1.0 では `var resp any` で受け取っていたため：

- `--format table`（デフォルト）で `map[company:map[...]]` のような Go 内部表現が表示される
- 数値が `1.2261605e+07` のような指数表記になる
- list 系コマンドでテーブルヘッダーが出ない

#### 修正方針

1. **型付きモデル導入** (`internal/model/`): API レスポンスに対応する Go 構造体を定義
2. **show コマンドの key-value 出力**: `--format table` 時は `key: value` 形式で読みやすく表示
3. **list コマンドのテーブル出力**: ラッパーを剥がして配列部分のみフォーマッタに渡す

#### 対象コマンド

| コマンド | 現状 | 修正後 |
|----------|------|--------|
| `company show` | `map[company:map[...]]` | key-value 形式（ID, Name, Role, Phone, Address 等） |
| `company list` | 動作中（既に型あり） | 変更なし |
| `account list` | `map[account_items:[...]]` | テーブル（ID, Name, Category, Tax Code） |
| `item list` | `map[items:[...]]` | テーブル（ID, Name, Available） |
| `section list` | `map[sections:[...]]` | テーブル（ID, Name） |
| `tag list` | `map[tags:[...]]` | テーブル（ID, Name） |
| `walletable list` | `map[walletables:[...]]` | テーブル（ID, Type, Name） |
| `partner list` | `map[partners:[...]]` | テーブル（ID, Name, Code） |
| `deal list` | `map[deals:[...]]` | テーブル（ID, Date, Type, Amount, Status） |
| `invoice list` | `map[invoices:[...]]` | テーブル（ID, Number, Partner, Amount, Status） |
| `expense list` | `map[expense_applications:[...]]` | テーブル（ID, Title, Amount, Status） |

#### 新規ファイル

| ファイル | 内容 |
|----------|------|
| `internal/model/company.go` | Company 構造体 |
| `internal/model/account.go` | AccountItem 構造体 |
| `internal/model/item.go` | Item 構造体 |
| `internal/model/section.go` | Section 構造体 |
| `internal/model/tag.go` | Tag 構造体 |
| `internal/model/walletable.go` | Walletable 構造体 |
| `internal/model/partner.go` | Partner 構造体 |
| `internal/model/deal.go` | Deal 構造体 |
| `internal/model/invoice.go` | Invoice 構造体 |
| `internal/model/expense.go` | ExpenseApplication 構造体 |

#### `company show` の出力例

**修正前**:
```
map[company:map[amount_fraction:0 company_number:7305785747 ...]]
```

**修正後** (`--format table`、デフォルト):
```
ID:          12261605
Name:        姜文喜
Display:     姜 文喜
Name Kana:   カンムンヒ
Role:        admin
Phone:       080-4209-1342
Zipcode:     154-0002
Address:     世田谷区下馬3-23-21 プレシス三軒茶屋204
Fiscal Year: 2025-01-01 ~ 2025-12-31
```

**修正後** (`--format json`):
```json
{"company": {"id": 12261605, "name": "姜文喜", ...}}
```
（API レスポンスをそのまま出力）

#### `account list` の出力例

**修正前**:
```
map[account_items:[map[account_category:事業主借 ...] ...]]
```

**修正後**:
```
ID           NAME                        CATEGORY    TAX_CODE
983573891    [確]一時収入                  事業主借      2
983573893    [確]保険金補填（医療費）        事業主借      2
```

#### 修正ファイル

| ファイル | 変更内容 |
|----------|----------|
| `cmd/company/company.go` | show で型付きモデル使用、key-value 出力 |
| `cmd/account/account.go` | list/show でラッパー剥がし＋行構造体 |
| `cmd/item/item.go` | 同上 |
| `cmd/section/section.go` | 同上 |
| `cmd/tag/tag.go` | 同上 |
| `cmd/walletable/walletable.go` | 同上 |
| `cmd/partner/partner.go` | 同上 |
| `cmd/deal/deal.go` | 同上 |
| `cmd/invoice/invoice.go` | 同上 |
| `cmd/expense/expense.go` | 同上 |
