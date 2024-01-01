package models

import (
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"photoId"`
	Title          string `gorm:"not null" form:"title" json:"title"`
	Caption        string `form:"caption" json:"caption"`
	PhotoUrl       string `gorm:"not null" form:"photoUrl" json:"photoUrl" valid:"required"`
	UserID         uint   `form:"uId" json:"uId"`
	IsProfilePhoto bool   `gorm:"not null"`
	User           User   `gorm:"foreignKey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
