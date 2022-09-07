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

func TestLogin(t *testing.T) {
	tester.TestPackage(t)
	defer tester.CleanupPackage(t)

	password, _ := models.HashPassword("testing")
	user := models.User{
		Email:    "test@test.com",
		Password: password,
		Admin:    false,
	}
	db.DB.Create(&user)

	t.Run("Test Login Empty Body", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{})
		w := tester.SetupTestRouter("POST", PostLogin, body)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		response := ErrorResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response.Error, "Field validation for 'Email' failed")
		assert.Contains(t, response.Error, "Field validation for 'Password' failed")
	})
	t.Run("Test Login Bad Params", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": "", "password": ""})
		w := tester.SetupTestRouter("POST", PostLogin, body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := ErrorResponse{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response.Error, "Field validation for 'Email' failed")
		assert.Contains(t, response.Error, "Field validation for 'Password' failed")
	})
	t.Run("Test Login Email Not Found", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": "bad@test.com", "password": "testing"})
		w := tester.SetupTestRouter("POST", PostLogin, body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid email or password")
	})
	t.Run("Test Login Invalid Password", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": user.Email, "password": "badpass"})
		w := tester.SetupTestRouter("POST", PostLogin, body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid email or password")
	})
	t.Run("Test Login Valid Password", func(t *testing.T) {
		body, _ := json.Marshal(gin.H{"email": user.Email, "password": "testing"})
		w := tester.SetupTestRouter("POST", PostLogin, body)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "test@test.com")
	})
}
