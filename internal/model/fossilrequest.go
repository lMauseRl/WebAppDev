package model

import "time"

type FossilRequest struct {
	IDFossil       uint      `gorm:"type:serial;primarykey" json:"id_fossil"`
	Species        string    `json:"species"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
	Status         string    `json:"status"`
	FullName       string    `json:"full_name"`
}
