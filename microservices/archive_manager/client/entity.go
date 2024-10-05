package client

type Record struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type ResponseError struct {
	Error string `json:"error"`
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
