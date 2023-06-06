package forms

type RegistrateRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=30"`
	Role     string `json:"role" validate:"required,oneof=student teacher"`
}
