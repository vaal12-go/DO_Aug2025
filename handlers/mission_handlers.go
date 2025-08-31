package handlers

import (
	"do_aug25/controllers"
	"do_aug25/db"
	"do_aug25/models"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func AssignMissionToCat(c *gin.Context) {
	cat_id, err := StringToUint64(c.Param("cat_id"))
	if err != nil {
		SendErrorResponse(c, 400, "Incorrect cat id.", err, cat_id)
		return
	}

	mission_id, err := StringToUint64(c.Param("mission_id"))
	if err != nil {
		SendErrorResponse(c, 400, "Incorrect mission id.", err, mission_id)
		return
	}

	cat, statusCode, err := controllers.AssignMissionToCat(cat_id, mission_id)
	if err != nil {
		SendErrorResponse(c, statusCode, "Error assigning mission to cat.", err, nil)
		return
	}

	c.JSON(200, gin.H{"status": "success", "message": cat})
}

func GetMission(c *gin.Context) {
	mission_id_int64, err := strconv.ParseInt(c.Param("id"), 10, 64)
	fmt.Println("mission_id:", mission_id_int64)
	if err != nil {
		SendErrorResponse(c, 400, "Mission ID is not integer", err, c.Param("id"))
		return
	}

	if mission_id_int64 <= 0 {
		SendErrorResponse(c, 400, "MissionID is invalid", nil, mission_id_int64)
		return
	}

	mission_id := uint64(mission_id_int64)

	conn, err := db.SQLiteConnect()
	if err != nil {
		SendErrorResponse(c, 500, "Error connecting to DB", err, nil)
		return
	}

	mission, err := controllers.GetMission(mission_id, conn)
	if err != nil {
		SendErrorResponse(c, 400, "No mission with this id", err, mission_id)
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": mission})
}

func GetMissions(c *gin.Context) {
	var missions []models.Mission
	conn, err := db.SQLiteConnect()
	if err != nil {
		c.JSON(400, gin.H{"error": "Error connecting to DB"})
		return
	}
	err = conn.Model(&models.Mission{}).Preload("Targets").Find(&missions).Error
	if err != nil {
		fmt.Println("mission_handlers:19 err::", err)
		c.JSON(400, gin.H{"error": "Error retrieving missions"})
		return
	}
	c.JSON(200, gin.H{"message": "success", "object": missions})
}

func CreateMission(c *gin.Context) {
	// fmt.Println("mission_handlers:6 qwe::")
	// id := c.Param("catid")
	// fmt.Println(id)
	var newMission models.Mission
	if err := c.ShouldBindJSON(&newMission); err != nil {
		fmt.Println("mission_handlers:57 err::", err)
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
		c.JSON(400, gin.H{"error": errMsg})
		return
	}

	fmt.Println("mission_handlers:72 newMission::", newMission)

	conn, err := db.SQLiteConnect()
	if err != nil {
		c.JSON(400, gin.H{"error": "Error connecting to DB"})
		return
	}
	conn.Create(&newMission)
	c.JSON(200, gin.H{"message": "success", "object": newMission})
	// c.JSON(200, gin.H{"message": "success - mission CREATE handler works 22"})
}

func DeleteMission(c *gin.Context) {
	// id := c.Param("id")
	mission_id_int64, err := strconv.ParseInt(c.Param("id"), 10, 64)
	fmt.Println("mission_id:", mission_id_int64)
	if err != nil {
		c.JSON(400, gin.H{"error": "mission ID is not integer"})
		return
	}

	if mission_id_int64 <= 0 {
		c.JSON(400, gin.H{"error": "MissionID is invalid"})
		return
	}

	mission_id := uint64(mission_id_int64)

	conn, err := db.SQLiteConnect()
	if err != nil {
		c.JSON(400, gin.H{"error": "Error connecting to DB"})
		return
	}
	// var mission2Delete models.Mission

	mission2Delete, err := controllers.MissionExists(mission_id, conn)
	// res := conn.Model(&models.Mission{}).Preload("Targets").First(&mission2Delete, id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Error looking for mission"})
		return
	}

	if mission2Delete == nil {
		c.JSON(400, gin.H{"error": "No mission with this ID"})
		return
	}

	if mission2Delete.CatID > 0 {
		c.JSON(400, gin.H{"error": "Cannot delete mission, which is assigned to cat"})
		return
	}

	if len(mission2Delete.Targets) > 0 {
		for _, trgt := range mission2Delete.Targets {
			fmt.Println("mission_handlers:108 DELETING trgt.ID::", trgt.ID)
			err := controllers.DeleteTarget(trgt, conn)
			if err != nil {
				c.JSON(400, gin.H{"error": "Cannot delete target object. Aborting"})
				return
			}
		}
	}

	res := conn.Delete(&mission2Delete)
	fmt.Println("mission_handlers:111 res::", res)
	if res.Error != nil {
		c.JSON(400, gin.H{"error": "Error deleting mission"})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "Mission is deleted."})
}

func MarkMissionCompleted(c *gin.Context) {
	mission_id, err := StringToUint64(c.Param("id"))
	if err != nil {
		SendErrorResponse(c, 400, "", err, nil)
		return
	}

	conn, err := db.SQLiteConnect()
	if err != nil {
		SendErrorResponse(c, 500, "Error connecting to DB", err, nil)
		return
	}
	mission2Update, err := controllers.MissionExists(mission_id, conn)
	if err != nil {
		SendErrorResponse(c, 400, "No mission with this ID", err, mission_id)
		return
	}

	if mission2Update.Complete {
		SendErrorResponse(c, 400, "Mission is already marked as complete", err, mission_id)
		return
	}

	mission2Update.Complete = true
	res := conn.Save(&mission2Update)
	if res.Error != nil {
		SendErrorResponse(c, 500, "Error updating mission.", res.Error, mission_id)
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": mission2Update})
}
