package utilitymodels

import (
	"gorm.io/gorm"
	"time"
)

type Common struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type CommonSoftDelete struct {
	Common
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
