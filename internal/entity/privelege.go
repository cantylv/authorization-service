package entity

type GroupPrivelege struct {
	ID      int `json:"id"`
	GroupID int `json:"group_id"`
	AgentID int `json:"agent_id"`
}

type UserPrivelege struct {
	ID      int `json:"id"`
	UserID  int `json:"user_id"`
	AgentID int `json:"agent_id"`
}

type Agent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
