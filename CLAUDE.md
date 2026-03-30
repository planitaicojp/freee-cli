# planitai-freee-cli

## 개요

freee 공개 API에 대한 agent-friendly CLI 도구. Go 언어로 단일 실행파일 제공.
`aws`, `gh`, `conoha` CLI와 유사한 인터페이스를 지향한다.

**참고 프로젝트**: `/home/tkim/dev/crowdy/conoha-cli` — 구조, 패턴, 빌드 시스템 참조

## 핵심 설계 방침

- **SDK 미사용**: 공식 SDK는 대부분 archived, 커뮤니티 SDK는 업데이트가 느릴 수 있음. OpenAPI 스펙을 참조하여 직접 HTTP 호출
- **OpenAPI 스펙 참조**: `github.com/freee/freee-api-schema` (Accounting, HR, Invoice, Time Tracking, Sales)
- **OAuth2 인증**: Authorization Code + PKCE (S256) 플로우 (브라우저 기반 로그인)
- **OAuth2 주의사항**:
  - `scope=read write` 필수 (없으면 "必須パラメータが不足" 에러)
  - `prompt=select_company` 사용하지 않음 (에러 원인)
  - `AuthStyle: AuthStyleInParams` (freee는 body에 credentials 전송)
  - 콜백 포트 고정 8080 (freee 앱에 `http://localhost:8080/callback` 등록 필요)
  - 앱은 draft 상태로 사용 가능 (active 불필요)
- **멀티 프로필**: `gh auth list` 처럼 여러 계정/사업소 전환 가능
- **Agent-friendly**: `--format json`, `--no-input`, 명확한 exit code, stderr/stdout 분리

## 기술 스택

| 항목 | 선택 |
|------|------|
| 언어 | Go 1.26+ |
| CLI 프레임워크 | `github.com/spf13/cobra` |
| 빌드 | Makefile + GoReleaser |
| HTTP 클라이언트 | `net/http` (SDK 미사용, 직접 호출) |
| OAuth2 | `golang.org/x/oauth2` (표준 라이브러리 수준) |
| 설정 파일 | YAML (`~/.config/freee/`) |
| 출력 형식 | JSON / YAML / CSV / Table |

## freee API 정보

- **Base URL**: `https://api.freee.co.jp/api/1/`
- **인증**: OAuth2 Bearer token
- **Authorization**: `https://accounts.secure.freee.co.jp/public_api/authorize`
- **Token**: `https://accounts.secure.freee.co.jp/public_api/token`
- **OpenAPI Schema**: `https://github.com/freee/freee-api-schema`
- **Developer Portal**: `https://developer.freee.co.jp`
- **App 등록**: `https://app.secure.freee.co.jp/developers`

## 대상 API (5종)

| API | 설명 | 우선순위 |
|-----|------|----------|
| Accounting (会計) | 거래, 경비, 청구서, 분개, 결산 | Phase 1~4 |
| HR (人事労務) | 직원, 근태, 급여, 연말정산 | Phase 5 |
| Invoice (請求書) | 청구서, 견적서, 납품서 | Phase 6 |
| PM (工数管理) | 프로젝트, 팀, 작업시간 | Phase 6 |
| Sales (販売) | 고객, 수주, 매출 | Phase 6 |

## 현재 버전

**v0.4.2** (2026-03-30)

## 디렉토리 구조

```
planitai-freee-cli/
├── main.go
├── cmd/
│   ├── root.go                    # 루트 커맨드, 글로벌 플래그
│   ├── version.go
│   ├── completion.go
│   ├── cmdutil/                   # 공통 유틸 (클라이언트 생성, 인자 검증)
│   ├── auth/                      # login, logout, status, list, switch, remove
│   ├── company/                   # 사업소 조회/전환
│   ├── deal/                      # 거래 CRUD
│   ├── invoice/                   # 청구서 CRUD
│   ├── partner/                   # 거래처 CRUD
│   ├── account/                   # 勘定科目 조회
│   ├── section/                   # 부문 CRUD
│   ├── tag/                       # 메모태그 CRUD
│   ├── item/                      # 품목 CRUD
│   ├── journal/                   # 분개 조회
│   ├── expense/                   # 경비신청 CRUD
│   ├── walletable/                # 구좌 조회
│   ├── employee/                  # 직원 (HR API) — 미구현
│   └── config/                    # 설정 표시/변경
├── internal/
│   ├── api/                       # HTTP 클라이언트, 인증, API 래퍼
│   ├── config/                    # 프로필/인증정보/토큰 관리
│   ├── model/                     # API 응답 구조체 (typed models)
│   ├── output/                    # JSON/YAML/CSV/Table 포매터
│   ├── errors/                    # 에러 타입, exit code
│   └── prompt/                    # 대화형 입력
├── docs/                          # 다국어 문서 (ja/ko/en)
├── Makefile
├── .goreleaser.yaml
├── SPEC.md                        # 설계 문서, 커맨드 일람, 진행 상황
└── CLAUDE.md
```

## GitHub

- **Repo**: https://github.com/planitaicojp/freee-cli (public)
- **브랜치 전략**: feature 브랜치 → PR → main 머지 → 태그 → release

## 설정 파일 경로

```
~/.config/freee/
├── config.yaml        # 프로필, 기본값 (0600)
├── credentials.yaml   # OAuth 토큰 (0600)
└── tokens.yaml        # 토큰 캐시 (0600)
```

## 환경변수

| 변수 | 설명 |
|------|------|
| `FREEE_PROFILE` | 활성 프로필 |
| `FREEE_COMPANY_ID` | 사업소 ID 오버라이드 |
| `FREEE_TOKEN` | 직접 토큰 지정 (CI용) |
| `FREEE_FORMAT` | 출력 형식 |
| `FREEE_CONFIG_DIR` | 설정 디렉토리 |
| `FREEE_NO_INPUT` | 비대화형 모드 |
| `FREEE_DEBUG` | 디버그 로깅 (1=verbose, api=상세) |

## Exit Code

| Code | 의미 |
|------|------|
| 0 | 성공 |
| 1 | 일반 오류 |
| 2 | 인증 오류 |
| 3 | Not Found |
| 4 | 검증 오류 |
| 5 | API 오류 |
| 6 | 네트워크 오류 |
| 10 | 취소됨 |

## 개발 명령어

```bash
make build      # 빌드
make test       # 테스트 (-race 포함)
make lint       # 린트
make coverage   # 커버리지 리포트
make install    # 설치
```

## 출력 패턴 (v0.1.1~)

- **list 커맨드**: `internal/model/` 의 typed struct 사용 → `*Row` 구조체 → table 포맷
- **list --all**: 자동 페이지네이션 (offset/limit 루프로 전체 취득)
- **show 커맨드**: `*Response` wrapper struct 로 JSON 디코딩 → `fmt.Printf` key-value 출력
- **JSON/YAML/CSV**: `var resp any` 로 raw API 응답 그대로 출력
- **--dry-run**: 변경계 커맨드에서 API 호출 없이 리퀘스트 내용 프리뷰
- freee API 응답은 항상 래퍼로 감싸짐 (예: `{"deal": {...}}`, `{"partners": [...]}`)

## API 클라이언트 동작 (v0.2.0~)

- **Retry**: 429/5xx 시 최대 3회 재시도, Retry-After 헤더 파싱
- **에러 힌트**: 모든 에러 타입에 다음 액션 안내 (예: `hint: run 'freee auth login'`)
- **UserAgent**: `planitaicojp/freee-cli/<version>` 형식

## 릴리스 히스토리

| 버전 | 날짜 | 내용 |
|------|------|------|
| v0.1.0 | 2026-03-10 | 초기 스캐폴딩, OAuth2 인증, Phase 1 커맨드 구현 |
| v0.1.1 | 2026-03-10 | list/show 출력 개선 (typed models, key-value format) |
| v0.2.0 | 2026-03-11 | Agent 신뢰성 (--all 페이지네이션, Retry-After, 에러 힌트, --dry-run) |
| v0.2.1 | 2026-03-11 | 유닛테스트 (api/errors/output) + GitHub Actions CI |

## 참고 자료

- freee OpenAPI Schema: https://github.com/freee/freee-api-schema
- freee Developer Portal: https://developer.freee.co.jp
- freee MCP Server (참고): https://github.com/freee/freee-mcp
- conoha-cli (구조 참고): /home/tkim/dev/crowdy/conoha-cli
