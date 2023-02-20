package types

type SignInInput struct {
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}
type SignUpInput struct {
	SignInInput
	FullName string `json:"fullName"`
}
