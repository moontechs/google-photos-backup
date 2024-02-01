package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"google-backup/internal/account"
	"google-backup/internal/auth"
	"google-backup/internal/google_client"
	"google-backup/internal/settings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type googleCallbackHandler struct {
	accountRepository      account.Repository
	googleClientRepository google_client.Repository
	settingsRepository     settings.Repository
	googleAuth             auth.Auth
}

func NewGoogleCallbackHandler(
	accountRepository account.Repository,
	googleClientRepository google_client.Repository,
	settingsRepository settings.Repository,
	googleAuth auth.Auth,
) *googleCallbackHandler {
	return &googleCallbackHandler{
		accountRepository:      accountRepository,
		googleClientRepository: googleClientRepository,
		settingsRepository:     settingsRepository,
		googleAuth:             googleAuth,
	}
}

func (h *googleCallbackHandler) Handle(c *gin.Context) {
	var client struct {
		ClientID string `uri:"clientId" binding:"required"`
	}

	if err := c.ShouldBindUri(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	token, err := h.googleAuth.GetToken(c.Request.Context(), client.ClientID, c.Query("code"))
	if err != nil {
		log.Error(err)

		return
	}

	httpClient, err := h.googleAuth.GetHttpClient(c.Request.Context(), client.ClientID, token)
	if err != nil {
		log.Error(fmt.Errorf("get http client: %w", err))

		c.String(http.StatusInternalServerError, "Could not get Google http client")

		return
	}

	userInfo, err := h.googleAuth.GetUserInfo(httpClient)
	if err != nil {
		log.Error(fmt.Errorf("get user info: %w", err))

		c.String(http.StatusInternalServerError, "Could not get user info")

		return
	}

	tokenData, err := json.Marshal(token)
	if err != nil {
		log.Error(fmt.Errorf("marshal token: %w", err))

		c.String(http.StatusInternalServerError, "Could not marshal Google token")

		return
	}

	err = h.accountRepository.SaveToken(userInfo.Email, tokenData)
	if err != nil {
		log.Error(fmt.Errorf("save token: %w", err))

		c.String(http.StatusInternalServerError, "Could not save Google token")

		return
	}

	err = h.assigneAccountToClient(client.ClientID, userInfo.Email)
	if err != nil {
		log.Error(fmt.Errorf("assign account to client: %w", err))

		c.String(http.StatusInternalServerError, "Could not assign account to client")

		return
	}

	userInfoJson, err := json.Marshal(userInfo)
	if err != nil {
		log.Error(fmt.Errorf("marshal user info: %w", err))

		c.String(http.StatusInternalServerError, "Could not marshal user info")

		return
	}

	err = h.accountRepository.SaveAccount(userInfo.Email, userInfoJson)
	if err != nil {
		log.Error(fmt.Errorf("save account: %w", err))

		c.String(http.StatusInternalServerError, "Could not save user info")
	}

	host, err := h.getHostFromSettings()
	if err != nil {
		log.Error(fmt.Errorf("get host from settings: %w", err))

		c.String(http.StatusInternalServerError, "Host was not found in the settings")

		return
	}

	c.Redirect(http.StatusFound, host)
}

func (h *googleCallbackHandler) assigneAccountToClient(clientId string, email string) error {
	err := h.unassignAccountFromOtherClients(clientId, email)
	if err != nil {
		return fmt.Errorf("unassign account from other clients: %w", err)
	}

	emails, err := h.googleClientRepository.FindAssignedAccounts(clientId)
	if err != nil {
		return fmt.Errorf("find assigned accounts: %w", err)
	}

	if emails == nil {
		emails = []byte("[]")
	}

	var emailsSlice []string
	err = json.Unmarshal(emails, &emailsSlice)
	if err != nil {
		return fmt.Errorf("unmarshal emails: %w", err)
	}

	emailsSlice = append(emailsSlice, email)
	emailsSlice = slices.Compact(emailsSlice)

	emails, err = json.Marshal(emailsSlice)
	if err != nil {
		return fmt.Errorf("marshal emails: %w", err)
	}

	err = h.googleClientRepository.SaveAssignedAccounts(clientId, emails)
	if err != nil {
		return fmt.Errorf("save assigned accounts: %w", err)
	}

	return nil
}

func (h *googleCallbackHandler) unassignAccountFromOtherClients(clientId, email string) error {
	accountsJsonSlice, err := h.googleClientRepository.FindAllAssignedAccounts()
	if err != nil {
		return fmt.Errorf("find all assigned accounts: %w", err)
	}

	for clientIdOfJsonSlice, accountsJson := range accountsJsonSlice {
		if clientIdOfJsonSlice == clientId {
			continue
		}

		var accounts []string
		err = json.Unmarshal(accountsJson, &accounts)
		if err != nil {
			return fmt.Errorf("unmarshal accounts: %w", err)
		}

		for i, accountEmail := range accounts {
			if accountEmail != email {
				continue
			}

			accounts = append(accounts[:i], accounts[i+1:]...)

			accountsJson, err := json.Marshal(accounts)
			if err != nil {
				return fmt.Errorf("marshal accounts: %w", err)
			}

			h.googleClientRepository.SaveAssignedAccounts(clientIdOfJsonSlice, accountsJson)

			return nil
		}
	}

	return nil
}

func (h *googleCallbackHandler) getHostFromSettings() (string, error) {
	settingsJson, err := h.settingsRepository.Find()
	if err != nil {
		return "", fmt.Errorf("find settings: %w", err)
	}

	if settingsJson == nil {
		return "", fmt.Errorf("settings not found")
	}

	var settingsData settings.SettingsData

	err = json.Unmarshal(settingsJson, &settingsData)
	if err != nil {
		return "", fmt.Errorf("unmarshal settings: %w", err)
	}

	return settingsData.Host, nil
}
