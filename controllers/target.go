package controllers

import (
	"do_aug25/db"
	"do_aug25/models"
	"fmt"

	"gorm.io/gorm"
)

func DeleteTargetFromMission(mission_id uint64, target_id uint64) (mission *models.Mission,
	statusCode uint, err error) {
	conn, err := db.SQLiteConnect()
	if err != nil {
		return nil, 500, err
	}

	mission, err = MissionExists(mission_id, conn)
	if err != nil {
		return nil, 400, fmt.Errorf("Mission with id:%d does not exist. Error:%v", mission_id, err)
	}

	target, err := TargetExists(target_id, conn)
	if err != nil {
		return nil, 400, fmt.Errorf("Target with id:%d does not exists. Error:%v", target_id, err)
	}

	if target.MissionID != uint64(mission.ID) {
		return nil, 400, fmt.Errorf("Target id:%d does not belong to mission with id:%d", target_id, mission_id)
	}

	if target.Complete {
		return nil, 400, fmt.Errorf("Target id:%d is complete. Cannot delete", target_id)
	}

	// TODO: add notes deletion as well
	err = conn.Delete(&target).Error
	if err != nil {
		return nil, 500, fmt.Errorf("Error deleting target %d. Error:%v", target_id, err)
	}

	mission, err = GetMission(mission_id, conn)
	if err != nil {
		return nil, 500, fmt.Errorf("Error retrieving mission id:%d. Error:%v", mission_id, err)
	}

	return mission, 200, nil
}

func MarkTargetAsComplete(id uint64) (target *models.Target, statusCode uint, err error) {
	conn, err := db.SQLiteConnect()
	if err != nil {
		return nil, 500, err
	}

	target, err = TargetExists(id, conn)
	if err != nil {
		return nil, 400, fmt.Errorf("Target with id %d does not exists. Error:%v", id, err)
	}

	if target.Complete {
		return nil, 400, fmt.Errorf("Target %d is already complete", id)
	}

	target.Complete = true
	err = conn.Save(&target).Error
	if err != nil {
		return nil, 500, fmt.Errorf("Error updating target's [id=%d] Complete status:%v", id, err)
	}

	target, err = GetTarget(id, conn)
	if err != nil {
		return nil, 500, fmt.Errorf("Error retrieving target with id=%d. Error:%v", id, err)
	}
	return target, 200, nil
}

func TargetExists(id uint64, conn *gorm.DB) (*models.Target, error) {
	var target models.Target
	res := conn.First(&target, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &target, nil
}

func GetTarget(id uint64, conn *gorm.DB) (*models.Target, error) {
	var target models.Target
	res := conn.Model(&models.Target{}).Preload("Notes").First(&target, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &target, nil
}

func DeleteTarget(trgt models.Target, conn *gorm.DB) error {
	// TODO: Add deletion of connected notes

	res := conn.Delete(&trgt)
	fmt.Println("mission_handlers:136 res::", res)
	return res.Error
}
