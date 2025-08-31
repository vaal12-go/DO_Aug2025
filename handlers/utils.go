package handlers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SendErrorResponse(c *gin.Context, httpStatus uint, message string, err error, otherdata any) {
	msg := message
	if err != nil {
		msg = fmt.Sprintf("%s Error:%s", msg, err)
	}
	if otherdata != nil {
		msg = fmt.Sprintf("%s Other data:%v", msg, otherdata)
	}
	c.JSON(200, gin.H{"status": "error", "message": msg})
}

func StringToUint64(str string) (uint64, error) {
	mission_id_int64, err := strconv.ParseInt(str, 10, 64)
	fmt.Println("id:", mission_id_int64)
	if err != nil {
		return 0, fmt.Errorf("Convertion of string '%s' to number failed with error:%s", str, err)
	}

	if mission_id_int64 < 0 {
		return 0, fmt.Errorf("String '%s' cannot be converted to non-negative integer", str)
	}

	return uint64(mission_id_int64), nil
}
