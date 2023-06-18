package domain

import (
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Email string `gorm:"size:255;not null;unique" json:"email"`
    Password string `gorm:"size:255;not null;" json:"-"`
    Uploads []Upload `gorm:"foreignKey:UserId" json:"uploads"`
    Analyses []Analysis `gorm:"foreignKey:UserId" json:"analyses"`
}
