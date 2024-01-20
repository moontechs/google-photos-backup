package handlers

import (
	"fmt"
	"net/http"

	"google-backup/internal/account"
	"google-backup/internal/scanner"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type rescanHandler struct {
	accountRepository account.Repository
	scheduler         scanner.Scheduler
}

type requestData struct {
	Type  string `json:"type" binding:"required,oneof=photos drive"`
	Email string `json:"email" binding:"required,email"`
}

func NewRescanHandler(
	accountRepository account.Repository,
	scheduler scanner.Scheduler,
) *rescanHandler {
	return &rescanHandler{accountRepository: accountRepository, scheduler: scheduler}
}

func (h *rescanHandler) Handle(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{})

		return
	}

	var requestData requestData

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		log.Error(fmt.Errorf("json bind request: %w", err))

		return
	}

	exist, err := h.accountRepository.AccountExist(requestData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not check if account exists"})
		log.Error(fmt.Errorf("account exists: %w", err))

		return
	}

	if exist {
		err = h.scheduler.ScheduleRescan(requestData.Type, requestData.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not schedule rescan"})
			log.Error(fmt.Errorf("schedule rescan: %w", err))

			return
		}

		c.JSON(http.StatusOK, gin.H{})

		return
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
}
