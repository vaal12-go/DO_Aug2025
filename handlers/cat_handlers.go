package handlers

import (
	"do_aug25/db"
	"do_aug25/models"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SalaryUpdate struct {
	Salary uint32
}

func UpdateCatsSalary(c *gin.Context) {
	id := c.Param("id")
	// fmt.Println(id)

	var sUpdate SalaryUpdate
	if err := c.ShouldBindJSON(&sUpdate); err != nil {
		c.JSON(400, gin.H{"status": "error", "message": "New salary is not (or incorrectly) provided."})
		return
	}

	var cat2Update models.Cat
	conn, err := db.SQLiteConnect()
	if err != nil {
		c.JSON(400, gin.H{"error": "Error connecting to DB"})
		return
	}
	res := conn.First(&cat2Update, id)
	if res.Error != nil {
		c.JSON(400, gin.H{"error": "No cat with this ID"})
		return
	}

	cat2Update.Salary = sUpdate.Salary
	res = conn.Save(&cat2Update)
	if res.Error != nil {
		c.JSON(400, gin.H{"error": "Error updating cat's salary"})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": cat2Update})
}

func DeleteCat(c *gin.Context) {
	id := c.Param("id")
	// fmt.Println(id)

	var cat2Delete models.Cat
	conn, err := db.SQLiteConnect()
	if err != nil {
		c.JSON(400, gin.H{"error": "Error connecting to DB"})
		return
	}
	res := conn.First(&cat2Delete, id)
	if res.Error != nil {
		c.JSON(400, gin.H{"error": "No cat with this ID"})
		return
	}

	res = conn.Delete(&cat2Delete)
	if res.Error != nil {
		c.JSON(400, gin.H{"error": "Error deleting cat"})
		return
	}

	c.JSON(200, gin.H{"status": "success", "message": "Cat is deleted."})
}

func GetCat(c *gin.Context) {
	fmt.Println("Get cat called")
	id := c.Param("id")
	fmt.Println(id)

	conn, err := db.SQLiteConnect()
	if err != nil {
		c.JSON(400, gin.H{"error": "Error connecting to DB"})
		return
	}

	var cat models.Cat
	res := conn.Model(&models.Cat{}).Preload("Mission").First(&cat, id)
	if res.Error != nil {
		c.JSON(400, gin.H{"error": "No cat with this id"})
		return
	}
	// // Fetch user from database
	c.JSON(200, gin.H{"status": "success", "message": cat})
}

func GetCats(c *gin.Context) {
	conn, err := db.SQLiteConnect()
	if err != nil {
		c.JSON(400, gin.H{"error": "Error connecting to DB"})
		return
	}
	// cats := make([]models.Cat, 0)
	var cats []models.Cat
	err = conn.Model(&models.Cat{}).Preload("Mission").Find(&cats).Error
	if err != nil {
		c.JSON(400, gin.H{"error": "Error retrieving cats"})
		return
	}
	c.JSON(200, gin.H{"message": "success", "object": cats})
}

func CreateCat(c *gin.Context) {
	var newCat models.Cat
	if err := c.ShouldBindJSON(&newCat); err != nil {
		fmt.Println("cat_handlers:114 err::", err)
		var ve validator.ValidationErrors
		validationMessages := make([]string, 0)
		if errors.As(err, &ve) {
			for _, fe := range ve {
				fmt.Println("cat_handlers:41 fe.Tag()::", fe.Tag())
				validationMessages = append(validationMessages, models.MessageForValidationTag(
					fe.Tag()))
			}
			// c.JSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		errMsg := fmt.Sprintf("Invalid JSON data:\n%s", strings.Join(validationMessages, "\n"))
		c.JSON(400, gin.H{"error": errMsg})
		return
	}
	// fmt.Println("cat_handlers:24 newCat::", newCat)
	// fmt.Println("cat_handlers:25 newCat.YearsOfExperience::", newCat.YearsOfExperience)
	// fmt.Println("cat_handlers:27 newCat.Mission::", newCat.Mission)

	newCat.Mission = nil

	conn, err := db.SQLiteConnect()
	if err != nil {
		c.JSON(400, gin.H{"error": "Error connecting to DB"})
		return
	}
	conn.Create(&newCat)
	c.JSON(200, gin.H{"message": "success", "object": newCat})
}
