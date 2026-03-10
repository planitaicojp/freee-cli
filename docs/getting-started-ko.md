# freee CLI 빠른 시작 가이드

## 사전 조건

- freee 계정 ([freee.co.jp](https://www.freee.co.jp) 에서 등록)
- freee 개발자 앱 (아래에서 생성)

## 1. freee 개발자 앱 생성

1. [freee 개발자 콘솔](https://app.secure.freee.co.jp/developers/applications/new) 에 접속
2. 다음 설정으로 앱 생성:
   - **앱 타입 (アプリタイプ)**: `プライベート` (본인만 사용하는 경우)
   - **앱 이름**: 임의 (예: `freee-cli`)
   - **콜백 URL (コールバックURL)**: `http://localhost:8080/callback`
3. 생성 후, **Client ID** 와 **Client Secret** 을 메모

> **참고**: 앱은 `draft` 상태 그대로 사용할 수 있습니다. `active` (공개) 로 변경할 필요가 없습니다.

## 2. 로그인

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

브라우저가 자동으로 열리고, freee 로그인/인가 화면이 표시됩니다.
인가를 허용하면, CLI에 자동으로 토큰이 저장됩니다.

### 브라우저가 열리지 않는 경우

터미널에 표시되는 URL을 수동으로 브라우저에 붙여넣으세요.

## 3. 인증 상태 확인

```bash
$ freee auth status
Profile:   default
Email:     taro@example.com
Company:   株式会社サンプル (ID: 12345678)
Token:     valid (expires in 5h45m, 2026-03-11 02:53 JST)
```

## 4. 기본 사용법

```bash
# 사업소 목록
$ freee company list

# 거래 목록
$ freee deal list

# JSON 출력 (에이전트/스크립트용)
$ freee deal list --format json

# 거래처 목록
$ freee partner list

# 계정과목 목록
$ freee account list
```

## 5. 다중 계정 관리

```bash
# 다른 프로필로 로그인
$ freee auth login --profile sub-account

# 프로필 목록
$ freee auth list

# 프로필 전환
$ freee auth switch sub-account
```

## 6. 사업소 전환

여러 사업소에 접근할 수 있는 경우:

```bash
# 사업소 목록 확인
$ freee company list

# 기본 사업소 변경
$ freee company switch <company-id>
```

## 7. CI/CD 활용

```bash
# 환경변수로 토큰 직접 지정
export FREEE_TOKEN="<access-token>"
export FREEE_COMPANY_ID="12345678"

# 비대화형 모드로 실행
freee deal list --no-input --format json
```

## 문제 해결

### "必須パラメータが不足" (필수 파라미터 부족) 에러

- 앱의 **콜백 URL**이 `http://localhost:8080/callback` 으로 설정되어 있는지 확인
- 앱의 **앱 타입**이 `プライベート` 인지 확인

### 포트 8080이 사용 중

다른 애플리케이션이 포트 8080을 사용하고 있다면, 해당 애플리케이션을 종료한 후 다시 로그인하세요.

### 토큰 만료

액세스 토큰은 6시간 유효합니다. 만료 시 다음 API 호출 시 리프레시 토큰으로 자동 갱신됩니다. 리프레시 토큰의 유효기간은 90일입니다.

## 설정 파일

설정은 `~/.config/freee/` 에 저장됩니다:

```
~/.config/freee/
├── config.yaml        # 프로필 설정
└── credentials.yaml   # OAuth 토큰 (0600 퍼미션)
```
