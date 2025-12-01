package domain

type TokenService interface {
	GenerateToken(userID string) (string, error)
}
