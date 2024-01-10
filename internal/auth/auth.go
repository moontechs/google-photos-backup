package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type Auth interface {
	GetRedirectUrl(clientName string) string
	GetToken(ctx context.Context, clientName, code string) (*oauth2.Token, error)
	GetClient(ctx context.Context, clientName string, token *oauth2.Token) (*http.Client, error)
	GetUserInfo(client *http.Client) (UserInfo, error)
	SaveOauthClientData(email string, oauthClientData OauthClientData) error
	GetOauthClientData(email string) (OauthClientData, error)
}

type googleAuth struct {
	config     oauth2.Config
	repository Repository
}

type UserInfo struct {
	Picture string `json:"picture"`
	Email   string `json:"email"`
}

type OauthClientData struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

func NewGoogleAuth(repository Repository) googleAuth {
	return googleAuth{
		repository: repository,
	}
}

func (g googleAuth) GetRedirectUrl(clientName string) string {
	gConfig, err := g.createConfig(clientName)
	if err != nil {
		return "/error?type=oauth_client_data_not_found"
	}

	return gConfig.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (g googleAuth) GetToken(ctx context.Context, clientName, code string) (*oauth2.Token, error) {
	gConfig, err := g.createConfig(clientName)
	if err != nil {
		return nil, fmt.Errorf("create config: %v", err)
	}

	token, err := gConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange token: %v", err)
	}

	return token, nil
}

func (g googleAuth) GetClient(ctx context.Context, clientName string, token *oauth2.Token) (*http.Client, error) {
	gConfig, err := g.createConfig(clientName)
	if err != nil {
		return nil, fmt.Errorf("create config: %v", err)
	}

	return gConfig.Client(ctx, token), nil
}

func (g googleAuth) GetUserInfo(client *http.Client) (UserInfo, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return UserInfo{}, fmt.Errorf("get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UserInfo{}, fmt.Errorf("get user info: %w", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserInfo{}, fmt.Errorf("read user info body: %w", err)
	}

	var result UserInfo
	if err := json.Unmarshal(data, &result); err != nil {
		return UserInfo{}, fmt.Errorf("unmarshal user info: %w", err)
	}

	return result, nil
}

func (g googleAuth) SaveOauthClientData(email string, oauthClientData OauthClientData) error {

	return nil
}

func (g googleAuth) GetOauthClientData(email string) (OauthClientData, error) {
	return OauthClientData{}, nil
}

func (g googleAuth) createConfig(clientName string) (oauth2.Config, error) {
	oauthClientDataJson, err := g.repository.GetOauthClientData(clientName)
	if err != nil {
		return oauth2.Config{}, fmt.Errorf("get oauth client data: %w", err)
	}

	if oauthClientDataJson == nil {
		return oauth2.Config{}, fmt.Errorf("oauth client data not found: %w", err)
	}

	var oauthClientData OauthClientData
	err = json.Unmarshal(oauthClientDataJson, &oauthClientData)
	if err != nil {
		return oauth2.Config{}, fmt.Errorf("unmarshal oauth client data: %w", err)
	}

	return oauth2.Config{
		ClientID:     oauthClientData.ClientID,
		ClientSecret: oauthClientData.ClientSecret,
		RedirectURL:  oauthClientData.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/photoslibrary.readonly",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: oauth2.Endpoint{},
	}, nil
}
