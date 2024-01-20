package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"google-backup/internal/account/accountfakes"
	"google-backup/internal/scanner/scannerfakes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRescanHandle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("request photos rescan", func(t *testing.T) {
		fakeScheduler := new(scannerfakes.FakeScheduler)
		fakeAccountRepository := new(accountfakes.FakeRepository)
		handler := NewRescanHandler(fakeAccountRepository, fakeScheduler)

		fakeAccountRepository.AccountExistReturns(true, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/rescan", bytes.NewBuffer(
			[]byte(`{"type":"photos","email":"test@gmail.com"}`),
		))

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)

		requestType, email := fakeScheduler.ScheduleRescanArgsForCall(0)
		assert.Equal(t, "photos", requestType)
		assert.Equal(t, "test@gmail.com", email)
	})

	t.Run("request drive rescan", func(t *testing.T) {
		fakeScheduler := new(scannerfakes.FakeScheduler)
		fakeAccountRepository := new(accountfakes.FakeRepository)
		handler := NewRescanHandler(fakeAccountRepository, fakeScheduler)

		fakeAccountRepository.AccountExistReturns(true, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/rescan", bytes.NewBuffer(
			[]byte(`{"type":"drive","email":"test@gmail.com"}`),
		))

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)

		requestType, email := fakeScheduler.ScheduleRescanArgsForCall(0)
		assert.Equal(t, "drive", requestType)
		assert.Equal(t, "test@gmail.com", email)
	})

	t.Run("request photos rescan account not found", func(t *testing.T) {
		fakeScheduler := new(scannerfakes.FakeScheduler)
		fakeAccountRepository := new(accountfakes.FakeRepository)
		handler := NewRescanHandler(fakeAccountRepository, fakeScheduler)

		fakeAccountRepository.AccountExistReturns(false, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/rescan", bytes.NewBuffer(
			[]byte(`{"type":"photos","email":"test@gmail.com"}`),
		))

		handler.Handle(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, 0, fakeScheduler.ScheduleRescanCallCount())
	})

	t.Run("request photos rescan schedule error", func(t *testing.T) {
		fakeScheduler := new(scannerfakes.FakeScheduler)
		fakeAccountRepository := new(accountfakes.FakeRepository)
		handler := NewRescanHandler(fakeAccountRepository, fakeScheduler)

		fakeAccountRepository.AccountExistReturns(true, nil)
		fakeScheduler.ScheduleRescanReturns(errors.New("schedule error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/rescan", bytes.NewBuffer(
			[]byte(`{"type":"photos","email":"test@gmail.com"}`),
		))

		handler.Handle(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, 1, fakeScheduler.ScheduleRescanCallCount())
	})
}
