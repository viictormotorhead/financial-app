package outputs

type AuthUserOutputDTO struct {
	UserID   string
	Name     string
	Username string
}

type AuthOutputDTO struct {
	Token string
	User  AuthUserOutputDTO
}
