package ds

import "time"

type Periods struct {
	Id          uint `gorm:"primarykey"`
	Name        string
	Description string
	Age_start   float32
	Age_end     float32
	Photo       string
	Status      string
}

type FossilPeriod struct {
	PeriodsID uint `gorm:"primarykey"`
	FossilID  uint `gorm:"primarykey"`
}

type Fossil struct {
	RequestID     uint   `gorm:"primarykey"`
	Status        string `gorm:"size:30"`
	StartDate     time.Time
	FormationDate time.Time
	EndDate       time.Time
	UserID        uint
	ModeratorID   uint
	Genus         string
	Species       string
}

type Users struct {
	UserID   uint   `gorm:"primarykey"`
	Name     string `gorm:"size:60"`
	Email    string `gorm:"unique;size:60"`
	Role     string `gorm:"size:60"`
	Password string `gorm:"size:60"`
}
