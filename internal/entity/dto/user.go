package dto

// INPUT DATAFLOW
type CreateData struct {
	Email     string `json:"email" valid:"email"`
	Password  string `json:"password" valid:"runelength(6|30)"`
	FirstName string `json:"first_name" valid:"runelength(2|50)"`
	LastName  string `json:"last_name" valid:"runelength(2|50)"`
}

// OUTPUT DATAFLOW
type UserWithoutPassword struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
