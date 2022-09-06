package tester

import (
	"bytes"
	"fmt"
	"gin-n-juice/config"
	"gin-n-juice/db"
	"gin-n-juice/utils/file"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var packageDb string

// TestPackage creates a copy of the empty test database so we can use it then delete it without having to
// rerun migrations for each test
func TestPackage(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	exPath := filepath.Dir(filename)

	packageDb = fmt.Sprintf("%s/../../tmp/test_%s.db", exPath, RandStringBytesRmndr(10))
	file.CopyFile(fmt.Sprintf("%s/../../tmp/test.db", exPath), packageDb)
	config.SetupTestEnv(packageDb)
	db.ConnectDatabase(logger.Silent)
}

// CleanupPackage deletes the temp database we are using for this package
func CleanupPackage(t *testing.T) {
	db.Disconnect()
	t.Log("Removing db file created for this package: ", packageDb)
	if err := os.Remove(packageDb); err != nil {
		t.Log("Error removing file: ", err)
	}
}

func SetupTestRouter(method string, handler func(c *gin.Context), body []byte) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	HttpMethods := map[string]interface{}{
		"POST":    router.POST,
		"GET":     router.GET,
		"PUT":     router.PUT,
		"DELETE":  router.DELETE,
		"PATCH":   router.PATCH,
		"HEAD":    router.HEAD,
		"OPTIONS": router.OPTIONS,
	}

	f := reflect.ValueOf(HttpMethods[method])
	f.Call([]reflect.Value{
		reflect.ValueOf("/"),
		reflect.ValueOf(handler),
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, "/", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	return w
}

func RandStringBytesRmndr(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
