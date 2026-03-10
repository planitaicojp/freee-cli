package model

type SectionsResponse struct {
	Sections []Section `json:"sections"`
}

type Section struct {
	ID          int64  `json:"id"`
	CompanyID   int64  `json:"company_id"`
	Name        string `json:"name"`
	Available   bool   `json:"available"`
	LongName    string `json:"long_name"`
	Shortcut1   string `json:"shortcut1"`
	Shortcut2   string `json:"shortcut2"`
	IndentCount int    `json:"indent_count"`
	ParentID    *int64 `json:"parent_id"`
}

type SectionRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
