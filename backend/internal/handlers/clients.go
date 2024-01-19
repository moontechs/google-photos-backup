package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google-backup/internal/google_client"
	"google-backup/internal/settings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type clientsApiHandler struct {
	googleClientRepository google_client.Repository
	settingsRepository     settings.Repository
}

type updateClientRequest struct {
	ID     string `json:"id" binding:"required,alphanum"`
	Secret string `json:"secret" binding:"required,alphanum"`
}

func NewClientsApiHandler(
	googleClientRepository google_client.Repository,
	settingsRepository settings.Repository,
) *clientsApiHandler {
	return &clientsApiHandler{
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

		clientsData := make([]google_client.ClientData, 0, len(clients))

		for _, client := range clients {
			var clientData google_client.ClientData
			err = json.Unmarshal(client, &clientData)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				log.Error(fmt.Errorf("google clients: handle get all: unmarshal: %w", err))

				return
			}

			clientsData = append(clientsData, clientData)
		}

		c.JSON(http.StatusOK, gin.H{"data": clientsData})

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

	var clientData google_client.ClientData
	err = json.Unmarshal(client, &clientData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		log.Error(fmt.Errorf("google clients: handle get: unmarshal: %w", err))

		return
	}

	c.JSON(http.StatusOK, gin.H{"data": clientData})
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

func (h *clientsApiHandler) generateRedirectUrl(clientID string) (string, error) {
	settingsJson, err := h.settingsRepository.Find()
	if err != nil {
		return "", fmt.Errorf("find settings: %w", err)
	}

	var settingsData settings.SettingsData
	err = json.Unmarshal(settingsJson, &settingsData)
	if err != nil {
		return "", fmt.Errorf("marshal settings: %w", err)
	}

	return fmt.Sprintf("%s/auth/google/callback/%s", settingsData.Host, clientID), nil
}
