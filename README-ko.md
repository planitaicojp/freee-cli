# freee - freee API CLI

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[日本語](README.md) | [English](README-en.md)

freee 공개 API용 커맨드라인 인터페이스입니다. Go로 작성된 단일 바이너리로, 에이전트 친화적 설계를 채택하고 있습니다.

> **주의**: 본 도구는 비공식이며, freee 주식회사와 제휴·추천 관계에 있지 않습니다.

## 특징

- 단일 바이너리, 크로스 플랫폼 지원 (Linux / macOS / Windows)
- OAuth2 Authorization Code + PKCE 브라우저 기반 로그인
- 멀티 프로필 지원 (`gh auth` 스타일)
- 구조화된 출력 (`--format json/yaml/csv/table`)
- 에이전트 친화적 설계 (`--no-input`, 명확한 exit code, stderr/stdout 분리)
- 토큰 자동 갱신 (액세스 토큰 6시간, 리프레시 토큰 90일)
- SDK 미사용 — OpenAPI 스펙 기반 직접 HTTP 호출

## 설치

### 소스에서 빌드

```bash
go install github.com/planitaicojp/freee-cli@latest
```

### Git에서 빌드

```bash
git clone https://github.com/planitaicojp/freee-cli.git
cd freee-cli
make build
sudo mv freee /usr/local/bin/
```

### 릴리스 바이너리

[Releases](https://github.com/planitaicojp/freee-cli/releases) 페이지에서 다운로드하거나 아래 명령어를 사용하세요:

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

## 사전 준비

1. [freee 개발자 콘솔](https://app.secure.freee.co.jp/developers/applications/new)에서 앱 생성
   - **앱 타입**: `プライベート` (프라이빗)
   - **콜백 URL**: `http://localhost:8080/callback`
2. **Client ID** 와 **Client Secret** 메모

> 앱은 `draft` 상태 그대로 사용할 수 있습니다.

자세한 내용은 [빠른 시작 가이드](docs/getting-started-ko.md)를 참조하세요.

## 빠른 시작

```bash
# 로그인 (브라우저가 열립니다)
freee auth login

# 인증 상태 확인
freee auth status

# 사업소 목록
freee company list

# 거래 목록
freee deal list

# JSON 형식 출력
freee deal list --format json

# 거래처 목록
freee partner list
```

## 명령어 목록

| 명령어 | 설명 |
|--------|------|
| `freee auth` | 인증 관리 (login / logout / status / list / switch / token / remove) |
| `freee company` | 사업소 관리 (list / show / switch) |
| `freee deal` | 거래 관리 (list / show / create / update / delete) |
| `freee invoice` | 청구서 관리 (list / show / create / update / delete) |
| `freee partner` | 거래처 관리 (list / show / create / update / delete) |
| `freee account` | 계정과목 (list / show) |
| `freee section` | 부문 관리 (list / create / update / delete) |
| `freee tag` | 메모태그 관리 (list / create / update / delete) |
| `freee item` | 품목 관리 (list / create / update / delete) |
| `freee journal` | 분개장 (list) |
| `freee expense` | 경비신청 관리 (list / show / create / update / delete) |
| `freee walletable` | 구좌 관리 (list / show) |
| `freee config` | CLI 설정 관리 (show / set / path) |

## 설정

설정 파일은 `~/.config/freee/`에 저장됩니다:

| 파일 | 설명 | 퍼미션 |
|------|------|--------|
| `config.yaml` | 프로필 설정 | 0600 |
| `credentials.yaml` | OAuth 토큰 | 0600 |

### 환경변수

| 변수 | 설명 |
|------|------|
| `FREEE_PROFILE` | 사용할 프로필명 |
| `FREEE_COMPANY_ID` | 사업소 ID |
| `FREEE_TOKEN` | 액세스 토큰 (직접 지정, CI용) |
| `FREEE_FORMAT` | 출력 형식 |
| `FREEE_CONFIG_DIR` | 설정 디렉토리 |
| `FREEE_NO_INPUT` | 비대화형 모드 (`1` 또는 `true`) |
| `FREEE_DEBUG` | 디버그 로그 (`1` 또는 `api`) |

우선순위: 환경변수 > 플래그 > 프로필 설정 > 기본값

### 글로벌 플래그

```
--profile      사용할 프로필
--format       출력 형식 (table / json / yaml / csv)
--company-id   사업소 ID (오버라이드)
--no-input     대화형 프롬프트 비활성화
--quiet        불필요한 출력 억제
--verbose      상세 출력
--no-color     컬러 출력 비활성화
```

## Exit Code

| 코드 | 의미 |
|------|------|
| 0 | 성공 |
| 1 | 일반 오류 |
| 2 | 인증 실패 |
| 3 | 리소스 미발견 |
| 4 | 검증 오류 |
| 5 | API 오류 |
| 6 | 네트워크 오류 |
| 10 | 사용자 취소 |

## 에이전트 연동

본 CLI는 스크립트 및 AI 에이전트에서의 활용을 염두에 두고 설계되었습니다:

```bash
# 비대화형 모드로 JSON 출력
freee deal list --format json --no-input

# 토큰을 얻어 스크립트에서 활용
TOKEN=$(freee auth token)

# exit code로 에러 핸들링
freee deal show 12345 || echo "Exit code: $?"

# CI/CD 환경에서 사용
export FREEE_TOKEN="<access-token>"
export FREEE_COMPANY_ID="12345678"
freee deal list --format json
```

## 개발

```bash
make build     # 바이너리 빌드
make test      # 테스트 실행
make lint      # 린터 실행
make clean     # 빌드 결과물 삭제
```

## 대응 API

| API | 상태 |
|-----|------|
| [회계 API](https://developer.freee.co.jp/reference/accounting/reference) | Phase 1 대응 중 |
| [인사노무 API](https://developer.freee.co.jp/reference/hr) | Phase 2 예정 |
| [청구서 API](https://developer.freee.co.jp/reference/iv) | Phase 2 예정 |
| [공수관리 API](https://developer.freee.co.jp/reference/time-tracking) | Phase 3 예정 |
| [판매 API](https://developer.freee.co.jp/reference/sales) | Phase 3 예정 |

## 참고

- [freee API 레퍼런스](https://developer.freee.co.jp)
- [freee OpenAPI 스키마](https://github.com/freee/freee-api-schema)

## 라이선스

MIT License - 자세한 내용은 [LICENSE](LICENSE)를 참조하세요.
