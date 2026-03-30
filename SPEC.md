# freee CLI — 설계 명세서

**Version**: 0.3.0
**Status**: Phase 1 완료 / Phase 2 개발 중
**Binary**: `freee`
**Repository**: [planitaicojp/freee-cli](https://github.com/planitaicojp/freee-cli)
**Language**: Go 1.26+

---

## 목차

1. [배경과 목적](#1-배경과-목적)
2. [설계 결정 (ADR)](#2-설계-결정-adr)
3. [빠른 시작](#3-빠른-시작)
4. [인증](#4-인증)
5. [글로벌 플래그](#5-글로벌-플래그)
6. [출력 명세](#6-출력-명세)
7. [에러 카탈로그](#7-에러-카탈로그)
8. [커맨드 레퍼런스](#8-커맨드-레퍼런스)
9. [워크플로우 레시피](#9-워크플로우-레시피)
10. [설정 레퍼런스](#10-설정-레퍼런스)
11. [진행 현황](#11-진행-현황)
12. [로드맵](#12-로드맵)
13. [버전 이력](#13-버전-이력)

---

## 1. 배경과 목적

### 과제

freee는 일본 최대의 클라우드 회계/HR SaaS로, 5종의 공개 API를 제공한다. 그러나:

- **공식 SDK는 사실상 archived** (Java, JS, C# — PHP만 유지)
- **CLI 도구가 존재하지 않음** — curl을 직접 작성하거나 GUI를 사용하는 수밖에 없음
- **월별 처리가 수작업** — 분개 확인, 경비 승인, 리포트 출력을 모두 브라우저에서 수행
- **CI/CD에 통합 불가** — 자동화에는 독자적인 스크립트가 필요

### 해결책

`freee` CLI는 freee API를 **셸 스크립트·CI/CD·AI 에이전트**에서 조작할 수 있는 단일 실행 파일을 제공한다.

```bash
# 이번 달 미결제 거래를 CSV로 경리 담당자에게 전달
freee deal list --status unsettled --format csv > unsettled-$(date +%Y%m).csv

# 승인 대기 경비신청을 일괄 승인 (상급자가 실행)
freee expense list --status pending --format json | \
  jq -r '.[].id' | \
  xargs -I{} freee expense action {} --action approve

# 월별 P&L을 Slack에 게시 (CI cron)
freee report pl --format json | jq '.net_income' | \
  curl -s -X POST $SLACK_WEBHOOK -d "{\"text\":\"이번 달 순이익: $(cat -)엔\"}"
```

### freee-mcp와의 관계

freee는 공식 MCP 서버(`freee/freee-mcp`)를 제공하지만, 이는 **AI 에이전트 전용**으로 셸 스크립트나 CI에서는 사용할 수 없다. 본 CLI는 전통적인 셸 도구로서 양쪽을 보완한다.

---

## 2. 설계 결정 (ADR)

### ADR-1: SDK를 사용하지 않음

| 항목 | 내용 |
|------|------|
| **결정** | 공식 SDK·커뮤니티 SDK를 사용하지 않고, OpenAPI 스펙을 참조하여 `net/http`로 직접 호출 |
| **이유** | 공식 SDK는 archived. `LayerXcom/freee-go`는 API 변경에 대한 추적이 늦을 수 있음 |
| **결과** | API 변경에 즉시 대응 가능, 의존성 제로, 바이너리 크기 최소화 |

### ADR-2: OAuth2 Authorization Code + PKCE

| 항목 | 내용 |
|------|------|
| **결정** | 브라우저 기반의 OAuth2 PKCE 플로우 |
| **이유** | freee는 Client Credentials를 제공하지 않음. CI용은 `FREEE_TOKEN` 환경 변수로 대체 |
| **주의점** | `scope=read write` 필수. `prompt=select_company`는 사용 금지(에러 원인). `AuthStyleInParams` 필수 |

### ADR-3: 회계 기간(Fiscal Year)을 1등급 개념으로 취급

| 항목 | 내용 |
|------|------|
| **결정** | `--fiscal-year <YYYY>` 글로벌 플래그를 제공. 기본값은 현재 회계 기간 |
| **이유** | freee의 모든 데이터는 회계 기간에 귀속됨. 복수 기간에 걸친 처리는 일상적으로 발생 |
| **예시** | `freee deal list --fiscal-year 2025`로 전기 거래 참조 |

### ADR-4: 이름 해결 (Name Resolution)

| 항목 | 내용 |
|------|------|
| **결정** | ID 플래그와 동등한 `--<resource>-name` 플래그를 제공하여 내부에서 ID로 변환 |
| **이유** | 5년을 써도 거래처 ID는 외울 수 없음. `--partner-id 12345`보다 `--partner-name "주식회사A"`가 실용적 |
| **구현** | list API로 완전 일치 우선 검색 (대소문자 무시), 0건일 경우 부분 일치(contains) fallback. 1건이면 사용, 복수 건이면 에러로 선택을 유도 |

### ADR-5: 출력 스트림 분리

| 항목 | 내용 |
|------|------|
| **결정** | 데이터는 stdout, 진행 상황/경고/에러는 stderr에 출력 |
| **이유** | 파이프 처리에서 `| jq`가 항상 안전하게 동작함을 보장 |

---

## 3. 빠른 시작

### 설치

```bash
# Homebrew (macOS/Linux)
brew install planitaicojp/tap/freee

# Go install
go install github.com/planitaicojp/freee-cli@latest

# 바이너리 직접 다운로드 (GitHub Releases)
curl -fsSL https://github.com/planitaicojp/freee-cli/releases/latest/download/freee_linux_amd64.tar.gz | tar xz
```

### 5분 만에 첫 결과 얻기

```bash
# 1. 로그인 (브라우저가 열림)
freee auth login

# 2. 사업소 확인
freee company list

# 3. 이번 달 거래 목록
freee deal list

# 4. JSON으로 가져와 파이프 처리
freee deal list --format json | jq '[.[].amount] | add'

# 5. 경비신청 등록
freee expense create \
  --title "3월 출장 교통비" \
  --amount 12500 \
  --date 2026-03-15 \
  --account-item-id 1234
```

---

## 4. 인증

### 앱 사전 준비

1. [freee 앱 스토어](https://app.secure.freee.co.jp/developers)에서 앱을 등록
2. Redirect URI에 `http://localhost:8080/callback`을 추가
3. Client ID / Client Secret 취득 → `freee auth login` 실행 시 입력

### 브라우저 로그인

```bash
# 기본 프로필로 로그인
freee auth login

# 다른 계정을 추가
freee auth login --profile corp-b

# 로그인 상태 확인
freee auth status
#   Profile:   default
#   User:      taro@example.com
#   Company:   주식회사 샘플 (ID: 1234567)
#   Token:     Valid (expires in 45m)
#   Fiscal:    2026 (2026-01-01 ~ 2026-12-31)

# 프로필 목록
freee auth list
#   PROFILE      USER                  COMPANY              STATUS
# ✓ default      taro@example.com      주식회사 샘플         Active
#   corp-b       hanako@example.com    합동회사 테스트        Active

# 프로필 전환
freee auth switch corp-b

# 액세스 토큰만 출력 (스크립트용)
freee auth token

# 로그아웃
freee auth logout
```

### CI/CD 환경에서의 인증

브라우저를 사용할 수 없는 환경에서는 환경 변수로 인증한다:

```bash
# GitHub Actions 예시
env:
  FREEE_TOKEN: ${{ secrets.FREEE_ACCESS_TOKEN }}
  FREEE_COMPANY_ID: ${{ secrets.FREEE_COMPANY_ID }}

# 또는 임시 토큰을 생성하여 export
export FREEE_TOKEN=$(freee auth token)
export FREEE_COMPANY_ID=1234567
```

**토큰 관리 주의사항:**
- 액세스 토큰의 유효 기간은 통상 24시간
- CI에서는 secrets으로 관리하며 로그에 출력하지 않음 (`--quiet` 사용)
- Refresh token은 `~/.config/freee/credentials.yaml` (권한 0600)에 저장됨

---

## 5. 글로벌 플래그

모든 커맨드에서 사용 가능:

| 플래그 | 환경 변수 | 기본값 | 설명 |
|--------|----------|--------|------|
| `--profile <name>` | `FREEE_PROFILE` | `default` | 사용할 프로필 |
| `--company-id <id>` | `FREEE_COMPANY_ID` | 프로필 설정값 | 사업소 ID 오버라이드 |
| `--fiscal-year <YYYY>` | `FREEE_FISCAL_YEAR` | 현재 회계 기간 | 회계 기간 지정 |
| `--format <fmt>` | `FREEE_FORMAT` | `table` | 출력 형식: `table` / `json` / `yaml` / `csv` |
| `--no-header` | — | false | 테이블 헤더 비표시 (스크립트용) |
| `--quiet` / `-q` | — | false | 성공 시 출력 억제 |
| `--verbose` | `FREEE_DEBUG=1` | false | HTTP 요청/응답 표시 |
| `--no-input` | `FREEE_NO_INPUT=1` | false | 비대화형 모드 (확인 프롬프트 출력 안 함) |
| `--dry-run` | — | false | 변경계 커맨드에서 API 호출 없이 요청 내용을 표시 |

### 회계 기간 플래그 사용 예

```bash
# 전기 (2025년도)의 거래 참조
freee deal list --fiscal-year 2025

# 전기의 P&L 리포트
freee report pl --fiscal-year 2025 --format json

# 기본 회계 기간 변경 (config에 영구 저장)
freee config set fiscal-year 2025
```

### 변경계 커맨드의 --dry-run

```bash
# 실제로는 생성하지 않고, 전송 예정인 JSON을 확인
freee deal create --type expense --date 2026-03-10 \
  --account-item-id 1234 --amount 5000 --dry-run

# [dry-run] POST /api/1/deals
# {
#   "company_id": 1234567,
#   "type": "expense",
#   "issue_date": "2026-03-10",
#   "details": [{"account_item_id": 1234, "amount": 5000}]
# }
```

---

## 6. 출력 명세

### stdout / stderr 분리

| 스트림 | 내용 |
|--------|------|
| `stdout` | 데이터 (JSON, Table, CSV, YAML) |
| `stderr` | 진행 상황, 경고, 에러 메시지, `[dry-run]` 프리픽스 |

### 형식별 보장

| 형식 | list 반환값 | show 반환값 | 빈 목록 |
|------|-----------|-----------|--------|
| `table` | 헤더 + 행 | key-value 형식 | `(no results)` |
| `json` | `[...]` 배열 | `{...}` 객체 | `[]` |
| `csv` | 헤더 + 행 | — | 헤더만 |
| `yaml` | `- ...` 리스트 | `key: value` | `[]` |

### JSON 출력 보장

- `--format json` 시 **항상 유효한 JSON만** stdout에 출력한다
- API 응답의 래퍼는 **제거**한다 (예: `{"deals": [...]}` → `[...]`)
- 숫자는 숫자형 (지수 표기 없음)
- `jq`로 직접 파이프할 수 있음을 보장한다

```bash
# 항상 안전하게 파이프 가능
freee deal list --format json | jq '.[0].id'
freee partner show 12345 --format json | jq '.name'
```

### 테이블 출력 표시 규칙

- 금액 열: 콤마 구분 (`1,234,567`)
- 날짜 열: `YYYY-MM-DD` 형식
- 상태: 영어 그대로 출력 (기계 처리를 위해)
- 긴 문자열은 말줄임표로 생략 (`...`)

---

## 7. 에러 카탈로그

### Exit Code

| Code | 의미 | 예시 |
|------|------|------|
| `0` | 성공 | — |
| `1` | 일반 에러 | 잘못된 인자 |
| `2` | 인증 에러 | 토큰 만료, 미로그인 |
| `3` | Not Found | 지정한 ID의 리소스가 존재하지 않음 |
| `4` | 유효성 검사 에러 | 필수 플래그 미지정, 잘못된 날짜 형식 |
| `5` | API 에러 | freee 서버가 에러를 반환함 |
| `6` | 네트워크 에러 | 연결 실패, 타임아웃 |
| `10` | 취소됨 | Ctrl+C |

### 에러 메시지 형식

```
error: <메시지>
hint:  <다음에 취할 액션>
```

예시:
```
error: authentication required
hint:  run 'freee auth login' to authenticate

error: deal 99999 not found
hint:  run 'freee deal list' to see available deals

error: invalid date format "2026/03/10"
hint:  use YYYY-MM-DD format (e.g. 2026-03-10)

error: API error [invalid_param]: partner_id is required
hint:  run 'freee partner list' to find your partner ID
```

### 레이트 제한

freee API에는 **1,000요청/10분** 제한이 있다.

```
error: rate limit exceeded
hint:  waiting 47s (Retry-After header), then retrying...
```

CLI는 최대 3회 자동으로 재시도한다. `--all` 플래그 사용 시 특히 주의:

```bash
# 대량 데이터 취득 시 --limit으로 페이지당 건수를 조정
freee deal list --all --limit 100  # 기본값: 50건/페이지
```

---

## 8. 커맨드 레퍼런스

### 범례

- `<필수 인자>` / `[생략 가능 인자]`
- 플래그는 모든 커맨드에서 글로벌 플래그와 함께 사용 가능
- ID 대신 이름으로 지정할 수 있는 경우 `--<resource>-name` 플래그가 존재함

---

### 8.1 인증 (`auth`)

```bash
freee auth login [--profile <name>]        # OAuth2 브라우저 로그인
freee auth logout [--profile <name>]       # 토큰 삭제
freee auth status                          # 현재 인증 상태
freee auth list                            # 전체 프로필 목록
freee auth switch <profile>                # 프로필 전환
freee auth remove <profile>                # 프로필 삭제
freee auth token                           # 액세스 토큰 출력 (파이프용)
```

---

### 8.2 설정 (`config`)

```bash
freee config show                          # 전체 설정 표시
freee config set <key> <value>             # 설정 변경
freee config path                          # 설정 파일 경로 표시
```

설정 가능한 키:

| 키 | 설명 | 기본값 |
|----|------|--------|
| `default-profile` | 기본 프로필 | `default` |
| `default-format` | 기본 출력 형식 | `table` |
| `fiscal-year` | 기본 회계 기간 | 현재 기간 |
| `timeout` | API 타임아웃 (초) | `30` |
| `max-retries` | 최대 재시도 횟수 | `3` |

---

### 8.3 사업소 (`company`)

```bash
freee company list                         # 사업소 목록
freee company show [<id>]                  # 사업소 상세 (생략 시 현재 사업소)
freee company switch <id>                  # 기본 사업소 전환
```

```bash
# 출력 예시
$ freee company show
ID:          1234567
Name:        주식회사 샘플
Role:        admin
Phone:       03-1234-5678
Address:     도쿄도 시부야구...
Fiscal Year: 2026-01-01 ~ 2026-12-31
```

---

### 8.4 거래 (`deal`)

```bash
freee deal list [flags]
freee deal show <id>
freee deal create [flags]
freee deal update <id> [flags]
freee deal delete <id>
```

**list 플래그:**

| 플래그 | 설명 | 예시 |
|--------|------|------|
| `--type` | `income` 또는 `expense` | `--type expense` |
| `--status` | `settled` / `unsettled` | `--status unsettled` |
| `--from` | 시작일 (YYYY-MM-DD) | `--from 2026-01-01` |
| `--to` | 종료일 (YYYY-MM-DD) | `--to 2026-03-31` |
| `--partner-id` | 거래처 ID | `--partner-id 12345` |
| `--partner-name` | 거래처명 (이름 해결) | `--partner-name "주식회사A"` |
| `--all` | 전체 취득 (자동 페이지네이션) | `--all` |
| `--limit` | 페이지당 건수 | `--limit 100` |

**create/update 플래그:**

| 플래그 | 설명 | 필수 여부 |
|--------|------|----------|
| `--type` | `income` / `expense` | create만 |
| `--date` | 발생일 (YYYY-MM-DD) | create만 |
| `--partner-id` | 거래처 ID | — |
| `--partner-name` | 거래처명 (이름 해결) | — |
| `--account-item-id` | 계정과목 ID | create만 |
| `--amount` | 금액 (엔) | create만 |
| `--tax-code` | 세금 구분 코드 | — |

```bash
# 이번 달 미결제 지출 전체를 CSV로 취득
freee deal list --type expense --status unsettled \
  --from 2026-03-01 --to 2026-03-31 \
  --all --format csv > march-unsettled.csv

# 거래 등록 (--dry-run으로 사전 확인)
freee deal create \
  --type expense \
  --date 2026-03-15 \
  --partner-name "B상사" \
  --account-item-id 1234 \
  --amount 50000 \
  --tax-code 1 \
  --dry-run

# 거래 ID를 취득하여 이어서 처리
DEAL_ID=$(freee deal create ... --format json --quiet | jq -r '.deal.id')
freee deal payment create $DEAL_ID --amount 50000 --date 2026-03-20
```

---

### 8.5 청구서 (`invoice`)

```bash
freee invoice list [flags]
freee invoice show <id>
freee invoice create [flags]
freee invoice update <id> [flags]
freee invoice delete <id>
```

**list 플래그:**

| 플래그 | 설명 |
|--------|------|
| `--status` | `draft` / `applying` / `approved` / `remanded` / `rejected` / `canceled` / `sending` / `sent` / `overdue` / `settled` |
| `--partner-id` | 거래처 ID |
| `--partner-name` | 거래처명 (이름 해결) |
| `--from` | 청구일 시작 |
| `--to` | 청구일 종료 |
| `--all` | 전체 취득 |

```bash
# 기한 초과 미수금 청구서 목록
freee invoice list --status overdue --format json | \
  jq '.[] | {id, partner: .partner_name, amount: .total_amount, due: .payment_date}'

# 청구서 생성
freee invoice create \
  --partner-name "주식회사A" \
  --title "2026년 3월분 컨설팅 비용" \
  --date 2026-03-31 \
  --due-date 2026-04-30 \
  --item "컨설팅" \
  --amount 500000
```

---

### 8.6 거래처 (`partner`)

```bash
freee partner list [--all] [--format <fmt>]
freee partner show <id>
freee partner create [flags]
freee partner update <id> [flags]
freee partner delete <id>
```

**create/update 플래그:**

| 플래그 | 설명 |
|--------|------|
| `--name` | 거래처명 (필수) |
| `--code` | 거래처 코드 |
| `--email` | 이메일 주소 |
| `--phone` | 전화번호 |
| `--contact-name` | 담당자명 |

```bash
# 거래처를 이름으로 검색
freee partner list --format json | jq '.[] | select(.name | contains("샘플"))'

# 거래처 등록
freee partner create --name "신규 주식회사" --code "SHINKI001" --email "info@shinki.co.jp"
```

---

### 8.7 계정과목 (`account`)

```bash
freee account list [--all] [--format <fmt>]
freee account show <id>
```

```bash
# 계정과목을 이름으로 검색
freee account list --format json | jq '.[] | select(.name | contains("교통"))'
```

---

### 8.8 부문 (`section`)

```bash
freee section list [--format <fmt>]
freee section create --name <name> [--shortcut1 <s>] [--shortcut2 <s>]
freee section update <id> [--name <name>] [--shortcut1 <s>] [--shortcut2 <s>]
freee section delete <id>
```

---

### 8.9 메모 태그 (`tag`)

```bash
freee tag list [--format <fmt>]
freee tag create --name <name> [--shortcut1 <s>] [--shortcut2 <s>]
freee tag update <id> [--name <name>]
freee tag delete <id>
```

---

### 8.10 품목 (`item`)

```bash
freee item list [--format <fmt>]
freee item create --name <name>
freee item update <id> [--name <name>]
freee item delete <id>
```

---

### 8.11 분개 (`journal`)

```bash
freee journal list [flags]           # 분개를 CSV/JSON으로 다운로드
```

| 플래그 | 설명 |
|--------|------|
| `--from` | 시작일 |
| `--to` | 종료일 |
| `--format` | `json` / `csv` |

---

### 8.12 경비신청 (`expense`)

```bash
freee expense list [flags]
freee expense show <id>
freee expense create [flags]
freee expense update <id> [flags]
freee expense delete <id>
freee expense action <id> --action <approve|reject|cancel>  # Phase 4
```

**create/update 플래그:**

| 플래그 | 설명 | 필수 여부 |
|--------|------|----------|
| `--title` | 신청 제목 | create만 |
| `--amount` | 금액 | create만 |
| `--date` | 발생일 | create만 |
| `--account-item-id` | 계정과목 ID | create만 |
| `--description` | 적요 | — |

```bash
# 내 승인 대기 경비신청 확인
freee expense list --status pending --format table

# 경비신청 등록 (영수증은 별도로 freee receipt create로 첨부)
freee expense create \
  --title "도쿄 출장 교통비" \
  --amount 12500 \
  --date 2026-03-15 \
  --account-item-id 1234
```

---

### 8.13 구좌 (`walletable`)

```bash
freee walletable list [--format <fmt>]
freee walletable show <type> <id>
```

구좌 타입: `bank_account` / `credit_card` / `wallet` / `other`

---

### 8.14 유틸리티

```bash
freee version               # 버전 표시
freee completion bash        # bash 자동완성 스크립트 생성
freee completion zsh         # zsh 자동완성 스크립트 생성
freee completion fish        # fish 자동완성 스크립트 생성
```

---

## 9. 워크플로우 레시피

### 9.1 월별 마감 처리

월말에 경리 담당자가 실시하는 일반적인 절차:

```bash
#!/bin/bash
# monthly-close.sh — 월별 마감 보조 스크립트
YEAR=2026
MONTH=03
FROM="${YEAR}-${MONTH}-01"
TO="${YEAR}-${MONTH}-31"

echo "=== 미결제 거래 확인 ==="
freee deal list --status unsettled --from $FROM --to $TO --format table

echo ""
echo "=== 승인 대기 경비신청 ==="
freee expense list --status pending --format table

echo ""
echo "=== 이번 달 P&L 요약 ==="
freee report pl --from $FROM --to $TO --format json | \
  jq '{매출액: .revenue, 비용: .expenses, 순이익: .net_income}'
```

### 9.2 거래처별 월별 청구

```bash
#!/bin/bash
# monthly-invoice.sh — 거래처 목록에서 월별 청구서 일괄 생성
MONTH_LABEL="2026년 3월분"
DUE_DATE="2026-04-30"

freee partner list --format json | jq -r '.[] | select(.name | startswith("거래처")) | .id' | \
while read partner_id; do
  echo "Creating invoice for partner $partner_id..."
  freee invoice create \
    --partner-id "$partner_id" \
    --title "${MONTH_LABEL} 월별 서비스 비용" \
    --date "2026-03-31" \
    --due-date "$DUE_DATE" \
    --amount 100000 \
    --quiet
done
echo "Done."
```

### 9.3 CI/CD에서의 월별 리포트 발송

GitHub Actions로 매월 말에 P&L을 Slack에 알리는 예시:

```yaml
# .github/workflows/monthly-report.yml
name: Monthly P&L Report

on:
  schedule:
    - cron: '0 0 1 * *'   # 매월 1일 09:00 JST

jobs:
  report:
    runs-on: ubuntu-latest
    steps:
      - name: Install freee CLI
        run: |
          curl -fsSL https://github.com/planitaicojp/freee-cli/releases/latest/download/freee_linux_amd64.tar.gz | tar xz
          sudo mv freee /usr/local/bin/

      - name: Generate report
        env:
          FREEE_TOKEN: ${{ secrets.FREEE_ACCESS_TOKEN }}
          FREEE_COMPANY_ID: ${{ secrets.FREEE_COMPANY_ID }}
        run: |
          LAST_MONTH=$(date -d "last month" +%Y-%m)
          FROM="${LAST_MONTH}-01"
          TO=$(date -d "${FROM} +1 month -1 day" +%Y-%m-%d)

          NET_INCOME=$(freee report pl \
            --from "$FROM" --to "$TO" \
            --format json | jq '.net_income')

          curl -s -X POST ${{ secrets.SLACK_WEBHOOK }} \
            -H 'Content-type: application/json' \
            -d "{\"text\":\"📊 ${LAST_MONTH} 월별 순이익: ${NET_INCOME} 엔\"}"
```

### 9.4 Line Messaging API와의 연동

Line Pay 결제 → freee에 거래를 자동 등록하는 패턴:

```bash
# line-pay-webhook.sh — Line Pay webhook에서 호출되는 스크립트
# 인자: $1=amount $2=partner_name $3=date
AMOUNT=$1
PARTNER_NAME=$2
DATE=$3

# 거래처 ID를 이름으로 해결
PARTNER_ID=$(freee partner list --format json | \
  jq -r ".[] | select(.name == \"$PARTNER_NAME\") | .id")

if [ -z "$PARTNER_ID" ]; then
  echo "거래처 '$PARTNER_NAME'를 찾을 수 없습니다. 신규 등록합니다..."
  PARTNER_ID=$(freee partner create \
    --name "$PARTNER_NAME" \
    --format json | jq -r '.partner.id')
fi

# 거래 등록
DEAL_ID=$(freee deal create \
  --type income \
  --date "$DATE" \
  --partner-id "$PARTNER_ID" \
  --account-item-id "$SALES_ACCOUNT_ID" \
  --amount "$AMOUNT" \
  --format json | jq -r '.deal.id')

echo "거래 #$DEAL_ID 등록 완료 (¥$AMOUNT / $PARTNER_NAME)"
```

### 9.5 JSON 입력을 활용한 파이프라인 통합

외부 시스템의 JSON 데이터를 freee에 직접 흘려보내는 패턴:

```bash
# 외부 API의 주문 데이터를 거래로 일괄 등록
curl -s https://api.external-shop.example/orders?date=2026-03-15 | \
  jq -r '.orders[] | [.amount, .customer, .date] | @tsv' | \
  while IFS=$'\t' read amount customer date; do
    freee deal create \
      --type income \
      --date "$date" \
      --partner-name "$customer" \
      --account-item-id 5678 \
      --amount "$amount" \
      --no-input --quiet
  done
```

---

## 10. 설정 레퍼런스

### 파일 구성

```
~/.config/freee/
├── config.yaml       # 프로필 설정 (0600)
├── credentials.yaml  # OAuth 토큰 (0600)
└── tokens.yaml       # 토큰 캐시 (0600)
```

### config.yaml 포맷

```yaml
default_profile: default
default_format: table
fiscal_year: 2026          # 생략 시 현재 회계 기간
timeout: 30                # 초
max_retries: 3

profiles:
  default:
    company_id: 1234567
    user: taro@example.com

  corp-b:
    company_id: 2345678
    user: hanako@example.com
```

### 환경 변수 목록

| 변수 | 설명 | 우선도 |
|------|------|--------|
| `FREEE_TOKEN` | 액세스 토큰 직접 지정 (CI용) | 최고 |
| `FREEE_COMPANY_ID` | 사업소 ID 오버라이드 | 높음 |
| `FREEE_PROFILE` | 사용할 프로필 | 높음 |
| `FREEE_FISCAL_YEAR` | 회계 기간 | 높음 |
| `FREEE_FORMAT` | 기본 출력 형식 | 중간 |
| `FREEE_NO_INPUT` | 비대화형 모드 (`1`로 활성화) | 중간 |
| `FREEE_DEBUG` | 디버그 출력 (`1`=verbose, `api`=HTTP 상세) | 낮음 |
| `FREEE_CONFIG_DIR` | 설정 디렉터리 경로 | 낮음 |

---

## 11. 진행 현황

### Phase 1: 기반 + Accounting API 핵심 (52개 커맨드)

| # | 커맨드 | 설명 | 상태 |
|---|--------|------|------|
| **인증/설정** | | | |
| 1 | `freee auth login` | OAuth2 브라우저 로그인 | ☑ |
| 2 | `freee auth logout` | 토큰 삭제 | ☑ |
| 3 | `freee auth status` | 인증 상태 표시 | ☑ |
| 4 | `freee auth list` | 프로필 목록 | ☑ |
| 5 | `freee auth switch <profile>` | 프로필 전환 | ☑ |
| 6 | `freee auth remove <profile>` | 프로필 삭제 | ☑ |
| 7 | `freee auth token` | 액세스 토큰 출력 | ☑ |
| 8 | `freee config show` | 설정 표시 | ☑ |
| 9 | `freee config set <key> <value>` | 설정 변경 | ☑ |
| 10 | `freee config path` | 설정 파일 경로 출력 | ☑ |
| **사업소** | | | |
| 11 | `freee company list` | 사업소 목록 | ☑ |
| 12 | `freee company show [id]` | 사업소 상세 | ☑ |
| 13 | `freee company switch <id>` | 기본 사업소 전환 | ☑ |
| **거래 (Deals)** | | | |
| 14 | `freee deal list` | 거래 목록 | ☑ |
| 15 | `freee deal show <id>` | 거래 상세 | ☑ |
| 16 | `freee deal create` | 거래 등록 | ☑ |
| 17 | `freee deal update <id>` | 거래 수정 | ☑ |
| 18 | `freee deal delete <id>` | 거래 삭제 | ☑ |
| **청구서 (Invoices)** | | | |
| 19 | `freee invoice list` | 청구서 목록 | ☑ |
| 20 | `freee invoice show <id>` | 청구서 상세 | ☑ |
| 21 | `freee invoice create` | 청구서 작성 | ☑ |
| 22 | `freee invoice update <id>` | 청구서 수정 | ☑ |
| 23 | `freee invoice delete <id>` | 청구서 삭제 | ☑ |
| **거래처 (Partners)** | | | |
| 24 | `freee partner list` | 거래처 목록 | ☑ |
| 25 | `freee partner show <id>` | 거래처 상세 | ☑ |
| 26 | `freee partner create` | 거래처 등록 | ☑ |
| 27 | `freee partner update <id>` | 거래처 수정 | ☑ |
| 28 | `freee partner delete <id>` | 거래처 삭제 | ☑ |
| **계정과목 (Account Items)** | | | |
| 29 | `freee account list` | 계정과목 목록 | ☑ |
| 30 | `freee account show <id>` | 계정과목 상세 | ☑ |
| **부문 (Sections)** | | | |
| 31 | `freee section list` | 부문 목록 | ☑ |
| 32 | `freee section create` | 부문 등록 | ☑ |
| 33 | `freee section update <id>` | 부문 수정 | ☑ |
| 34 | `freee section delete <id>` | 부문 삭제 | ☑ |
| **메모 태그 (Tags)** | | | |
| 35 | `freee tag list` | 태그 목록 | ☑ |
| 36 | `freee tag create` | 태그 등록 | ☑ |
| 37 | `freee tag update <id>` | 태그 수정 | ☑ |
| 38 | `freee tag delete <id>` | 태그 삭제 | ☑ |
| **품목 (Items)** | | | |
| 39 | `freee item list` | 품목 목록 | ☑ |
| 40 | `freee item create` | 품목 등록 | ☑ |
| 41 | `freee item update <id>` | 품목 수정 | ☑ |
| 42 | `freee item delete <id>` | 품목 삭제 | ☑ |
| **분개 (Journals)** | | | |
| 43 | `freee journal list` | 분개 다운로드 | ☑ |
| **경비신청 (Expenses)** | | | |
| 44 | `freee expense list` | 경비신청 목록 | ☑ |
| 45 | `freee expense show <id>` | 경비신청 상세 | ☑ |
| 46 | `freee expense create` | 경비신청 작성 | ☑ |
| 47 | `freee expense update <id>` | 경비신청 수정 | ☑ |
| 48 | `freee expense delete <id>` | 경비신청 삭제 | ☑ |
| **구좌 (Walletables)** | | | |
| 49 | `freee walletable list` | 구좌 목록 | ☑ |
| 50 | `freee walletable show <id>` | 구좌 상세 | ☑ |
| **유틸리티** | | | |
| 51 | `freee version` | 버전 표시 | ☑ |
| 52 | `freee completion` | 셸 자동완성 생성 | ☑ |

### Phase 2: Accounting 확장 — 분개장·이체·명세·세금 (19개 커맨드)

| # | 커맨드 | 설명 | 상태 |
|---|--------|------|------|
| 53–57 | `freee manual-journal` CRUD | 분개장 (仕訳帳) | ☐ |
| 58–62 | `freee transfer` CRUD | 이체 전표 (振替伝票) | ☐ |
| 63–66 | `freee wallet-txn` CRUD | 구좌 명세 (口座明細) | ☐ |
| 67–69 | `freee walletable` create/update/delete | 구좌 등록/수정/삭제 | ☐ |
| 70–71 | `freee tax` list/show | 세금 구분 (税区分) | ☐ |

### Phase 3: Accounting 확장 — 리포트 + 증빙 파일함 (25개 커맨드)

| # | 커맨드 | 설명 | 상태 |
|---|--------|------|------|
| 72–86 | `freee report` bs/pl/cr (각 5패턴) | 재무제표 (BS·PL·CF) | ☐ |
| 87–89 | `freee report` general-ledger/trial-bs/trial-pl | 총계정원장·잔액시산표 | ☐ |
| 90–95 | `freee receipt` CRUD + download | 증빙 관리 | ☐ |
| 96 | `freee bank list` | 은행 목록 | ☐ |

### Phase 4: Accounting 확장 — 워크플로우 + 서브리소스 (44개 커맨드)

| # | 커맨드 | 설명 | 상태 |
|---|--------|------|------|
| 97–102 | `freee approval` CRUD + action | 승인 의뢰 | ☐ |
| 103–106 | `freee approval-form/route` | 승인 양식/경로 | ☐ |
| 107–112 | `freee payment-request` CRUD + action | 지급 의뢰 | ☐ |
| 113–117 | `freee expense-template` CRUD | 경비 템플릿 | ☐ |
| 118 | `freee expense action` | 경비신청 승인/반려 | ☐ |
| 119–124 | `freee deal payment/renew` | 거래 결제·갱신 | ☐ |
| 125–126 | `freee quotation` list/show | 견적서 | ☐ |
| 127–130 | `freee segment-tag` CRUD | 세그먼트 태그 | ☐ |
| 131–135 | `freee code` account/section/item/tag/walletable upsert | 각종 코드 | ☐ |
| 136–137 | `freee partner code` create/update | 거래처 코드 | ☐ |
| 138–140 | `freee user` list/me/capabilities | 사용자 | ☐ |

### Phase 5: HR API 전체 (93개 커맨드)

| # | 커맨드 | 설명 | 상태 |
|---|--------|------|------|
| 141–145 | `freee hr-employee` CRUD | 직원 | ☐ |
| 146–150 | `freee hr-work-record` CRUD | 근태 | ☐ |
| 151–154 | `freee hr-time-clock` | 출퇴근 기록 | ☐ |
| 155–159 | `freee hr-attendance-tag` CRUD | 소정 근무 태그 | ☐ |
| 160–163 | `freee hr-salary/bonus` list/show | 급여·상여 명세 | ☐ |
| 164–170 | `freee hr-group/position` CRUD | 그룹·직위 | ☐ |
| 171–187 | `freee hr-employee <sub>` show/update | 직원 서브리소스 | ☐ |
| 188–217 | `freee hr-approval-*` | HR 승인 워크플로우 (4종) | ☐ |
| 218–219 | `freee hr-approval-route` | HR 승인 경로 | ☐ |
| 220 | `freee hr-user me` | HR 사용자 정보 | ☐ |
| 221–233 | `freee hr-yearend` | 연말정산 | ☐ |

### Phase 6: Invoice API + PM API + Sales API (39개 커맨드)

| # | 커맨드 | 설명 | 상태 |
|---|--------|------|------|
| 234–248 | `freee iv-invoice/iv-quotation/iv-delivery/iv-template` | Invoice API | ☐ |
| 249–258 | `freee pm-project/pm-workload/pm-team` | PM (공수 관리) API | ☐ |
| 259–270 | `freee sales-business/sales-order/sales-customer/sales-product` | Sales API | ☐ |
| 271–272 | `freee api-resources`, `freee forms selectables` | 메타 정보 | ☐ |

### 진행 상황 요약

| Phase | 범위 | 커맨드 수 | 완료 | 진척도 |
|-------|------|-----------|------|--------|
| Phase 1 | 기반 + Accounting 핵심 | 52 | **52** | **100%** |
| Phase 2 | Accounting 분개·이체·명세·세금 | 19 | 0 | 0% |
| Phase 3 | Accounting 리포트 + 증빙 파일함 | 25 | 0 | 0% |
| Phase 4 | Accounting 워크플로우 + 서브리소스 | 44 | 0 | 0% |
| Phase 5 | HR API 전체 | 93 | 0 | 0% |
| Phase 6 | Invoice / PM / Sales API | 39 | 0 | 0% |
| **합계** | | **272** | **52** | **19%** |

---

## 12. 로드맵

### v0.4.0 — 이름 해결 + 사용성 개선

5년 사용자의 가장 큰 불만: **ID를 외워야 한다**.

| # | 항목 | 설명 |
|---|------|------|
| 1 | `--partner-name` 이름 해결 | `--partner-id` 대신 이름으로 지정 |
| 2 | `--account-name` 이름 해결 | 계정과목 이름으로 지정 |
| 3 | `--fiscal-year` 글로벌 플래그 구현 | 회계 기간을 1등급 인자로 구현 |
| 4 | 금액 콤마 포맷 | table 모드에서 `1,234,567` 형식 |
| 5 | 상태 한국어 라벨 | `settled` → `결제완료` (table 모드만) |
| 6 | `--no-header` 플래그 | 테이블 헤더 비표시 |
| 7 | `freee schema <resource>` | JSON Schema 출력 (통합 개발용) |

### v0.5.0 — Accounting 분개·이체·명세 (Phase 2, 19개 커맨드)

```bash
# 분개장 조작
freee manual-journal list --from 2026-03-01 --to 2026-03-31
freee manual-journal create --debit-account 매출 --credit-account 현금 --amount 100000

# 이체 전표
freee transfer create ...

# 구좌 명세
freee wallet-txn list --walletable-id 12345 --all
```

### v0.6.0 — 리포트 + 증빙 관리 (Phase 3, 25개 커맨드)

재무제표 출력과 영수증 업로드를 추가. 증빙 관리는 경비신청의 실용성에 직결되므로 워크플로우(Phase 4)보다 먼저 구현한다.

```bash
# 월별 P&L
freee report pl --from 2026-03-01 --to 2026-03-31 --format json

# 영수증 업로드 → 경비신청에 연결
freee receipt create --file receipt.jpg
freee expense update 12345 --receipt-id <receipt_id>
```

### v0.7.0 — 승인 워크플로우 (Phase 4, 44개 커맨드)

월별 업무에서 가장 중요한 승인 플로우.

```bash
# 경비신청 일괄 승인
freee expense list --status pending --format json | \
  jq -r '.[].id' | \
  xargs -I{} freee expense action {} --action approve

# 지급 의뢰 처리
freee payment-request list --status draft
freee payment-request action 12345 --action approve
```

### v0.8.0 — 배치 처리 + 감사 로그

| # | 항목 | 설명 |
|---|------|------|
| 1 | `freee deal import --file txn.csv` | CSV 일괄 임포트 |
| 2 | `freee invoice bulk-send` | 청구서 일괄 발송 |
| 3 | `~/.local/share/freee/audit.log` | 변경계 조작의 감사 로그 |
| 4 | `--idempotency-key` | 멱등 키 지원 (외부 연동 중복 방지) |

### v0.9.0 — HR API (Phase 5, 93개 커맨드)

```bash
freee hr-employee list
freee hr-work-record list --employee-id 123 --from 2026-03-01
freee hr-time-clock create --type clock_in
```

### v1.0.0 — Invoice / PM / Sales API (Phase 6, 39개 커맨드)

전체 272개 커맨드 구현 완료. GoReleaser를 통한 정식 릴리스.

---

## 13. 버전 이력

| 버전 | 날짜 | 내용 |
|------|------|------|
| **v0.4.1** | 2026-03-30 | Table 사용성 (콤마 포맷, 한국어 라벨, --no-header) + --fiscal-year + resolve 리팩토링 |
| v0.4.0 | 2026-03-29 | 이름 해결 (`--partner-name`, `--account-name`) + ADR-4 부분 일치 확장 |
| v0.3.0 | 2026-03-20 | Phase 1 전체 스텁 구현 + lint 에러 전체 수정 |
| v0.2.1 | 2026-03-11 | 유닛 테스트 (api/errors/output) + GitHub Actions CI |
| v0.2.0 | 2026-03-11 | Agent 신뢰성 향상 (`--all` 페이지네이션, Retry-After, 에러 힌트, `--dry-run`) |
| v0.1.1 | 2026-03-10 | list/show 출력 개선 (typed models, key-value 포맷) |
| v0.1.0 | 2026-03-10 | Phase 1 스캐폴딩, OAuth2 인증, 기본 커맨드 골격 |

### v0.4.1 변경 내용

**Table 출력 개선:**
- 금액 콤마 포맷: int64 필드에 `1,234,567` 형식 적용 (ID/코드 필드 제외)
- 한국어 상태 라벨: `settled` → `결제완료`, `draft` → `임시저장` 등 11개 번역
- `--no-header` 글로벌 플래그: table/CSV 모드에서 헤더 행 생략
- `output.New()` 팩토리에 `Options` 변수 추가 (backward-compatible variadic)

**--fiscal-year 글로벌 플래그:**
- 결산월 API 조회 → from/to 자동 계산
- deal list, invoice list, expense list 지원
- invoice list, expense list에 --from/--to 플래그 추가
- --fiscal-year와 --from/--to 상호 배타 (exit 4)

**Resolve 패키지 리팩토링:**
- `Named` 인터페이스 + Go 제네릭으로 중복 제거 (~77줄 삭감)
- 테스트 서버 URL path 검증 추가
- `--account-item-name` alias 테스트 추가

### v0.4.0 변경 내용

**이름 해결 (Name Resolution):**
- `internal/resolve/` 패키지 신설 — `PartnerID()`, `AccountItemID()` 함수
- `--partner-name` 플래그: deal create/update, invoice create/update
- `--account-name` 플래그: deal create/update, expense create (`--account-item-name` alias)
- 매칭 전략: 완전 일치 우선 (대소문자 무시), 부분 일치(contains) fallback
- 에러 처리: mutually exclusive (exit 4), not found (exit 3), multiple match (exit 4 + 후보 목록)
- ADR-4 업데이트: 부분 일치 fallback 전략 반영
- 14개 유닛 테스트 (PartnerID 7개 + AccountItemID 7개)

### v0.3.0 변경 내용

**구현:**
- `deal update/delete` — PUT/DELETE API 호출 구현
- `invoice update/delete` — 동일
- `partner create/update/delete` — POST/PUT/DELETE 구현
- `expense create/update/delete` — 동일
- `section create/update/delete` — 동일
- `tag create/update/delete` — 동일
- `item update/delete` — 동일

**품질:**
- 전체 커맨드의 `fmt.Sscanf`를 `strconv.ParseInt`로 교체 (잘못된 ID 입력 에러 처리 추가)
- `errcheck` lint 위반을 전체 파일에서 해소
- Phase 1 진척도: 39/52 → **52/52 (100%)**

### v0.2.0 변경 내용

| # | 항목 | 상세 |
|---|------|------|
| 1 | `--all` 자동 페이지네이션 | offset/limit 루프로 전체 취득 |
| 2 | 429 Retry-After 대응 | 서버 지정 대기 시간 파싱, 최대 3회 재시도 |
| 3 | 에러 힌트 메시지 | 인증 에러 → `hint: run 'freee auth login'` 등 |
| 4 | UserAgent에 버전 반영 | `planitaicojp/freee-cli/0.2.0` |
| 5 | `--dry-run` 플래그 | 변경계 커맨드에서 요청 내용 프리뷰 |

### v0.1.1 변경 내용

**문제:** freee API 응답은 래퍼 객체로 감싸진다 (`{"deals": [...]}`).
v0.1.0에서는 `var resp any`로 받았기 때문에 table 출력이 `map[deals:map[...]]`이 되었다.

**수정:**
- `internal/model/`에 typed struct 정의
- `list` 커맨드: 래퍼를 제거하고 테이블 포매터에 전달
- `show` 커맨드: key-value 형식으로 읽기 좋게 표시

---

*이 문서는 구현과 동기화하여 업데이트한다. SPEC과 구현이 어긋나는 경우 Issue를 등록할 것.*
