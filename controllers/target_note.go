package controllers

import (
	"do_aug25/db"
	"do_aug25/models"
	"fmt"
)

func CreateTargetNote(target_id uint64, targetNote models.TargetNote) (target *models.Target,
	statusCode uint, err error) {

	conn, err := db.SQLiteConnect()
	if err != nil {
		return nil, 500, fmt.Errorf("Error connecting to DB:%v", err)
	}

	target, err = TargetExists(target_id, conn)
	if err != nil {
		return nil, 400, fmt.Errorf("Target with id:%d not found. Error:%v", target_id, err)
	}

	if target.Complete {
		return nil, 400, fmt.Errorf("Target with id:%d Is complete. Notes cannot be created",
			target_id)
	}

	fmt.Println("target_note:28 target.MissionID::", target.MissionID)
	mission, err := GetMission(target.MissionID, conn)

	if err != nil {
		return nil, 400, fmt.Errorf("Error retrieving mission with id:%d. Error:%v",
			target.MissionID, err)
	}

	if mission.Complete {
		return nil, 400, fmt.Errorf("Mission with id:%d is complete. No notes can be created.",
			mission.ID)
	}

	targetNote.TargetID = uint64(target.ID)

	err = conn.Create(&targetNote).Error
	if err != nil {
		return nil, 500, fmt.Errorf("Error creating target note:%v", err)
	}

	target, err = GetTarget(target_id, conn)
	return target, 200, nil
}
