package client

type Agent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ResponseDetail struct {
	Detail string `json:"detail"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type UserWithoutPassword struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateData struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Group struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

type Bid struct {
	ID        int    `json:"id"`
	GroupName string `json:"group_name"`
	UserId    string `json:"user_id"`
	Status    string `json:"status"`
}

type RequestMeta struct {
	UserAgent string
	RealIp    string
}

type RequestStatus struct {
	Err        error
	StatusCode int
}

func newRequestStatus(err error, status int) *RequestStatus {
	return &RequestStatus{
		Err:        err,
		StatusCode: status,
	}
}
