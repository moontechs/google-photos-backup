package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google-backup/internal/auth"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type googleRedirectUrlHandler struct {
	googleAuth auth.Auth
}

func NewGoogleRedirectUrlHandler(googleAuth auth.Auth) *googleRedirectUrlHandler {
	return &googleRedirectUrlHandler{googleAuth: googleAuth}
}

func (h *googleRedirectUrlHandler) Handle(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{})

		return
	}

	var client struct {
		ClientID string `uri:"clientId" binding:"required"`
	}

	if err := c.ShouldBindUri(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	redirctUrl, err := h.googleAuth.GetRedirectUrl(client.ClientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		log.Error(fmt.Errorf("google redirect url: %w", err))

		return
	}

	rawMessage := json.RawMessage([]byte(`{"data": "` + redirctUrl + `"}`))

	// needs to return unescaped string to be able to redirect or open in a browser
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(rawMessage))
}
