package dto

import ent "github.com/cantylv/authorization-service/internal/entity"

// OUTPUT DATAFLOW
type UserResponse struct {
	Payload  ent.User `json:"payload"`
	JwtToken string   `json:"jwt_token"`
}
