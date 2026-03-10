package model

type PartnersResponse struct {
	Partners []Partner `json:"partners"`
}

type PartnerResponse struct {
	Partner Partner `json:"partner"`
}

type Partner struct {
	ID          int64  `json:"id"`
	CompanyID   int64  `json:"company_id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Shortcut1   string `json:"shortcut1"`
	Shortcut2   string `json:"shortcut2"`
	LongName    string `json:"long_name"`
	OrgCode     int    `json:"org_code"`
	CountryCode string `json:"country_code"`
	Available   bool   `json:"available"`
	UpdateDate  string `json:"update_date"`
}

type PartnerRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
