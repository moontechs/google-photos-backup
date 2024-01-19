package account

import (
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
)

type Account interface {
	GetAccounts() ([][]byte, error)
	GetAccountOauthClientName(email string) (string, error)
	GetTokenByEmail(email string) (oauth2.Token, error)
}

type account struct {
	repository Repository
}

func NewAccount(repository Repository) account {
	return account{repository: repository}
}

func (a account) GetAccounts() ([][]byte, error) {
	return a.repository.GetAccounts()
}

func (a account) GetAccountOauthClientName(email string) (string, error) {
	clientName, err := a.repository.GetAccountOauthClientName(email)
	if err != nil {
		return "", fmt.Errorf("get account oauth client name: %w", err)
	}

	if clientName == nil {
		return "", fmt.Errorf("account oauth client name is not assigned")
	}

	return string(clientName), nil
}

func (a account) GetTokenByEmail(email string) (oauth2.Token, error) {
	token, err := a.repository.GetTokenByEmail(email)
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
