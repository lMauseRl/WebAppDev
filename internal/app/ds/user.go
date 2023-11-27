package ds

type User struct {
	UserID   uint   `gorm:"autoIncrement;primarykey" json:"id_user"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
