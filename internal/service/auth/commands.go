package auth

// RegisterCommand carries the fields needed to create a new account.
type RegisterCommand struct {
	Email       string
	Password    string
	DisplayName string
}

// LoginCommand carries the credentials for a password login.
type LoginCommand struct {
	Email    string
	Password string
}

// RefreshCommand carries the refresh token to rotate.
type RefreshCommand struct {
	RefreshToken string
}

// LogoutCommand carries the refresh token to invalidate.
type LogoutCommand struct {
	RefreshToken string
}
