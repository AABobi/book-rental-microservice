package data

type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}
