package google_client

type ClientData struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	RedirectURL string `json:"redirectUrl"`
}

type AssignedAccountsData struct {
	Accounts []string `json:"accounts"`
}
