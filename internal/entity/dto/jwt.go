package dto

import "time"

// INPUT DATAFLOW
type JwtPayload struct {
	UserIpAddress string `json:"user_ip_address" valid:"ip"`
	UserId        string `json:"user_id" valid:"uuidv4"`
}

type JwtHeader struct {
	Alg  string    `json:"alg" valid:"-"`
	Type string    `json:"type" valid:"-"`
	Exp  time.Time `json:"exp" valid:"-"`
}
