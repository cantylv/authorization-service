package entity

type Group struct {
	ID      int    `json:"id"`
	OwnerID string `json:"owner_id"`
}

type Participation struct {
	ID      int    `json:"id"`
	UserID  string `json:"user_id"`
	GroupID int    `json:"group_id"`
}
