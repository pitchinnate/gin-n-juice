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
)

func TestRegister(t *testing.T) {
	tester.TestPackage(t)
	defer tester.CleanupPackage(t)

	user := models.User{
		Email:    "test@test.com",
		Password: "testing",
		Admin:    false,
	}
	db.DB.Create(&user)

	t.Run("Test Register Empty Body", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{})
		w := tester.SetupTestRouter("POST", PostRegister, body)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		response := ErrorResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response.Error, "Field validation for 'Email' failed")
		assert.Contains(t, response.Error, "Field validation for 'Password' failed")
		assert.Contains(t, response.Error, "Field validation for 'ConfirmPassword' failed")
		assert.Contains(t, response.Error, "Field validation for 'VerifyUrl' failed")
	})
	t.Run("Test Register Bad Email", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": "test", "password": "testing1", "confirm_password": "testing1", "verify_url": "http://test.com"})
		w := tester.SetupTestRouter("POST", PostRegister, body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := ErrorResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response.Error, "Field validation for 'Email' failed")
	})
	t.Run("Test Register Short Password", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": "test@test.com", "password": "test", "confirm_password": "test", "verify_url": "http://test.com"})
		w := tester.SetupTestRouter("POST", PostRegister, body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := ErrorResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response.Error, "Field validation for 'Password' failed")
	})
	t.Run("Test Register Not Matching Password", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": "test@test.com", "password": "testing1", "confirm_password": "testing2", "verify_url": "http://test.com"})
		w := tester.SetupTestRouter("POST", PostRegister, body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := ErrorResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response.Error, "Field validation for 'ConfirmPassword' failed")
	})
	t.Run("Test Register Email already used", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": "test@test.com", "password": "testing1", "confirm_password": "testing1", "verify_url": "http://test.com"})
		w := tester.SetupTestRouter("POST", PostRegister, body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := ErrorResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response.Error, "Email address already used, use forgot password")
	})
	t.Run("Test Register Email valid", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": "foo@test.com", "password": "testing1", "confirm_password": "testing1", "verify_url": "http://test.com"})
		w := tester.SetupTestRouter("POST", PostRegister, body)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Check your email to verify your email address")
	})
}
