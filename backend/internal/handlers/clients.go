package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google-backup/internal/account"
	"google-backup/internal/google_client"
	"google-backup/internal/settings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type clientsApiHandler struct {
	accountRepository      account.Repository
	googleClientRepository google_client.Repository
	settingsRepository     settings.Repository
}

type updateClientRequest struct {
	ID     string `json:"id" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}

type clientDataResponse struct {
	google_client.ClientData
	AssignedAccounts []account.AccountData `json:"assignedAccounts"`
}

func NewClientsApiHandler(
	accountRepository account.Repository,
	googleClientRepository google_client.Repository,
	settingsRepository settings.Repository,
) *clientsApiHandler {
	return &clientsApiHandler{
		accountRepository:      accountRepository,
		googleClientRepository: googleClientRepository,
		settingsRepository:     settingsRepository,
	}
}

func (h *clientsApiHandler) Handle(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		h.handleGet(c)

		return
	case "POST":
		h.handlePost(c)

		return
	case "DELETE":
		h.handleDelete(c)

		return
	}

	c.JSON(http.StatusMethodNotAllowed, gin.H{})
}

func (h *clientsApiHandler) handleGet(c *gin.Context) {
	clientId := c.Param("clientId")

	if clientId == "" {
		clients, err := h.googleClientRepository.FindAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Error(fmt.Errorf("google clients: handle get all: find all: %w", err))

			return
		}

		clientsDataResponse := make([]clientDataResponse, 0, len(clients))

		for _, client := range clients {
			clientDataResponse, err := h.getClientData(client)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				log.Error(fmt.Errorf("google clients: handle get all: get client data: %w", err))

				return
			}

			clientsDataResponse = append(clientsDataResponse, clientDataResponse)
		}

		c.JSON(http.StatusOK, gin.H{"data": clientsDataResponse})

		return
	}

	client, err := h.googleClientRepository.Find(clientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		log.Error(fmt.Errorf("google clients: handle get: find: %w", err))

		return
	}

	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{})

		return
	}

	clientDataResponse, err := h.getClientData(client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		log.Error(fmt.Errorf("google clients: handle get: get client data: %w", err))

		return
	}

	c.JSON(http.StatusOK, gin.H{"data": clientDataResponse})
}

func (h *clientsApiHandler) handlePost(c *gin.Context) {
	// TODO: add edit client functionality
	updateClientRequestData := updateClientRequest{}

	err := c.ShouldBindJSON(&updateClientRequestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	redirectUrl, err := h.generateRedirectUrl(updateClientRequestData.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		log.Error(fmt.Errorf("google clients: handle post: generate redirect url: %w", err))

		return
	}

	clientData := google_client.ClientData{
		ID:          updateClientRequestData.ID,
		Secret:      updateClientRequestData.Secret,
		RedirectURL: redirectUrl,
	}

	client, err := json.Marshal(clientData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		log.Error(fmt.Errorf("google clients: handle post: marshal: %w", err))

		return
	}

	err = h.googleClientRepository.Save(clientData.ID, client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		log.Error(fmt.Errorf("google clients: handle post: save: %w", err))

		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": clientData})
}

func (h *clientsApiHandler) handleDelete(c *gin.Context) {
	var clientIdParam struct {
		ClientID string `uri:"clientId" binding:"required"`
	}
	if err := c.ShouldBindUri(&clientIdParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	clientId := c.Param("clientId")

	err := h.googleClientRepository.Delete(clientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		log.Error(fmt.Errorf("google clients: handle delete: %w", err))

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *clientsApiHandler) generateRedirectUrl(clientId string) (string, error) {
	settingsJson, err := h.settingsRepository.Find()
	if err != nil {
		return "", fmt.Errorf("find settings: %w", err)
	}

	var settingsData settings.SettingsData
	err = json.Unmarshal(settingsJson, &settingsData)
	if err != nil {
		return "", fmt.Errorf("marshal settings: %w", err)
	}

	return fmt.Sprintf("%s/auth/google/callback/%s", settingsData.Host, clientId), nil
}

func (h *clientsApiHandler) getClientData(client []byte) (clientDataResponse, error) {
	var clientData google_client.ClientData

	err := json.Unmarshal(client, &clientData)
	if err != nil {
		return clientDataResponse{}, fmt.Errorf("unmarshal client data: %w", err)
	}

	clientDataResponse := clientDataResponse{clientData, make([]account.AccountData, 0)}

	assignedAccountsJson, err := h.googleClientRepository.FindAssignedAccounts(clientData.ID)
	if err != nil {
		return clientDataResponse, fmt.Errorf("find assigned accounts: %w", err)
	}

	if assignedAccountsJson == nil {
		return clientDataResponse, nil
	}

	var assignedAccountEmails []string
	err = json.Unmarshal(assignedAccountsJson, &assignedAccountEmails)
	if err != nil {
		return clientDataResponse, fmt.Errorf("unmarshal assigned accounts: %w", err)
	}

	for _, email := range assignedAccountEmails {
		accountJson, err := h.accountRepository.FindAccount(email)
		if err != nil {
			return clientDataResponse, fmt.Errorf("find account: %w", err)
		}

		if accountJson == nil {
			continue
		}

		var accountData account.AccountData
		err = json.Unmarshal(accountJson, &accountData)
		if err != nil {
			return clientDataResponse, fmt.Errorf("unmarshal account data: %w", err)
		}

		clientDataResponse.AssignedAccounts = append(clientDataResponse.AssignedAccounts, accountData)
	}

	return clientDataResponse, nil
}
