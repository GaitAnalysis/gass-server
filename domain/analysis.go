package domain

import (
	"gorm.io/gorm"
)

type AnalysisInput struct {
	UserId   uint   `json:"user_id"`
	UploadId uint   `json:"upload_id"`
	Result   Result `json:"result"`
}

type Result struct {
	Confidence float64 `gorm:"not null" json:"confidence"`
	Label      string  `gorm:"size:255;not null" json:"label"`
}

type Analysis struct {
	gorm.Model
	Upload   Upload `gorm:"foreignKey:UploadId" json:"upload"`
	User     User   `gorm:"foreignKey:UserId" json:"-"`
	UploadId uint   `gorm:"not null" json:"-"`
	UserId   uint   `gorm:"not null" json:"-"`
	Result   Result `gorm:"embedded;embeddedPrefix:result_" json:"result"`
}
