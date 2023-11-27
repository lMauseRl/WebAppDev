package ds

type Period struct {
	PeriodID     uint   `gorm:"type:serial;primarykey" json:"id_period"`
	PeriodName   string `json:"name"`
	PeriodDesc   string `json:"description"`
	PeriodAge    string `json:"age"`
	PeriodStatus string `json:"status"`
	PhotoURL     string `json:"photo"`
}
