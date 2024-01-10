package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moontechs/photos-backup/internal/account"
	"github.com/moontechs/photos-backup/internal/scanner"
)

type rescanHandler struct {
	accountRepository account.Repository
	scheduler         scanner.Scheduler
}

type requestData struct {
	Email string `json:"email" binding:"required"`
}

func NewRescanHandler(
	accountRepository account.Repository,
	scheduler scanner.Scheduler,
) *rescanHandler {
	return &rescanHandler{accountRepository: accountRepository, scheduler: scheduler}
}

func (h *rescanHandler) Handle(c *gin.Context) {
	var requestData requestData

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("json bind request: %v", err)

		return
	}

	exist, err := h.accountRepository.AccountExist(requestData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not check if account exists"})
		log.Printf("account exists: %v", err)

		return
	}

	if exist {
		h.scheduler.ScheduleRescan(requestData.Email)

		c.JSON(http.StatusOK, gin.H{})

		return
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
}
