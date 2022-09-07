package auth

import (
	"encoding/json"
	"gin-n-juice/db"
	"gin-n-juice/models"
	"gin-n-juice/utils/tester"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestReset(t *testing.T) {
	tester.TestPackage(t)
	defer tester.CleanupPackage(t)

	password, _ := models.HashPassword("testing")
	user := models.User{
		Email:    "test@test.com",
		Password: password,
		Admin:    false,
	}
	db.DB.Create(&user)

	forgot := models.PasswordReset{
		Email:     user.Email,
		Token:     "testingToken",
		CreatedAt: time.Now(),
	}
	db.DB.Create(&forgot)

	t.Run("Test Reset Empty Body", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{})
		w := tester.SetupTestRouter("POST", PostReset, body)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		response := ErrorResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response.Error, "Field validation for 'Email' failed")
		assert.Contains(t, response.Error, "Field validation for 'Password' failed")
		assert.Contains(t, response.Error, "Field validation for 'ConfirmPassword' failed")
		assert.Contains(t, response.Error, "Field validation for 'Token' failed")
	})
	t.Run("Test Reset Bad Email", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": "test", "password": "testing1", "confirm_password": "testing1", "token": "test"})
		w := tester.SetupTestRouter("POST", PostReset, body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := ErrorResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response.Error, "Field validation for 'Email' failed")
	})
	t.Run("Test Reset Short Password", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": "test@test.com", "password": "test", "confirm_password": "test", "token": "test"})
		w := tester.SetupTestRouter("POST", PostReset, body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := ErrorResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response.Error, "Field validation for 'Password' failed")
	})
	t.Run("Test Reset Not Matching Password", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": "test@test.com", "password": "testing1", "confirm_password": "testing2", "token": "test"})
		w := tester.SetupTestRouter("POST", PostReset, body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := ErrorResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response.Error, "Field validation for 'ConfirmPassword' failed")
	})
	t.Run("Test Reset Invalid Token", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": user.Email, "password": "testing1", "confirm_password": "testing1", "token": "badToken"})
		w := tester.SetupTestRouter("POST", PostReset, body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := ErrorResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response.Error, "invalid token")
	})
	t.Run("Test Reset Email valid", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": user.Email, "password": "testing1", "confirm_password": "testing1", "token": forgot.Token})
		w := tester.SetupTestRouter("POST", PostReset, body)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "password updated")
	})
}
