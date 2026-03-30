package model

type AccountItemsResponse struct {
	AccountItems []AccountItem `json:"account_items"`
}

type AccountItemResponse struct {
	AccountItem AccountItem `json:"account_item"`
}

type AccountItem struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	AccountCategory string `json:"account_category"`
	GroupName       string `json:"group_name"`
	TaxCode         int    `json:"tax_code"`
	DefaultTaxCode  int    `json:"default_tax_code"`
	Shortcut        string `json:"shortcut"`
	ShortcutNum     string `json:"shortcut_num"`
	Available       bool   `json:"available"`
	UpdateDate      string `json:"update_date"`
}

func (a AccountItem) GetID() int64    { return a.ID }
func (a AccountItem) GetName() string { return a.Name }

// AccountItemRow is a display-friendly row for table output.
type AccountItemRow struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	TaxCode  int    `json:"tax_code"`
	Shortcut string `json:"shortcut"`
}
