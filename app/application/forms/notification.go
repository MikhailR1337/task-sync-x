package forms

type Mailer struct {
	Email    string `json:"email"`
	Template string `json:"template"`
	Subject  string `json:"subject"`
}
