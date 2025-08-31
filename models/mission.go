package models

import "time"

type Mission struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name"`
	Complete  bool       `json:"complete" gorm:"type:bool;default:false"`
	CatID     uint64     `json:"-"`
	Targets   []Target   `json:"targets"`
}
