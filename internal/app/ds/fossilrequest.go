package ds

import "time"

type FossilRequest struct {
	Species        string    `json:"species"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
}
