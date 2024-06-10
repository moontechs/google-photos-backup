package account

import (
	"encoding/json"
	"fmt"
	"google-backup/internal/google_client"

	"golang.org/x/oauth2"
)

type Account interface {
	GetAccountOauthClientId(email string) (string, error)
	GetTokenByEmail(email string) (oauth2.Token, error)
}

type AccountData struct {
	Email      string `json:"email"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	Picture    string `json:"picture"`
}

type account struct {
	repository             Repository
	googleClientRepository google_client.Repository
}

func NewAccount(repository Repository, googleClientRepository google_client.Repository) account {
	return account{repository: repository, googleClientRepository: googleClientRepository}
}

func (a account) GetAccountOauthClientId(email string) (string, error) {
	assignedAccountsMap, err := a.googleClientRepository.FindAllAssignedAccounts()
	if err != nil {
		return "", fmt.Errorf("find all assigned accounts: %w", err)
	}

	for clientId, assignedAccountsJson := range assignedAccountsMap {
		var assignedAccounts []string

		err = json.Unmarshal(assignedAccountsJson, &assignedAccounts)
		if err != nil {
			return "", fmt.Errorf("unmarshal assigned accounts: %w", err)
		}

		for _, assignedAccount := range assignedAccounts {
			if assignedAccount == email {
				return clientId, nil
			}
		}
	}

	return "", fmt.Errorf("account oauth client id is not assigned")
}

func (a account) GetTokenByEmail(email string) (oauth2.Token, error) {
	token, err := a.repository.FindTokenByEmail(email)
	if err != nil {
		return oauth2.Token{}, fmt.Errorf("get token by email: %w", err)
	}

	if token == nil {
		return oauth2.Token{}, fmt.Errorf("account token is not assigned")
	}

	var authToken oauth2.Token
	err = json.Unmarshal(token, &authToken)
	if err != nil {
		return oauth2.Token{}, fmt.Errorf("oauth token unmarshal: %w", err)
	}

	return authToken, nil
}
