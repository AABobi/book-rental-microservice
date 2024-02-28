package data

type User struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
