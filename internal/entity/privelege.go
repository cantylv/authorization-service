package entity

type Privelege struct {
	ID      int `json:"id"`
	GroupID int `json:"group_id"`
	AgentID int `json:"agent_id"`
}

type Agent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
