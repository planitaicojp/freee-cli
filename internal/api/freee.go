package api

import "fmt"

// FreeeAPI wraps the freee Accounting API.
type FreeeAPI struct {
	Client *Client
}

func (a *FreeeAPI) url(path string, args ...any) string {
	return a.Client.BaseURL() + "/api/1" + fmt.Sprintf(path, args...)
}

// --- Companies ---

func (a *FreeeAPI) ListCompanies(result any) error {
	return a.Client.Get(a.url("/companies"), result)
}

func (a *FreeeAPI) GetCompany(id int64, result any) error {
	return a.Client.Get(a.url("/companies/%d", id), result)
}

// --- Deals ---

func (a *FreeeAPI) ListDeals(companyID int64, params string, result any) error {
	url := a.url("/deals?company_id=%d", companyID)
	if params != "" {
		url += "&" + params
	}
	return a.Client.Get(url, result)
}

func (a *FreeeAPI) GetDeal(companyID, dealID int64, result any) error {
	return a.Client.Get(a.url("/deals/%d?company_id=%d", dealID, companyID), result)
}

func (a *FreeeAPI) CreateDeal(body, result any) error {
	_, err := a.Client.Post(a.url("/deals"), body, result)
	return err
}

func (a *FreeeAPI) UpdateDeal(dealID int64, body, result any) error {
	return a.Client.Put(a.url("/deals/%d", dealID), body, result)
}

func (a *FreeeAPI) DeleteDeal(companyID, dealID int64) error {
	return a.Client.Delete(a.url("/deals/%d?company_id=%d", dealID, companyID))
}

// --- Invoices ---

func (a *FreeeAPI) ListInvoices(companyID int64, params string, result any) error {
	url := a.url("/invoices?company_id=%d", companyID)
	if params != "" {
		url += "&" + params
	}
	return a.Client.Get(url, result)
}

func (a *FreeeAPI) GetInvoice(companyID, invoiceID int64, result any) error {
	return a.Client.Get(a.url("/invoices/%d?company_id=%d", invoiceID, companyID), result)
}

func (a *FreeeAPI) CreateInvoice(body, result any) error {
	_, err := a.Client.Post(a.url("/invoices"), body, result)
	return err
}

func (a *FreeeAPI) UpdateInvoice(invoiceID int64, body, result any) error {
	return a.Client.Put(a.url("/invoices/%d", invoiceID), body, result)
}

func (a *FreeeAPI) DeleteInvoice(companyID, invoiceID int64) error {
	return a.Client.Delete(a.url("/invoices/%d?company_id=%d", invoiceID, companyID))
}

// --- Partners ---

func (a *FreeeAPI) ListPartners(companyID int64, params string, result any) error {
	url := a.url("/partners?company_id=%d", companyID)
	if params != "" {
		url += "&" + params
	}
	return a.Client.Get(url, result)
}

func (a *FreeeAPI) GetPartner(companyID, partnerID int64, result any) error {
	return a.Client.Get(a.url("/partners/%d?company_id=%d", partnerID, companyID), result)
}

func (a *FreeeAPI) CreatePartner(body, result any) error {
	_, err := a.Client.Post(a.url("/partners"), body, result)
	return err
}

func (a *FreeeAPI) UpdatePartner(partnerID int64, body, result any) error {
	return a.Client.Put(a.url("/partners/%d", partnerID), body, result)
}

func (a *FreeeAPI) DeletePartner(companyID, partnerID int64) error {
	return a.Client.Delete(a.url("/partners/%d?company_id=%d", partnerID, companyID))
}

// --- Account Items ---

func (a *FreeeAPI) ListAccountItems(companyID int64, result any) error {
	return a.Client.Get(a.url("/account_items?company_id=%d", companyID), result)
}

func (a *FreeeAPI) GetAccountItem(companyID, itemID int64, result any) error {
	return a.Client.Get(a.url("/account_items/%d?company_id=%d", itemID, companyID), result)
}

// --- Sections ---

func (a *FreeeAPI) ListSections(companyID int64, result any) error {
	return a.Client.Get(a.url("/sections?company_id=%d", companyID), result)
}

func (a *FreeeAPI) CreateSection(body, result any) error {
	_, err := a.Client.Post(a.url("/sections"), body, result)
	return err
}

func (a *FreeeAPI) UpdateSection(sectionID int64, body, result any) error {
	return a.Client.Put(a.url("/sections/%d", sectionID), body, result)
}

func (a *FreeeAPI) DeleteSection(companyID, sectionID int64) error {
	return a.Client.Delete(a.url("/sections/%d?company_id=%d", sectionID, companyID))
}

// --- Tags ---

func (a *FreeeAPI) ListTags(companyID int64, result any) error {
	return a.Client.Get(a.url("/tags?company_id=%d", companyID), result)
}

func (a *FreeeAPI) CreateTag(body, result any) error {
	_, err := a.Client.Post(a.url("/tags"), body, result)
	return err
}

func (a *FreeeAPI) UpdateTag(tagID int64, body, result any) error {
	return a.Client.Put(a.url("/tags/%d", tagID), body, result)
}

func (a *FreeeAPI) DeleteTag(companyID, tagID int64) error {
	return a.Client.Delete(a.url("/tags/%d?company_id=%d", tagID, companyID))
}

// --- Items ---

func (a *FreeeAPI) ListItems(companyID int64, result any) error {
	return a.Client.Get(a.url("/items?company_id=%d", companyID), result)
}

func (a *FreeeAPI) CreateItem(body, result any) error {
	_, err := a.Client.Post(a.url("/items"), body, result)
	return err
}

func (a *FreeeAPI) UpdateItem(itemID int64, body, result any) error {
	return a.Client.Put(a.url("/items/%d", itemID), body, result)
}

func (a *FreeeAPI) DeleteItem(companyID, itemID int64) error {
	return a.Client.Delete(a.url("/items/%d?company_id=%d", itemID, companyID))
}

// --- Journals ---

func (a *FreeeAPI) ListJournals(companyID int64, params string, result any) error {
	url := a.url("/journals?company_id=%d", companyID)
	if params != "" {
		url += "&" + params
	}
	return a.Client.Get(url, result)
}

// --- Expense Applications ---

func (a *FreeeAPI) ListExpenseApplications(companyID int64, params string, result any) error {
	url := a.url("/expense_applications?company_id=%d", companyID)
	if params != "" {
		url += "&" + params
	}
	return a.Client.Get(url, result)
}

func (a *FreeeAPI) GetExpenseApplication(companyID, expenseID int64, result any) error {
	return a.Client.Get(a.url("/expense_applications/%d?company_id=%d", expenseID, companyID), result)
}

func (a *FreeeAPI) CreateExpenseApplication(body, result any) error {
	_, err := a.Client.Post(a.url("/expense_applications"), body, result)
	return err
}

func (a *FreeeAPI) UpdateExpenseApplication(expenseID int64, body, result any) error {
	return a.Client.Put(a.url("/expense_applications/%d", expenseID), body, result)
}

func (a *FreeeAPI) DeleteExpenseApplication(companyID, expenseID int64) error {
	return a.Client.Delete(a.url("/expense_applications/%d?company_id=%d", expenseID, companyID))
}

// --- Walletables ---

func (a *FreeeAPI) ListWalletables(companyID int64, result any) error {
	return a.Client.Get(a.url("/walletables?company_id=%d", companyID), result)
}

func (a *FreeeAPI) GetWalletable(companyID int64, walletableType string, walletableID int64, result any) error {
	return a.Client.Get(a.url("/walletables/%s/%d?company_id=%d", walletableType, walletableID, companyID), result)
}

// --- Manual Journals ---

func (a *FreeeAPI) ListManualJournals(companyID int64, params string, result any) error {
	url := a.url("/manual_journals?company_id=%d", companyID)
	if params != "" {
		url += "&" + params
	}
	return a.Client.Get(url, result)
}

func (a *FreeeAPI) GetManualJournal(companyID, id int64, result any) error {
	return a.Client.Get(a.url("/manual_journals/%d?company_id=%d", id, companyID), result)
}

func (a *FreeeAPI) CreateManualJournal(body, result any) error {
	_, err := a.Client.Post(a.url("/manual_journals"), body, result)
	return err
}

func (a *FreeeAPI) UpdateManualJournal(id int64, body, result any) error {
	return a.Client.Put(a.url("/manual_journals/%d", id), body, result)
}

func (a *FreeeAPI) DeleteManualJournal(companyID, id int64) error {
	return a.Client.Delete(a.url("/manual_journals/%d?company_id=%d", id, companyID))
}

// --- Transfers ---

func (a *FreeeAPI) ListTransfers(companyID int64, params string, result any) error {
	url := a.url("/transfers?company_id=%d", companyID)
	if params != "" {
		url += "&" + params
	}
	return a.Client.Get(url, result)
}

func (a *FreeeAPI) GetTransfer(companyID, id int64, result any) error {
	return a.Client.Get(a.url("/transfers/%d?company_id=%d", id, companyID), result)
}

func (a *FreeeAPI) CreateTransfer(body, result any) error {
	_, err := a.Client.Post(a.url("/transfers"), body, result)
	return err
}

func (a *FreeeAPI) UpdateTransfer(id int64, body, result any) error {
	return a.Client.Put(a.url("/transfers/%d", id), body, result)
}

func (a *FreeeAPI) DeleteTransfer(companyID, id int64) error {
	return a.Client.Delete(a.url("/transfers/%d?company_id=%d", id, companyID))
}
