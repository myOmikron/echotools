package utilitymodels

import (
	"database/sql"
)

type User struct {
	Common
	LastLoginAt sql.NullTime `json:"-" gorm:"default:null"` // This is only relevant if the session middleware is in use
	Email       string       `json:"email" gorm:"unique;default:null"`
	Username    string       `json:"username" gorm:"unique;not null"`
	Password    string       `json:"-" gorm:"not null"`
	Active      sql.NullBool `json:"active" gorm:"default:true"`
}
