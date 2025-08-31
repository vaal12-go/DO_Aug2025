package controllers

import (
	"do_aug25/models"

	"gorm.io/gorm"
)

func CatExists(cat_id uint64, conn *gorm.DB) (*models.Cat, error) {
	var cat models.Cat
	res := conn.First(&cat, cat_id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &cat, nil
}

func GetCat(cat_id uint64, conn *gorm.DB) (*models.Cat, error) {
	var cat models.Cat
	res := conn.Model(&models.Cat{}).Preload("Mission").First(&cat, cat_id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &cat, nil
}
