package entity

type Role struct {
	ID          int    `json:"id"`
	UserID      string `json:"user_id"`
	PrivelegeID int    `json:"privelege_id"`
}

type Privelege struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
