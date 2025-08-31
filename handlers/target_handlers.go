package handlers

import (
	"do_aug25/controllers"
	"do_aug25/db"
	"do_aug25/models"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func DeleteTarget(c *gin.Context) {
	mission_id, err := StringToUint64(c.Param("mission_id"))
	if err != nil {
		SendErrorResponse(c, 400, "Incorrect mission id", err, nil)
		return
	}

	target_id, err := StringToUint64(c.Param("target_id"))
	if err != nil {
		SendErrorResponse(c, 400, "Incorrect target id", err, nil)
		return
	}

	mission, statusCode, err := controllers.DeleteTargetFromMission(mission_id, target_id)
	if err != nil {
		SendErrorResponse(c, statusCode, "Error deleting target from mission.", err, nil)
		return
	}
	c.JSON(200, gin.H{"message": "success", "object": mission})
}

func MarkTargetCompleted(c *gin.Context) {
	target_id, err := StringToUint64(c.Param("id"))
	if err != nil {
		SendErrorResponse(c, 400, "Incorrect target id", err, nil)
		return
	}

	target, statusCode, err := controllers.MarkTargetAsComplete(target_id)
	fmt.Println("target_handlers:23 err::", err)
	fmt.Println("target_handlers:24 statusCode::", statusCode)
	if err != nil {
		SendErrorResponse(c, statusCode, "Error marking target complete.",
			err, target_id)
		return
	}

	c.JSON(200, gin.H{"message": "success", "object": target})
}

func CreateTarget(c *gin.Context) {
	mission_id_int64, err := StringToUint64(c.Param("mission_id"))
	if err != nil {
		SendErrorResponse(c, 400, "Incorrect mission id", err, nil)
		return
	}
	mission_id := uint64(mission_id_int64)

	var newTarget models.Target
	if err := c.ShouldBindJSON(&newTarget); err != nil {
		var ve validator.ValidationErrors
		validationMessages := make([]string, 0)
		if errors.As(err, &ve) {
			for _, fe := range ve {
				fmt.Println("cat_handlers:41 fe.Tag()::", fe.Tag())
				validationMessages = append(validationMessages, models.MessageForValidationTag(
					fe.Tag()))
			}
		}
		errMsg := fmt.Sprintf("Invalid JSON data:\n%s", strings.Join(validationMessages, "\n"))
		SendErrorResponse(c, 400, errMsg, err, nil)
		// c.JSON(400, gin.H{"error": errMsg})
		return
	}
	// fmt.Println("target_handlers:32 newTarget::", newTarget)

	conn, err := db.SQLiteConnect()
	if err != nil {
		SendErrorResponse(c, 500, "Error connecting to DB", err, nil)
		// c.JSON(400, gin.H{"error": "Error connecting to DB"})
		return
	}

	mission2Update, err := controllers.MissionExists(mission_id, conn)
	if err != nil {
		SendErrorResponse(c, 500, "Error querying DB for mission", err, nil)
		// c.JSON(400, gin.H{"error": "Error querying DB for mission"})
		return
	}
	if mission2Update == nil {
		SendErrorResponse(c, 400, "No such mission exists.", nil, mission_id)
		// c.JSON(400, gin.H{"error": "No such mission exists."})
		return
	}

	if mission2Update.Complete {
		SendErrorResponse(c, 400, "Mission is completed. New targets cannot be added.", nil, mission_id)
		return
	}

	newTarget.MissionID = uint64(mission2Update.ID)
	conn.Create(&newTarget)

	mission2Update, err = controllers.GetMission(mission_id, conn)
	if err != nil {
		SendErrorResponse(c, 500, "Error querying DB for mission", err, nil)
		// c.JSON(400, gin.H{"error": "Error querying DB for mission"})
		return
	}
	c.JSON(200, gin.H{"message": "success", "object": mission2Update})
}
