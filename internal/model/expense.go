package model

type ExpenseApplicationsResponse struct {
	ExpenseApplications []ExpenseApplication `json:"expense_applications"`
}

type ExpenseApplicationResponse struct {
	ExpenseApplication ExpenseApplication `json:"expense_application"`
}

type ExpenseApplication struct {
	ID               int64  `json:"id"`
	CompanyID        int64  `json:"company_id"`
	Title            string `json:"title"`
	IssueDate        string `json:"issue_date"`
	TotalAmount      int64  `json:"total_amount"`
	Status           string `json:"status"`
	Description      string `json:"description"`
	SectionID        *int64 `json:"section_id"`
	TagIDs           []int  `json:"tag_ids"`
	ApplicantID      int64  `json:"applicant_id"`
	ApplicationDate  string `json:"application_date"`
}

type ExpenseApplicationRow struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Amount int64  `json:"amount"`
	Status string `json:"status"`
	Date   string `json:"date"`
}
