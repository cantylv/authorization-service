package entity

type Group struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}
