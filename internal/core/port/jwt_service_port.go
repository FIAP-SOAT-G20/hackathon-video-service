package port

// JWTService provides token generation and validation methods
type JWTService interface {
	// GenerateToken creates a new JWT token with the given claims
	GenerateToken(customerID uint64) (string, error)

	// ValidateToken verifies if a token is valid without extracting data
	ValidateToken(token string) error
}
