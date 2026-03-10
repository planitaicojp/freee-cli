package model

type TagsResponse struct {
	Tags []Tag `json:"tags"`
}

type Tag struct {
	ID        int64  `json:"id"`
	CompanyID int64  `json:"company_id"`
	Name      string `json:"name"`
	Shortcut1 string `json:"shortcut1"`
	Shortcut2 string `json:"shortcut2"`
	UpdateDate string `json:"update_date"`
}

type TagRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
