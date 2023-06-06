package forms

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=30"`
	Role     string `json:"role" validate:"required,oneof=student teacher"`
}
