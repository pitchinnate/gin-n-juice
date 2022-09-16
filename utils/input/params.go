package input

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetPathInt(c *gin.Context, name string) (int, error) {
	val := c.Params.ByName(name)
	if val == "" {
		return 0, errors.New(name + " path parameter value is empty or not specified")
	}
	return strconv.Atoi(val)
}
