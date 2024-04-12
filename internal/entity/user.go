package entity

import (
	"time"

	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	ID          uint        `gorm:"primary_key"                                                          json:"id"`
	RoleID      uint        `gorm:"not null;DEFAULT:3"                                                   json:"role_id"`
	Username    string      `gorm:"size:255;not null;unique"                                             json:"username"`
	Email       string      `gorm:"size:255;not null;unique"                                             json:"email"`
	Password    string      `gorm:"size:255;not null"                                                    json:"-"`
	Role        Role        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"                        json:"-"`
	RememberMe  bool        `json:"-"`
	UserProfile UserProfile `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_profile"`
}

// User Profile model
type UserProfile struct {
	gorm.Model
	UserID      uint       `gorm:"primary_key" ` // Foreign key to User.ID
	FirstName   string     `gorm:"size:30"  json:"first_name"`
	LastName    string     `gorm:"size:30"  json:"last_name"`
	Country     string     `gorm:"size:50"  json:"country"`
	City        string     `gorm:"size:50"  json:"city"`
	PhoneNumber string     `gorm:"size:10"  json:"phone_number"`
	Birthday    *time.Time `json:"birthday,omitempty"`
	Company     string     `gorm:"size:255" json:"company"`
	Role        string     `gorm:"size:30"  json:"role"`
	Gender      string     `gorm:"size:15"  json:"gender"`
	Address     string     `gorm:"size:255" json:"address"`
}

type SocialAccount struct {
	gorm.Model
	UserID      uint   `gorm:"primary_key" json:"user_id"`      // Foreign key to User.ID
	AccountType string `gorm:"size:50"     json:"account_type"` // Type of social account (e.g., "facebook", "twitter")
	AccountLink string `gorm:"size:255"    json:"account_link"` // Link to the user's social account profile
	ButtonState string `json:"-"`
}
