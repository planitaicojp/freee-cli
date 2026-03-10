package model

type InvoicesResponse struct {
	Invoices []Invoice `json:"invoices"`
}

type InvoiceResponse struct {
	Invoice Invoice `json:"invoice"`
}

type Invoice struct {
	ID               int64  `json:"id"`
	CompanyID        int64  `json:"company_id"`
	IssueDate        string `json:"issue_date"`
	DueDate          string `json:"due_date"`
	InvoiceNumber    string `json:"invoice_number"`
	PartnerID        int64  `json:"partner_id"`
	PartnerName      string `json:"partner_name"`
	PartnerCode      string `json:"partner_code"`
	InvoiceStatus    string `json:"invoice_status"`
	TotalAmount      int64  `json:"total_amount"`
	TotalVat         int64  `json:"total_vat"`
	SubTotal         int64  `json:"sub_total"`
	Title            string `json:"title"`
	PaymentType      string `json:"payment_type"`
	InvoiceLayout    string `json:"invoice_layout"`
	BookingDate      string `json:"booking_date"`
	Description      string `json:"description"`
}

type InvoiceRow struct {
	ID          int64  `json:"id"`
	Number      string `json:"number"`
	Partner     string `json:"partner"`
	Amount      int64  `json:"amount"`
	Status      string `json:"status"`
	IssueDate   string `json:"issue_date"`
}
