package ds

import "time"

type Fossil struct {
	FossilID       uint      `gorm:"type:serial;primarykey" json:"id_fossil"`
	Species        string    `json:"species"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
	FossilStatus   string    `json:"status"`
	UserID         uint      `json:"user_id"`
	ModeratorID    uint      `json:"moderator_id"`
}

type FossilRequest struct {
	FossilID       uint      `gorm:"type:serial;primarykey" json:"id_fossil"`
	Genus          string    `json:"genus"`
	Species        string    `json:"species"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
	FossilStatus   string    `json:"status"`
	FullName       string    `json:"full_name"`
}
