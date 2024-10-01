package dto

type Bid struct {
	ID        int    `json:"id"`
	GroupName string `json:"group_name"`
	UserId    string    `json:"user_id"`
	Status    string `json:"status"`
}
