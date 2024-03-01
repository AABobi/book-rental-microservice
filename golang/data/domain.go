package data

type AuthResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}
