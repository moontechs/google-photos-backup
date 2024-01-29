package google_client

import "google-backup/internal/account"

type ClientData struct {
	ID               string                `json:"id"`
	Secret           string                `json:"secret"`
	RedirectURL      string                `json:"redirectUrl"`
	AssignedAccounts []account.AccountData `json:"assignedAccounts"`
}
