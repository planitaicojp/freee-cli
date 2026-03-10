package model

type ItemsResponse struct {
	Items []Item `json:"items"`
}

type Item struct {
	ID         int64  `json:"id"`
	CompanyID  int64  `json:"company_id"`
	Name       string `json:"name"`
	Available  bool   `json:"available"`
	Shortcut1  string `json:"shortcut1"`
	Shortcut2  string `json:"shortcut2"`
	UpdateDate string `json:"update_date"`
}

type ItemRow struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Available bool   `json:"available"`
}
