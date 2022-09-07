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

type ErrorResponse struct {
	Error string `json:"error"`
}

func TestForgot(t *testing.T) {
	tester.TestPackage(t)
	defer tester.CleanupPackage(t)

	t.Run("Test Forgot Password with Body Empty", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{})
		w := tester.SetupTestRouter("POST", PostForgot, body)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		response := ErrorResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response.Error, "Field validation for 'Email' failed")
		assert.Contains(t, response.Error, "Field validation for 'ReturnUrl' failed")
	})
	t.Run("Test Forgot Password with Bad Params", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": "testing", "return_url": "test"})
		w := tester.SetupTestRouter("POST", PostForgot, body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := ErrorResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response.Error, "Field validation for 'Email' failed")
		assert.Contains(t, response.Error, "Field validation for 'ReturnUrl' failed")
	})
	t.Run("Test Forgot Password with Email Not Found", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": "test@test.com", "return_url": "http://test.com"})
		w := tester.SetupTestRouter("POST", PostForgot, body)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "email sent")
	})
	t.Run("Test Forgot Password with Email Found", func(t *testing.T) {
		user := models.User{
			Email:    "test@test.com",
			Password: "test",
			Admin:    false,
		}
		response := db.DB.Create(&user)
		if response.Error != nil {
			t.Error("Error creating user:", response.Error)
		}

		body, _ := json.Marshal(gin.H{"email": "test@test.com", "return_url": "http://test.com"})
		w := tester.SetupTestRouter("POST", PostForgot, body)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "email sent")
	})
}
