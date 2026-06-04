package response

type AuthUserResponse struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type AuthResponse struct {
	Token string           `json:"token"`
	User  AuthUserResponse `json:"user"`
}
