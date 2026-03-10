package model

type CompanyResponse struct {
	Company Company `json:"company"`
}

type Company struct {
	ID                     int64        `json:"id"`
	Name                   string       `json:"name"`
	NameKana               string       `json:"name_kana"`
	DisplayName            string       `json:"display_name"`
	Role                   string       `json:"role"`
	Phone1                 string       `json:"phone1"`
	Phone2                 string       `json:"phone2"`
	Fax                    *string      `json:"fax"`
	Zipcode                string       `json:"zipcode"`
	PrefectureCode         int          `json:"prefecture_code"`
	StreetName1            string       `json:"street_name1"`
	StreetName2            string       `json:"street_name2"`
	CompanyNumber          string       `json:"company_number"`
	CorporateNumber        string       `json:"corporate_number"`
	PrivateSettlement      bool         `json:"private_settlement"`
	InvoiceLayout          string       `json:"invoice_layout"`
	OrgCode                int          `json:"org_code"`
	TxnNumberFormat        string       `json:"txn_number_format"`
	AmountFraction         int          `json:"amount_fraction"`
	WorkflowSetting        string       `json:"workflow_setting"`
	UsePartnerCode         bool         `json:"use_partner_code"`
	TaxAtSourceCalcType    int          `json:"tax_at_source_calc_type"`
	MinusFormat            int          `json:"minus_format"`
	DefaultWalletAccountID int64        `json:"default_wallet_account_id"`
	FiscalYears            []FiscalYear `json:"fiscal_years"`
}

type FiscalYear struct {
	ID        int64  `json:"id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
