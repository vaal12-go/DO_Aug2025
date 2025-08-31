package models

import "time"

type Target struct {
	// gorm.Model
	ID        uint         `gorm:"primary_key" json:"id"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt *time.Time   `json:"-"`
	Name      string       `json:"name"`
	Country   string       `json:"country"`
	Complete  bool         `json:"complete" gorm:"type:bool DEFAULT False"`
	MissionID uint64       `json:"-"`
	Notes     []TargetNote `json:"notes"`
}
