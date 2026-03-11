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

### Phase 2: Accounting 확장 — 仕訳・振替・明細・税 (19 commands)

| # | 명령어 | 설명 | 진행 |
|---|--------|------|------|
| **仕訳帳 (Manual Journals)** | | | |
| 53 | `freee manual-journal list` | 仕訳 목록 | ☐ |
| 54 | `freee manual-journal show <id>` | 仕訳 상세 | ☐ |
| 55 | `freee manual-journal create` | 仕訳 등록 | ☐ |
| 56 | `freee manual-journal update <id>` | 仕訳 수정 | ☐ |
| 57 | `freee manual-journal delete <id>` | 仕訳 삭제 | ☐ |
| **振替伝票 (Transfers)** | | | |
| 58 | `freee transfer list` | 振替 목록 | ☐ |
| 59 | `freee transfer show <id>` | 振替 상세 | ☐ |
| 60 | `freee transfer create` | 振替 등록 | ☐ |
| 61 | `freee transfer update <id>` | 振替 수정 | ☐ |
| 62 | `freee transfer delete <id>` | 振替 삭제 | ☐ |
| **口座明細 (Wallet Transactions)** | | | |
| 63 | `freee wallet-txn list` | 口座明細 목록 | ☐ |
| 64 | `freee wallet-txn show <id>` | 口座明細 상세 | ☐ |
| 65 | `freee wallet-txn create` | 口座明細 등록 | ☐ |
| 66 | `freee wallet-txn delete <id>` | 口座明細 삭제 | ☐ |
| **口座 (Walletables) 拡張** | | | |
| 67 | `freee walletable create` | 口座 등록 | ☐ |
| 68 | `freee walletable update <id>` | 口座 수정 | ☐ |
| 69 | `freee walletable delete <id>` | 口座 삭제 | ☐ |
| **税区分 (Taxes)** | | | |
| 70 | `freee tax list` | 税区分 목록 | ☐ |
| 71 | `freee tax show <code>` | 税区分 상세 | ☐ |

### Phase 3: Accounting 확장 — レポート + ファイルボックス (25 commands)

| # | 명령어 | 설명 | 진행 |
|---|--------|------|------|
| **レポート (Reports)** | | | |
| 72 | `freee report bs` | 貸借対照表 | ☐ |
| 73 | `freee report bs-2y` | 貸借対照表 2期比較 | ☐ |
| 74 | `freee report bs-3y` | 貸借対照表 3期比較 | ☐ |
| 75 | `freee report bs-sections` | 貸借対照表 部門別 | ☐ |
| 76 | `freee report bs-segments` | 貸借対照表 セグメント別 | ☐ |
| 77 | `freee report pl` | 損益計算書 | ☐ |
| 78 | `freee report pl-2y` | 損益計算書 2期比較 | ☐ |
| 79 | `freee report pl-3y` | 損益計算書 3期比較 | ☐ |
| 80 | `freee report pl-sections` | 損益計算書 部門別 | ☐ |
| 81 | `freee report pl-segments` | 損益計算書 セグメント別 | ☐ |
| 82 | `freee report cr` | キャッシュフロー計算書 | ☐ |
| 83 | `freee report cr-2y` | キャッシュフロー 2期比較 | ☐ |
| 84 | `freee report cr-3y` | キャッシュフロー 3期比較 | ☐ |
| 85 | `freee report cr-sections` | キャッシュフロー 部門別 | ☐ |
| 86 | `freee report cr-segments` | キャッシュフロー セグメント別 | ☐ |
| 87 | `freee report general-ledger` | 総勘定元帳 | ☐ |
| 88 | `freee report trial-bs` | 残高試算表 (BS) | ☐ |
| 89 | `freee report trial-pl` | 残高試算表 (PL) | ☐ |
| **ファイルボックス (Receipts)** | | | |
| 90 | `freee receipt list` | 証憑 목록 | ☐ |
| 91 | `freee receipt show <id>` | 証憑 상세 | ☐ |
| 92 | `freee receipt create` | 証憑 アップロード | ☐ |
| 93 | `freee receipt update <id>` | 証憑 수정 | ☐ |
| 94 | `freee receipt delete <id>` | 証憑 삭제 | ☐ |
| 95 | `freee receipt download <id>` | 証憑 ダウンロード | ☐ |
| **銀行 (Banks)** | | | |
| 96 | `freee bank list` | 銀行 목록 | ☐ |

### Phase 4: Accounting 확장 — ワークフロー + サブリソース (44 commands)

| # | 명령어 | 설명 | 진행 |
|---|--------|------|------|
| **承認 (Approvals)** | | | |
| 97 | `freee approval list` | 承認依頼 목록 | ☐ |
| 98 | `freee approval show <id>` | 承認依頼 상세 | ☐ |
| 99 | `freee approval create` | 承認依頼 작성 | ☐ |
| 100 | `freee approval update <id>` | 承認依頼 수정 | ☐ |
| 101 | `freee approval delete <id>` | 承認依頼 삭제 | ☐ |
| 102 | `freee approval action <id>` | 承認/却下 アクション | ☐ |
| **承認フォーム/経路** | | | |
| 103 | `freee approval-form list` | 承認フォーム 목록 | ☐ |
| 104 | `freee approval-form show <id>` | 承認フォーム 상세 | ☐ |
| 105 | `freee approval-route list` | 承認経路 목록 | ☐ |
| 106 | `freee approval-route show <id>` | 承認経路 상세 | ☐ |
| **支払依頼 (Payment Requests)** | | | |
| 107 | `freee payment-request list` | 支払依頼 목록 | ☐ |
| 108 | `freee payment-request show <id>` | 支払依頼 상세 | ☐ |
| 109 | `freee payment-request create` | 支払依頼 작성 | ☐ |
| 110 | `freee payment-request update <id>` | 支払依頼 수정 | ☐ |
| 111 | `freee payment-request delete <id>` | 支払依頼 삭제 | ☐ |
| 112 | `freee payment-request action <id>` | 支払依頼 承認/却下 | ☐ |
| **経費テンプレート (Expense Templates)** | | | |
| 113 | `freee expense-template list` | 経費テンプレート 목록 | ☐ |
| 114 | `freee expense-template show <id>` | 経費テンプレート 상세 | ☐ |
| 115 | `freee expense-template create` | 경비 テンプレート 등록 | ☐ |
| 116 | `freee expense-template update <id>` | 경비 テンプレート 수정 | ☐ |
| 117 | `freee expense-template delete <id>` | 경비 テンプレート 삭제 | ☐ |
| **経費 アクション** | | | |
| 118 | `freee expense action <id>` | 경비신청 承認/却下 | ☐ |
| **取引 決済 (Deal Payments)** | | | |
| 119 | `freee deal payment create <deal-id>` | 取引 決済 등록 | ☐ |
| 120 | `freee deal payment update <deal-id> <id>` | 取引 決済 수정 | ☐ |
| 121 | `freee deal payment delete <deal-id> <id>` | 取引 決済 삭제 | ☐ |
| **取引 更新 (Deal Renews)** | | | |
| 122 | `freee deal renew create <deal-id>` | 取引 +更新行 등록 | ☐ |
| 123 | `freee deal renew update <deal-id> <id>` | 取引 更新行 수정 | ☐ |
| 124 | `freee deal renew delete <deal-id> <id>` | 取引 更新行 삭제 | ☐ |
| **見積書 (Quotations)** | | | |
| 125 | `freee quotation list` | 見積書 목록 | ☐ |
| 126 | `freee quotation show <id>` | 見積書 상세 | ☐ |
| **セグメントタグ (Segment Tags)** | | | |
| 127 | `freee segment-tag list` | セグメントタグ 목록 | ☐ |
| 128 | `freee segment-tag create` | セグメントタグ 등록 | ☐ |
| 129 | `freee segment-tag update <id>` | セグメントタグ 수정 | ☐ |
| 130 | `freee segment-tag delete <id>` | セグメントタグ 삭제 | ☐ |
| **事業所コード (Company Codes)** | | | |
| 131 | `freee code account upsert` | 勘定科目コード 登録/更新 | ☐ |
| 132 | `freee code section upsert` | 部門コード 登録/更新 | ☐ |
| 133 | `freee code item upsert` | 品目コード 登録/更新 | ☐ |
| 134 | `freee code tag upsert` | タグコード 登録/更新 | ☐ |
| 135 | `freee code walletable upsert` | 口座コード 登録/更新 | ☐ |
| **取引先コード** | | | |
| 136 | `freee partner code create <partner-id>` | 取引先コード 등록 | ☐ |
| 137 | `freee partner code update <partner-id>` | 取引先コード 수정 | ☐ |
| **ユーザー** | | | |
| 138 | `freee user list` | ユーザー 목록 | ☐ |
| 139 | `freee user me` | 現在ユーザー情報 | ☐ |
| 140 | `freee user capabilities` | ユーザー権限情報 | ☐ |

### Phase 5: HR API 전체 (93 commands)

| # | 명령어 | 설명 | 진행 |
|---|--------|------|------|
| **従業員 (Employees)** | | | |
| 141 | `freee hr-employee list` | 従業員 목록 | ☐ |
| 142 | `freee hr-employee show <id>` | 従業員 상세 | ☐ |
| 143 | `freee hr-employee create` | 従業員 등록 | ☐ |
| 144 | `freee hr-employee update <id>` | 従業員 수정 | ☐ |
| 145 | `freee hr-employee delete <id>` | 従業員 삭제 | ☐ |
| **勤怠 (Work Records)** | | | |
| 146 | `freee hr-work-record list` | 勤怠 목록 | ☐ |
| 147 | `freee hr-work-record show <id>` | 勤怠 상세 | ☐ |
| 148 | `freee hr-work-record create` | 勤怠 등록 | ☐ |
| 149 | `freee hr-work-record update <id>` | 勤怠 수정 | ☐ |
| 150 | `freee hr-work-record delete <id>` | 勤怠 삭제 | ☐ |
| **打刻 (Time Clocks)** | | | |
| 151 | `freee hr-time-clock list` | 打刻 목록 | ☐ |
| 152 | `freee hr-time-clock show <id>` | 打刻 상세 | ☐ |
| 153 | `freee hr-time-clock create` | 打刻 등록 | ☐ |
| 154 | `freee hr-time-clock available-types` | 打刻可能種別 | ☐ |
| **所定勤務タグ (Attendance Tags)** | | | |
| 155 | `freee hr-attendance-tag list` | 所定勤務タグ 목록 | ☐ |
| 156 | `freee hr-attendance-tag show <id>` | 所定勤務タグ 상세 | ☐ |
| 157 | `freee hr-attendance-tag create` | 所定勤務タグ 등록 | ☐ |
| 158 | `freee hr-attendance-tag update <id>` | 所定勤務タグ 수정 | ☐ |
| 159 | `freee hr-attendance-tag delete <id>` | 所定勤務タグ 삭제 | ☐ |
| **給与/賞与 (Salary & Bonus)** | | | |
| 160 | `freee hr-salary list` | 給与明細 목록 | ☐ |
| 161 | `freee hr-salary show <id>` | 給与明細 상세 | ☐ |
| 162 | `freee hr-bonus list` | 賞与明細 목록 | ☐ |
| 163 | `freee hr-bonus show <id>` | 賞与明細 상세 | ☐ |
| **グループ/役職** | | | |
| 164 | `freee hr-group list` | グループ 목록 | ☐ |
| 165 | `freee hr-group show <id>` | グループ 상세 | ☐ |
| 166 | `freee hr-group create` | グループ 등록 | ☐ |
| 167 | `freee hr-group update <id>` | グループ 수정 | ☐ |
| 168 | `freee hr-position list` | 役職 목록 | ☐ |
| 169 | `freee hr-position show <id>` | 役職 상세 | ☐ |
| 170 | `freee hr-position create` | 役職 등록 | ☐ |
| **従業員ルール (Employee Sub-resources)** | | | |
| 171 | `freee hr-employee bank-account show <emp-id>` | 振込先口座 | ☐ |
| 172 | `freee hr-employee bank-account update <emp-id>` | 振込先口座 수정 | ☐ |
| 173 | `freee hr-employee basic-pay show <emp-id>` | 基本給 | ☐ |
| 174 | `freee hr-employee basic-pay update <emp-id>` | 基本給 수정 | ☐ |
| 175 | `freee hr-employee dependents show <emp-id>` | 扶養親族 | ☐ |
| 176 | `freee hr-employee dependents update <emp-id>` | 扶養親族 수정 | ☐ |
| 177 | `freee hr-employee health-insurance show <emp-id>` | 健康保険 | ☐ |
| 178 | `freee hr-employee health-insurance update <emp-id>` | 健康保険 수정 | ☐ |
| 179 | `freee hr-employee welfare-pension show <emp-id>` | 厚生年金 | ☐ |
| 180 | `freee hr-employee welfare-pension update <emp-id>` | 厚生年金 수정 | ☐ |
| 181 | `freee hr-employee profile show <emp-id>` | プロフィール | ☐ |
| 182 | `freee hr-employee profile update <emp-id>` | プロフィール 수정 | ☐ |
| 183 | `freee hr-employee work-record-summary show <emp-id>` | 勤怠サマリー | ☐ |
| 184 | `freee hr-employee tax-withholding show <emp-id>` | 源泉徴収 | ☐ |
| 185 | `freee hr-employee social-insurance show <emp-id>` | 社会保険 | ☐ |
| 186 | `freee hr-employee social-insurance update <emp-id>` | 社会保険 수정 | ☐ |
| 187 | `freee hr-employee employment-insurance show <emp-id>` | 雇用保険 | ☐ |
| **HR 承認ワークフロー** | | | |
| 188 | `freee hr-approval-monthly list` | 月次勤怠承認 목록 | ☐ |
| 189 | `freee hr-approval-monthly show <id>` | 月次勤怠承認 상세 | ☐ |
| 190 | `freee hr-approval-monthly create` | 月次勤怠承認 작성 | ☐ |
| 191 | `freee hr-approval-monthly update <id>` | 月次勤怠承認 수정 | ☐ |
| 192 | `freee hr-approval-monthly delete <id>` | 月次勤怠承認 삭제 | ☐ |
| 193 | `freee hr-approval-monthly action <id>` | 月次勤怠 承認/却下 | ☐ |
| 194 | `freee hr-approval-overtime list` | 残業承認 목록 | ☐ |
| 195 | `freee hr-approval-overtime show <id>` | 残業承認 상세 | ☐ |
| 196 | `freee hr-approval-overtime create` | 残業承認 작성 | ☐ |
| 197 | `freee hr-approval-overtime update <id>` | 残業承認 수정 | ☐ |
| 198 | `freee hr-approval-overtime delete <id>` | 残業承認 삭제 | ☐ |
| 199 | `freee hr-approval-overtime action <id>` | 残業 承認/却下 | ☐ |
| 200 | `freee hr-approval-paid-leave list` | 有休承認 목록 | ☐ |
| 201 | `freee hr-approval-paid-leave show <id>` | 有休承認 상세 | ☐ |
| 202 | `freee hr-approval-paid-leave create` | 有休承認 작성 | ☐ |
| 203 | `freee hr-approval-paid-leave update <id>` | 有休承認 수정 | ☐ |
| 204 | `freee hr-approval-paid-leave delete <id>` | 有休承認 삭제 | ☐ |
| 205 | `freee hr-approval-paid-leave action <id>` | 有休 承認/却下 | ☐ |
| 206 | `freee hr-approval-special-leave list` | 特別休暇承認 목록 | ☐ |
| 207 | `freee hr-approval-special-leave show <id>` | 特別休暇承認 상세 | ☐ |
| 208 | `freee hr-approval-special-leave create` | 特別休暇承認 작성 | ☐ |
| 209 | `freee hr-approval-special-leave update <id>` | 特別休暇承認 수정 | ☐ |
| 210 | `freee hr-approval-special-leave delete <id>` | 特別休暇承認 삭제 | ☐ |
| 211 | `freee hr-approval-special-leave action <id>` | 特別休暇 承認/却下 | ☐ |
| 212 | `freee hr-approval-work-time list` | 勤務時間承認 목록 | ☐ |
| 213 | `freee hr-approval-work-time show <id>` | 勤務時間承認 상세 | ☐ |
| 214 | `freee hr-approval-work-time create` | 勤務時間承認 작성 | ☐ |
| 215 | `freee hr-approval-work-time update <id>` | 勤務時間承認 수정 | ☐ |
| 216 | `freee hr-approval-work-time delete <id>` | 勤務時間承認 삭제 | ☐ |
| 217 | `freee hr-approval-work-time action <id>` | 勤務時間 承認/却下 | ☐ |
| **HR 承認経路** | | | |
| 218 | `freee hr-approval-route list` | HR承認経路 목록 | ☐ |
| 219 | `freee hr-approval-route show <id>` | HR承認経路 상세 | ☐ |
| **年末調整 (Year-End Adjustment)** | | | |
| 220 | `freee hr-yearend employees` | 年末調整対象者 | ☐ |
| 221 | `freee hr-yearend dependents <emp-id>` | 扶養控除 상세 | ☐ |
| 222 | `freee hr-yearend housing-loans <emp-id>` | 住宅ローン控除 | ☐ |
| 223 | `freee hr-yearend insurances <emp-id>` | 保険料控除 | ☐ |
| 224 | `freee hr-yearend life-insurances <emp-id>` | 生命保険料控除 | ☐ |
| 225 | `freee hr-yearend social-insurances <emp-id>` | 社会保険料控除 | ☐ |
| 226 | `freee hr-yearend earthquake-insurances <emp-id>` | 地震保険料控除 | ☐ |
| 227 | `freee hr-yearend payroll <emp-id>` | 給与所得 | ☐ |
| 228 | `freee hr-yearend previous-jobs <emp-id>` | 前職情報 | ☐ |
| 229 | `freee hr-yearend base <emp-id>` | 基本情報 | ☐ |
| 230 | `freee hr-yearend status <emp-id>` | ステータス | ☐ |
| 231 | `freee hr-yearend result <emp-id>` | 計算結果 | ☐ |
| 232 | `freee hr-yearend summary` | 年末調整サマリー | ☐ |
| 233 | `freee hr-user me` | HR現在ユーザー情報 | ☐ |

### Phase 6: Invoice API + PM API + Sales API (39 commands)

| # | 명령어 | 설명 | 진행 |
|---|--------|------|------|
| **Invoice API (/iv)** | | | |
| 234 | `freee iv-invoice list` | 請求書 목록 (Invoice API) | ☐ |
| 235 | `freee iv-invoice show <id>` | 請求書 상세 | ☐ |
| 236 | `freee iv-invoice create` | 請求書 작성 | ☐ |
| 237 | `freee iv-invoice update <id>` | 請求書 수정 | ☐ |
| 238 | `freee iv-invoice delete <id>` | 請求書 삭제 | ☐ |
| 239 | `freee iv-quotation list` | 見積書 목록 (Invoice API) | ☐ |
| 240 | `freee iv-quotation show <id>` | 見積書 상세 | ☐ |
| 241 | `freee iv-quotation create` | 見積書 작성 | ☐ |
| 242 | `freee iv-quotation update <id>` | 見積書 수정 | ☐ |
| 243 | `freee iv-quotation delete <id>` | 見積書 삭제 | ☐ |
| 244 | `freee iv-delivery list` | 納品書 목록 | ☐ |
| 245 | `freee iv-delivery show <id>` | 납品서 상세 | ☐ |
| 246 | `freee iv-delivery create` | 납품서 작성 | ☐ |
| 247 | `freee iv-delivery update <id>` | 납품서 수정 | ☐ |
| 248 | `freee iv-template list` | テンプレート 목록 | ☐ |
| **PM API (/pm — 工数管理)** | | | |
| 249 | `freee pm-project list` | プロジェクト 목록 | ☐ |
| 250 | `freee pm-project show <id>` | プロジェクト 상세 | ☐ |
| 251 | `freee pm-project create` | プロジェクト 등록 | ☐ |
| 252 | `freee pm-project update <id>` | プロジェクト 수정 | ☐ |
| 253 | `freee pm-workload list` | 工数 목록 | ☐ |
| 254 | `freee pm-workload create` | 工数 등록 | ☐ |
| 255 | `freee pm-workload update <id>` | 工数 수정 | ☐ |
| 256 | `freee pm-workload delete <id>` | 工수 삭제 | ☐ |
| 257 | `freee pm-team list` | チーム 목록 | ☐ |
| 258 | `freee pm-team show <id>` | チーム 상세 | ☐ |
| **Sales API (/sm — 販売管理)** | | | |
| 259 | `freee sales-business list` | 案件 목록 | ☐ |
| 260 | `freee sales-business show <id>` | 案件 상세 | ☐ |
| 261 | `freee sales-business create` | 案件 등록 | ☐ |
| 262 | `freee sales-business update <id>` | 案件 수정 | ☐ |
| 263 | `freee sales-order list` | 受注 목록 | ☐ |
| 264 | `freee sales-order show <id>` | 受注 상세 | ☐ |
| 265 | `freee sales-order create` | 受注 등록 | ☐ |
| 266 | `freee sales-order update <id>` | 受注 수정 | ☐ |
| 267 | `freee sales-customer list` | 顧客 목록 | ☐ |
| 268 | `freee sales-customer show <id>` | 顧客 상세 | ☐ |
| 269 | `freee sales-product list` | 商品 목록 | ☐ |
| 270 | `freee sales-product show <id>` | 商品 상세 | ☐ |
| **メタ情報** | | | |
| 271 | `freee api-resources` | 対応APIリソース一覧 | ☐ |
| 272 | `freee forms selectables` | フォーム選択肢取得 | ☐ |

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
| Phase 1 | 기반 + Accounting 핵심 | 52 | 39 | 75% |
| Phase 2 | Accounting 仕訳・振替・明細・税 | 19 | 0 | 0% |
| Phase 3 | Accounting レポート + ファイルボックス | 25 | 0 | 0% |
| Phase 4 | Accounting ワークフロー + サブリソース | 44 | 0 | 0% |
| Phase 5 | HR API 전체 | 93 | 0 | 0% |
| Phase 6 | Invoice API + PM API + Sales API | 39 | 0 | 0% |
| **합계** | | **272** | **39** | **14%** |

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
| 0.8.1 | TBD | メタ情報 + ドキュメント + GoReleaser 完了 |
| 0.8.0 | TBD | Phase 6 — Invoice API + PM API + Sales API (39 commands) |
| 0.7.3 | TBD | HR 年末調整 (13 commands) |
| 0.7.2 | TBD | HR 承認ワークフロー (33 commands) |
| 0.7.1 | TBD | HR 従業員ルール — サブリソース (17 commands) |
| 0.7.0 | TBD | Phase 5 — HR 基本機能 (30 commands) |
| 0.6.1 | TBD | Accounting マスター・コード操作 (14 commands) |
| 0.6.0 | TBD | Phase 4 — Accounting ワークフロー + サブリソース (30 commands) |
| 0.5.0 | TBD | Phase 3 — Accounting レポート + ファイルボックス (25 commands) |
| 0.4.0 | TBD | Phase 2 — Accounting 仕訳・振替・明細 (19 commands) |
| 0.3.0 | TBD | クライアント設定 + 出力改善 |
| 0.2.2 | TBD | Phase 1 CRUD 完成（残り 13 コマンド） |
| 0.2.1 | TBD | テスト + CI |
| 0.2.0 | TBD | Agent 信頼性向上 |
| 0.1.1 | 2026-03-10 | 出力改善、型付きモデル導入 |
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

---

## Roadmap

### v0.2.0 — Agent 신뢰성

Agent (AI) が安全かつ確実に CLI を使えるようにするための改善。

| # | 項目 | 説明 |
|---|------|------|
| 1 | list 自動ページネーション | `--all` フラグ追加、offset/limit ループで全件取得 |
| 2 | 429 Retry-After 対応 | naive sleep → サーバ指定の待機時間をパース |
| 3 | エラーヒントメッセージ | 認証エラー → "run `freee auth login`" 等、次のアクション提示 |
| 4 | UserAgent にバージョン反映 | `freee-cli/0.2.0` 形式で実バージョンを送信 |
| 5 | `--dry-run` フラグ | 変更系コマンドでリクエスト内容をプレビュー（実行しない） |

### v0.2.1 — テスト + CI

品質保証の基盤整備。

| # | 項目 | 説明 |
|---|------|------|
| 1 | `internal/api/client.go` ユニットテスト | `httptest.Server` によるモックテスト |
| 2 | `internal/errors/` ユニットテスト | エラー型・exit code の検証 |
| 3 | `internal/output/` ユニットテスト | JSON/YAML/CSV/Table 出力の検証 |
| 4 | GitHub Actions CI | `.github/workflows/ci.yml` — build, test, lint |
| 5 | Makefile coverage ターゲット | `make coverage` で `go test -cover` 実行 |

### v0.2.2 — Phase 1 CRUD 完成

残り 13 個のスタブコマンドを実装し、Phase 1 を 100% 完了にする。

| コマンド | 説明 |
|----------|------|
| `deal update` | 取引更新 |
| `invoice create` | 請求書作成 |
| `invoice update` | 請求書更新 |
| `partner create` | 取引先登録 |
| `partner update` | 取引先更新 |
| `expense create` | 経費申請作成 |
| `expense update` | 経費申請更新 |
| `section create` | 部門登録 |
| `section update` | 部門更新 |
| `tag create` | タグ登録 |
| `tag update` | タグ更新 |
| `item create` | 品目登録 |
| `item update` | 品目更新 |

### v0.3.0 — クライアント設定 + 出力改善

CLI の使い勝手を向上させる設定項目と出力フォーマット改善。

| # | 項目 | 説明 |
|---|------|------|
| 1 | baseURL/timeout/maxRetries 設定 | 環境変数 + config.yaml で設定可能に |
| 2 | 金額カンマフォーマット | table モードで `1,234,567` 形式表示 |
| 3 | 日付フォーマット | table モードで読みやすい日付表示 |
| 4 | ステータス日本語ラベル | `settled` → `決済済` 等の日本語変換 |
| 5 | `--no-header` フラグ | テーブルヘッダー非表示（スクリプト向け） |
| 6 | `--output-fields` フラグ | 出力カラムの選択指定 |

### v0.4.0 — Accounting 仕訳・振替・明細 (Phase 2, 19 commands)

Accounting API の仕訳帳・振替伝票・口座明細・税区分を網羅。

| # | 項目 | 説明 |
|---|------|------|
| 1 | `manual-journal` CRUD (5) | 仕訳帳の一覧・詳細・登録・更新・削除 |
| 2 | `transfer` CRUD (5) | 振替伝票の一覧・詳細・登録・更新・削除 |
| 3 | `wallet-txn` CRUD (4) | 口座明細の一覧・詳細・登録・削除 |
| 4 | `walletable` 拡張 (3) | 口座の登録・更新・削除 |
| 5 | `tax` list/show (2) | 税区分の一覧・詳細 |

### v0.5.0 — Accounting レポート + ファイルボックス (Phase 3, 25 commands)

財務レポート（BS/PL/CR × 5パターン + 総勘定元帳 + 試算表）と証憑管理。

| # | 項目 | 説明 |
|---|------|------|
| 1 | `report bs` / `bs-2y` / `bs-3y` / `bs-sections` / `bs-segments` | 貸借対照表 5パターン |
| 2 | `report pl` / `pl-2y` / `pl-3y` / `pl-sections` / `pl-segments` | 損益計算書 5パターン |
| 3 | `report cr` / `cr-2y` / `cr-3y` / `cr-sections` / `cr-segments` | キャッシュフロー 5パターン |
| 4 | `report general-ledger` / `trial-bs` / `trial-pl` | 総勘定元帳 + 残高試算表 |
| 5 | `receipt` CRUD + download (6) | 証憑の一覧・詳細・アップロード・更新・削除・ダウンロード |
| 6 | `bank list` (1) | 銀行一覧 |

### v0.6.0 — Accounting ワークフロー + サブリソース (Phase 4a, 30 commands)

承認ワークフロー、支払依頼、経費テンプレート、取引サブリソース。

| # | 項目 | 説明 |
|---|------|------|
| 1 | `approval` CRUD + action (6) | 承認依頼の一覧〜削除 + 承認/却下 |
| 2 | `approval-form` / `approval-route` (4) | 承認フォーム・経路の一覧・詳細 |
| 3 | `payment-request` CRUD + action (6) | 支払依頼の一覧〜削除 + 承認/却下 |
| 4 | `expense-template` CRUD (5) | 経費テンプレートの一覧〜削除 |
| 5 | `expense action` (1) | 経費申請の承認/却下 |
| 6 | `deal payment` create/update/delete (3) | 取引決済サブリソース |
| 7 | `deal renew` create/update/delete (3) | 取引更新行サブリソース |
| 8 | `quotation` list/show (2) | 見積書の一覧・詳細 |

### v0.6.1 — Accounting マスター・コード操作 (Phase 4b, 14 commands)

セグメントタグ、各種コード操作、取引先コード、ユーザー情報。

| # | 項目 | 説明 |
|---|------|------|
| 1 | `segment-tag` CRUD (4) | セグメントタグの一覧・登録・更新・削除 |
| 2 | `code` upsert (5) | 勘定科目/部門/品目/タグ/口座コードの登録/更新 |
| 3 | `partner code` create/update (2) | 取引先コードの登録・更新 |
| 4 | `user` list/me/capabilities (3) | ユーザー情報・権限 |

### v0.7.0 — HR 基本機能 (Phase 5a, 30 commands)

HR API の基本リソース。Base URL: `/hr/api/v1/`

| # | 項目 | 説明 |
|---|------|------|
| 1 | HR API クライアント | Base URL 切替、`hr-` prefix 命名規則 |
| 2 | `hr-employee` CRUD (5) | 従業員の一覧〜削除 |
| 3 | `hr-work-record` CRUD (5) | 勤怠の一覧〜削除 |
| 4 | `hr-time-clock` (4) | 打刻の一覧・詳細・登録・可能種別 |
| 5 | `hr-attendance-tag` CRUD (5) | 所定勤務タグ |
| 6 | `hr-salary` / `hr-bonus` (4) | 給与・賞与明細 |
| 7 | `hr-group` CRUD (4) | グループ |
| 8 | `hr-position` (3) | 役職 |

### v0.7.1 — HR 従業員ルール (Phase 5b, 17 commands)

従業員のサブリソース（口座、基本給、扶養、保険等）の show/update。

| # | 項目 | 説明 |
|---|------|------|
| 1 | `hr-employee bank-account` show/update | 振込先口座 |
| 2 | `hr-employee basic-pay` show/update | 基本給 |
| 3 | `hr-employee dependents` show/update | 扶養親族 |
| 4 | `hr-employee health-insurance` show/update | 健康保険 |
| 5 | `hr-employee welfare-pension` show/update | 厚生年金 |
| 6 | `hr-employee profile` show/update | プロフィール |
| 7 | `hr-employee work-record-summary` show | 勤怠サマリー |
| 8 | `hr-employee tax-withholding` show | 源泉徴収 |
| 9 | `hr-employee social-insurance` show/update | 社会保険 |
| 10 | `hr-employee employment-insurance` show | 雇用保険 |

### v0.7.2 — HR 承認ワークフロー (Phase 5c, 33 commands)

5種の勤怠承認ワークフロー（各 CRUD + action）と承認経路。

| # | 項目 | 説明 |
|---|------|------|
| 1 | `hr-approval-monthly` CRUD+action (6) | 月次勤怠承認 |
| 2 | `hr-approval-overtime` CRUD+action (6) | 残業承認 |
| 3 | `hr-approval-paid-leave` CRUD+action (6) | 有休承認 |
| 4 | `hr-approval-special-leave` CRUD+action (6) | 特別休暇承認 |
| 5 | `hr-approval-work-time` CRUD+action (6) | 勤務時間承認 |
| 6 | `hr-approval-route` list/show (2) | HR承認経路 |
| 7 | `hr-user me` (1) | HR現在ユーザー情報 |

### v0.7.3 — HR 年末調整 (Phase 5d, 13 commands)

年末調整（年調）関連の参照コマンド。

| # | 項目 | 説明 |
|---|------|------|
| 1 | `hr-yearend employees` | 年末調整対象者一覧 |
| 2 | `hr-yearend dependents` | 扶養控除 |
| 3 | `hr-yearend housing-loans` | 住宅ローン控除 |
| 4 | `hr-yearend insurances` | 保険料控除 |
| 5 | `hr-yearend life-insurances` | 生命保険料控除 |
| 6 | `hr-yearend social-insurances` | 社会保険料控除 |
| 7 | `hr-yearend earthquake-insurances` | 地震保険料控除 |
| 8 | `hr-yearend payroll` | 給与所得 |
| 9 | `hr-yearend previous-jobs` | 前職情報 |
| 10 | `hr-yearend base` / `status` / `result` / `summary` | 基本情報・ステータス・計算結果・サマリー |

### v0.8.0 — Invoice + PM + Sales API (Phase 6, 39 commands)

Accounting 以外の3つの API を一括実装。各 API は独立した Base URL を持つ。

| # | 項目 | 説明 |
|---|------|------|
| 1 | Invoice API クライアント | Base URL: `/iv/api/v1/` — `iv-` prefix |
| 2 | `iv-invoice` CRUD (5) | 請求書（Invoice API） |
| 3 | `iv-quotation` CRUD (5) | 見積書（Invoice API） |
| 4 | `iv-delivery` CRUD+list (5) | 納品書 |
| 5 | `iv-template list` (1) | テンプレート一覧 |
| 6 | PM API クライアント | Base URL: `/pm/api/v1/` — `pm-` prefix |
| 7 | `pm-project` CRUD (4) | プロジェクト |
| 8 | `pm-workload` CRUD (4) | 工数 |
| 9 | `pm-team` list/show (2) | チーム |
| 10 | Sales API クライアント | Base URL: `/sm/api/v1/` — `sales-` prefix |
| 11 | `sales-business` CRUD (4) | 案件 |
| 12 | `sales-order` CRUD (4) | 受注 |
| 13 | `sales-customer` list/show (2) | 顧客 |
| 14 | `sales-product` list/show (2) | 商品 |

### v0.8.1 — メタ情報 + 完了

| # | 項目 | 説明 |
|---|------|------|
| 1 | `api-resources` | 対応 API リソース一覧表示 |
| 2 | `forms selectables` | フォーム選択肢取得 |
| 3 | 全 API ドキュメント更新 | README + コマンドリファレンス |
| 4 | GoReleaser 全プラットフォーム対応 | darwin/linux/windows × amd64/arm64 |

---

## アーキテクチャノート

### Multi-API Base URL

freee は API ごとに異なる Base URL を使用する:

| API | Base URL | CLI prefix |
|-----|----------|------------|
| Accounting | `https://api.freee.co.jp/api/1/` | (なし) |
| HR | `https://api.freee.co.jp/hr/api/v1/` | `hr-` |
| Invoice | `https://api.freee.co.jp/iv/api/v1/` | `iv-` |
| PM | `https://api.freee.co.jp/pm/api/v1/` | `pm-` |
| Sales | `https://api.freee.co.jp/sm/api/v1/` | `sales-` |

`internal/api/client.go` で API 種別ごとに Base URL を切り替える。

### ネーミング規則

- **Accounting**: prefix なし（`freee deal`, `freee invoice`, `freee partner`）
- **その他 API**: API prefix 付き（`freee hr-employee`, `freee iv-invoice`, `freee pm-project`）
- **理由**: Accounting の `invoice` と Invoice API の `iv-invoice` は別リソース（前者は `/api/1/invoices`、後者は `/iv/api/v1/invoices`）

### Sub-resource パターン

親リソースの ID を引数に取るネストされたサブコマンド:

```
freee deal payment create <deal-id>           # 取引決済
freee hr-employee bank-account show <emp-id>  # 従業員口座
freee hr-yearend dependents <emp-id>          # 年末調整扶養
```

Cobra の `RunE` で `args[0]` を親 ID として取得する。

### Action コマンド

承認/却下などのアクションは `action` サブコマンドで統一:

```
freee approval action <id> --action approve
freee hr-approval-monthly action <id> --action reject --comment "要修正"
```
