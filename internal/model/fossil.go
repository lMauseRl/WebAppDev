package model

import "time"

type Fossil struct {
	IDFossil       uint      `gorm:"type:serial;primarykey" json:"id_fossil"`
	Species        string    `json:"species"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
	Status         string    `json:"status"`
	UserID         uint      `json:"user_id"`
	ModeratorID    uint      `json:"moderator_id"`
}

type FossilGetResponse struct {
	IDFossil       uint      `json:"id_fossil"`
	Species        string    `json:"species"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
	Status         string    `json:"status"`
	Periods        []Period  `json:"periods"`
}

type FossilUpdateSpeciesRequest struct {
	Species string `json:"species"`
}

type FossilUpdateStatusRequest struct {
	Status string `json:"status"`
}
