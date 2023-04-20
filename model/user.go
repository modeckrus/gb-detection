package model

type User struct {
	ID       uint   `json:"id" db:"id" gorm:"primaryKey`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}
