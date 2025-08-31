package controllers

import (
	"do_aug25/db"
	"do_aug25/models"
	"fmt"

	"gorm.io/gorm"
)

func AssignMissionToCat(cat_id uint64, mission_id uint64) (cat *models.Cat,
	statusCode uint, err error) {
	conn, err := db.SQLiteConnect()
	if err != nil {
		return nil, 500, err
	}

	cat, err = CatExists(cat_id, conn)
	if err != nil {
		return nil, 400, fmt.Errorf("Cat with id:%d does not exist. Error:%v", cat_id, err)
	}

	cat, err = GetCat(cat_id, conn)
	if err != nil {
		return nil, 500, fmt.Errorf("Error retrieving cat with id=%d. Error:%v",
			cat_id, err)
	}

	fmt.Println("mission:33 cat.Mission::", cat.Mission)
	if cat.Mission != nil && !cat.Mission.Complete {
		return nil, 500, fmt.Errorf("Cat with id=%d already has a mission [id=%d] and it is not complete.",
			cat_id, cat.Mission.ID)
	}

	mission, err := MissionExists(mission_id, conn)
	if err != nil {
		return nil, 400, fmt.Errorf("Mission with id:%d does not exist. Error:%v", mission_id, err)
	}

	if mission.CatID > 0 {
		return nil, 400, fmt.Errorf("Mission with id:%d already assigned to cat %d. ",
			mission_id, mission.CatID)
	}

	mission.CatID = uint64(cat.ID)

	err = conn.Save(&mission).Error
	if err != nil {
		return nil, 500, fmt.Errorf("Error assigning mission id:%d to cat %d. Error:%v",
			mission_id, mission.CatID, err)
	}

	cat, err = GetCat(cat_id, conn)
	if err != nil {
		return nil, 500, fmt.Errorf("Error retrieving cat with id=%d. Error:%v",
			mission.CatID, err)
	}

	return cat, 200, nil

}

func MissionExists(id uint64, conn *gorm.DB) (*models.Mission, error) {
	var mission models.Mission
	res := conn.First(&mission, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &mission, nil
}

func GetMission(id uint64, conn *gorm.DB) (*models.Mission, error) {
	var mission models.Mission
	err := conn.Model(&models.Mission{}).
		Preload("Targets").First(&mission, id).Error
	if err != nil {
		return nil, err
	}
	return &mission, nil
}
