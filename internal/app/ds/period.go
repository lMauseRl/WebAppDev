package ds

type Period struct {
	IDPeriod    uint   `gorm:"type:serial;primarykey" json:"id_period"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Age         string `json:"age"`
	Status      string `json:"status"`
	Photo       string `json:"photo"`
}
