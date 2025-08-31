package models

import "time"

type TargetNote struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Note      string     `json:"note" gorm:"type:text"`
	TargetID  uint64     `json:"-"`
}
