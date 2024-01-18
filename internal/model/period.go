package model

type Period struct {
	IDPeriod    uint   `gorm:"type:serial;primarykey" json:"id_period"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Age         string `json:"age"`
	Photo       string `json:"photo"`
}

type PeriodRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Age         string `json:"age"`
}

type PeriodGetResponse struct {
	Period   []Period `json:"periods"`
	IDFossil uint     `json:"id_fossil" example:"1"`
}
