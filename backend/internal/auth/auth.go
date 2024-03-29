package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"google-backup/internal/google_client"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Auth
type Auth interface {
	GetRedirectUrl(clientId string) (string, error)
	GetToken(ctx context.Context, clientName, code string) (*oauth2.Token, error)
	GetHttpClient(ctx context.Context, clientId string, token *oauth2.Token) (*http.Client, error)
	GetUserInfo(client *http.Client) (UserInfo, error)
	SaveOauthClientData(email string, oauthClientData OauthClientData) error
	GetOauthClientData(email string) (OauthClientData, error)
}

type googleAuth struct {
	repository             Repository
	googleClientRepository google_client.Repository
}

type UserInfo struct {
	Picture    string `json:"picture"`
	Email      string `json:"email"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
}

type OauthClientData struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectURL  string `json:"redirectUrl"`
}

func NewGoogleAuth(repository Repository, googleClientRepository google_client.Repository) googleAuth {
	return googleAuth{
		repository:             repository,
		googleClientRepository: googleClientRepository,
	}
}

func (g googleAuth) GetRedirectUrl(clientId string) (string, error) {
	gConfig, err := g.createConfig(clientId)
	if err != nil {
		return "", fmt.Errorf("create config during create redirect url: %w", err)
	}

	return gConfig.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

func (g googleAuth) GetToken(ctx context.Context, clientId, code string) (*oauth2.Token, error) {
	gConfig, err := g.createConfig(clientId)
	if err != nil {
		return nil, fmt.Errorf("create config: %v", err)
	}

	token, err := gConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange token: %v", err)
	}

	return token, nil
}

func (g googleAuth) GetHttpClient(ctx context.Context, clientId string, token *oauth2.Token) (*http.Client, error) {
	gConfig, err := g.createConfig(clientId)
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

	var response struct {
		Email      string `json:"email"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Picture    string `json:"picture"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return UserInfo{}, fmt.Errorf("unmarshal user info: %w", err)
	}

	return UserInfo{
		Email:      response.Email,
		GivenName:  response.GivenName,
		FamilyName: response.FamilyName,
		Picture:    response.Picture,
	}, nil
}

func (g googleAuth) SaveOauthClientData(email string, oauthClientData OauthClientData) error {

	return nil
}

func (g googleAuth) GetOauthClientData(email string) (OauthClientData, error) {
	return OauthClientData{}, nil
}

func (g googleAuth) createConfig(clientId string) (oauth2.Config, error) {
	client, err := g.googleClientRepository.Find(clientId)
	if err != nil {
		return oauth2.Config{}, fmt.Errorf("find google client data: %w", err)
	}

	if client == nil {
		return oauth2.Config{}, fmt.Errorf("find google client data")
	}

	var googleClientData google_client.ClientData
	err = json.Unmarshal(client, &googleClientData)
	if err != nil {
		return oauth2.Config{}, fmt.Errorf("unmarshal google client data: %w", err)
	}

	return oauth2.Config{
		ClientID:     googleClientData.ID,
		ClientSecret: googleClientData.Secret,
		RedirectURL:  googleClientData.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/photoslibrary.readonly",
			"https://www.googleapis.com/auth/drive.readonly",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}, nil
}
