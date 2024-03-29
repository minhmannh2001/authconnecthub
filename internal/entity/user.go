package entity

import "gorm.io/gorm"

// User model
type User struct {
	gorm.Model
	ID         uint   `gorm:"primary_key"                                   json:"id"`
	RoleID     uint   `gorm:"not null;DEFAULT:3"                            json:"role_id"`
	Username   string `gorm:"size:255;not null;unique"                      json:"username"`
	Email      string `gorm:"size:255;not null;unique"                      json:"email"`
	Password   string `gorm:"size:255;not null"                             json:"-"`
	Role       Role   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	RememberMe bool
}
