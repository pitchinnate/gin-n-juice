package users

import (
	"gin-n-juice/utils/tester"
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

	t.Run("Get Users", func(t *testing.T) {
		w := tester.SetupTestRouter("GET", GetList, nil)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
