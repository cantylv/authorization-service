package dto

type SignInData struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"runelength(6|30)"`
}

type SignInForm struct {
	Payload  SignInData `json:"payload"`
	JwtToken string     `json:"jwt_token"`
}

type SignUpData struct {
	Email     string `json:"email" valid:"email"`
	Password  string `json:"password" valid:"runelength(6|30)"`
	FirstName string `json:"first_name" valid:"runelength(2|50)"`
	LastName  string `json:"last_name" valid:"runelength(2|50)"`
}
