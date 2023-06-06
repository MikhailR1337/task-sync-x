package forms

type CreateHomeworkRequest struct {
	Name          string `json:"name" validate:"required,max=100"`
	Description   string `json:"description" validate:"required,max=500"`
	CurrentPoints string `json:"currentPoints" validate:"required,max=2"`
	MaxPoints     string `json:"maxPoints" validate:"required,max=2"`
	Type          string `json:"type" validate:"required,oneof=listening reading"`
	Status        string `json:"status" validate:"required,oneof=new processing finished checked"`
	Student       string `json:"student" validate:"required"`
}

type UpdateHomeworkStudentRequest struct {
	Status string `json:"status" validate:"required,oneof=processing finished"`
}

type UpdateHomeworkTeacherRequest struct {
	Status        string `json:"status" validate:"required,oneof=checked"`
	CurrentPoints string `json:"currentPoints" validate:"required,max=2"`
}
