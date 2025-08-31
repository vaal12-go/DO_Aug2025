package handlers

import (
	"do_aug25/controllers"
	"do_aug25/models"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CreateTargetNote(c *gin.Context) {
	target_id, err := StringToUint64(c.Param("target_id"))
	if err != nil {
		SendErrorResponse(c, 400, "Incorrect target id.", err, target_id)
		return
	}

	var newTargetNote models.TargetNote
	if err := c.ShouldBindJSON(&newTargetNote); err != nil {
		// fmt.Println("mission_handlers:57 err::", err)
		var ve validator.ValidationErrors
		validationMessages := make([]string, 0)
		if errors.As(err, &ve) {
			for _, fe := range ve {
				fmt.Println("cat_handlers:41 fe.Tag()::", fe.Tag())
				validationMessages = append(validationMessages, fe.Tag())
			}
		}
		errMsg := fmt.Sprintf("Invalid Mission JSON data:\n%s",
			strings.Join(validationMessages, "\n"))
		SendErrorResponse(c, 400, "Error parsing body for TargetNote", err, errMsg)
		return
	}

	target, statusCode, err := controllers.CreateTargetNote(target_id, newTargetNote)
	if err != nil {
		SendErrorResponse(c, statusCode, "Error creating new target note.", err, nil)
		return
	}

	c.JSON(200, gin.H{"message": "success", "object": target})
}
