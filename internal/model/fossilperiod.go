package model

type Fossilperiod struct {
	FossilID uint `gorm:"type:serial;primaryKey;index" json:"fossil_id"`
	PeriodID uint `gorm:"type:serial;primaryKey;index" json:"period_id"`
}
