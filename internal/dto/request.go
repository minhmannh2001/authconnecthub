package dto

import "time"

type UpdateUserProfile struct {
	FirstName   string     `json:"firstname"    form:"firstname"`
	LastName    string     `json:"lastname"     form:"lastname"`
	Gender      string     `json:"gender"       form:"gender"`
	Country     string     `json:"country"      form:"country"`
	City        string     `json:"city"         form:"city"`
	Address     string     `json:"address"      form:"address"`
	Email       string     `json:"email"        form:"email"         binding:"required,email"`
	PhoneNumber string     `json:"phone_number" form:"phone_number"`
	Birthday    *time.Time `json:"birthday"     form:"birthday"      time_format:"2006-01-02"`
	Company     string     `json:"company"      form:"company"`
	Role        string     `json:"role"         form:"role"`
}
