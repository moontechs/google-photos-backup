package handlers

import (
	"encoding/json"
	"net/http"

	"google-backup/internal/account"
	"google-backup/internal/auth"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type googleCallbackHandler struct {
	accountRepository account.Repository
	googleAuth        auth.Auth
}

func NewGoogleCallbackHandler(accountRepository account.Repository, googleAuth auth.Auth) *googleCallbackHandler {
	return &googleCallbackHandler{accountRepository: accountRepository, googleAuth: googleAuth}
}

func (h *googleCallbackHandler) Handle(c *gin.Context) {
	var client struct {
		ClientName string `uri:"clientName" binding:"required"`
	}

	if err := c.ShouldBindUri(&client); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})

		return
	}

	token, err := h.googleAuth.GetToken(c.Request.Context(), client.ClientName, c.Query("code"))
	if err != nil {
		log.Error(err)

		return
	}

	httpClient, err := h.googleAuth.GetClient(c.Request.Context(), client.ClientName, token)
	if err != nil {
		log.Error(err)

		return
	}

	userInfo, err := h.googleAuth.GetUserInfo(httpClient)
	if err != nil {
		log.Error(err)

		return
	}

	tokenData, err := json.Marshal(token)
	if err != nil {
		log.Error(err)

		return
	}

	err = h.accountRepository.SaveToken(userInfo.Email, tokenData)
	if err != nil {
		log.Error(err)

		return
	}

	c.String(http.StatusOK, "Sucess! Now you can close this tab.") // TODO: redirect to frontend
}
