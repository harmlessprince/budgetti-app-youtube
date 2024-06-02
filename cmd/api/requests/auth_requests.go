package requests

type RegisterUserRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=2"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=2"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=2"`
	Password        string `json:"password" validate:"required,min=2"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
}
