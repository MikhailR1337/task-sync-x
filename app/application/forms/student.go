package forms

type StudentUpdateRequest struct {
	Teacher string `json:"teacher" validate:"required"`
}
