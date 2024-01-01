package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID       uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string  `gorm:"not null" json:"username" valid:"required"`
	Email    string  `gorm:"unique;not null" json:"email" valid:"email,required"`
	Password string  `gorm:"not null;min:6" json:"password" valid:"required"`
	Photos   []Photo `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
