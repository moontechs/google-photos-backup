package google_client

type ClientData struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	RedirectURL string `json:"redirect_url"`
}

type AssignedAccountsData struct {
	Accounts []string `json:"accounts"`
}
