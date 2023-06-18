package domain

import (
	"gorm.io/gorm"
)

type Upload struct {
	gorm.Model
	User  User   `gorm:"foreignKey:UserId" json:"-"`
	FileRef string `gorm:"size:255;not null;" json:"file_ref"`
	Size   int64  `gorm:"not null;" json:"size"`
	UserId uint   `gorm:"not null;" json:"-"`
}
