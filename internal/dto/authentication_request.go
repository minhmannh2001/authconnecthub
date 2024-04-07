package dto

// LoginRequestBody
type LoginRequestBody struct {
	Username   string `json:"username"    form:"username" binding:"required"`
	Password   string `json:"password"    form:"password" binding:"required,min=8"`
	RememberMe string `json:"remember_me" form:"remember_me"`
}

// RegisterRequestBody
type RegisterRequestBody struct {
	Username        string `json:"username"         form:"username"         binding:"required"`
	Email           string `json:"email"            form:"email"            binding:"required,email"`
	Password        string `json:"password"         form:"password"         binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required,min=8,eqfield=Password"`
}
