package service

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/golang-jwt/jwt/v4"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
)

type CustomClaims struct {
	UserID string `json:"https://video-manager.hackathon.fiap.com.br/hui"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

type jwtService struct {
	secretKey  []byte
	expiration time.Duration
	validator  *validator.Validator
}

func NewJWTService(cfg *config.Config) port.JWTService {

	issuerURL, err := url.Parse("https://" + cfg.Auth0Domain + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{cfg.Auth0Audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	return &jwtService{
		secretKey:  []byte(cfg.JWTSecret),
		expiration: cfg.JWTExpiration,
		validator:  jwtValidator,
	}
}

func (s *jwtService) GenerateToken(id uint64) (string, error) {
	expiresAt := time.Now().Add(s.expiration)
	idStr := strconv.FormatUint(id, 10)

	tokenClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        idStr,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *jwtService) ValidateToken(tokenString string) error {
	_, err := s.validator.ValidateToken(context.Background(), tokenString)
	if err != nil {
		return err
	}
	return nil
}

func (s *jwtService) ExtractUserIDFromToken(tokenString string) (uint64, error) {
<<<<<<< HEAD
	// Parse token without verifying signature for demo purposes
	vc, err := s.validator.ValidateToken(context.Background(), tokenString)
	if err != nil {
		return 0, err
	}

	validatedClaims := vc.(*validator.ValidatedClaims)
	tokenClaims := validatedClaims.CustomClaims.(*CustomClaims)

	if tokenClaims.UserID != "" {
		userIDInt, err := strconv.ParseUint(tokenClaims.UserID, 10, 64)
=======
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signature method")
		}
		return s.secretKey, nil
	})
	if err != nil || token == nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}
	if userID, ok := claims["https://video-manager.hackathon.fiap.com.br/hui"]; ok {
		userIDStr, ok := userID.(string)
		if !ok {
			return 0, errors.New("user id claim is not a string")
		}
		userIDInt, err := strconv.ParseUint(userIDStr, 10, 64)
>>>>>>> f47cc4c3e4266c5e3fd1ddfcb1ff97420c261f6e
		if err != nil {
			return 0, err
		}
		return userIDInt, nil
	} else {
		return 0, errors.New("no user id claim found")
	}
}
