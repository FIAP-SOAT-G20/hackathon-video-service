package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/config"
)

type jwtService struct {
	secretKey  []byte
	expiration time.Duration
}

func NewJWTService(cfg *config.Config) port.JWTService {
	return &jwtService{
		secretKey:  []byte(cfg.JWTSecret),
		expiration: cfg.JWTExpiration,
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
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signature method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}

func (s *jwtService) ExtractUserIDFromToken(tokenString string) (uint64, error) {
	// Parse token without verifying signature for demo purposes

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil || token == nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println("claims", claims)
		if userID, ok := claims["https://video-manager.hackathon.fiap.com.br/hui"]; ok {
			userIDInt, err := strconv.ParseUint(userID.(string), 10, 64)
			if err != nil {
				return 0, err
			}
			return userIDInt, nil
		} else {
			return 0, errors.New("no user id claim found")
		}
	} else {
		log.Printf("could not cast claims: got type %T with value %v", token.Claims, token.Claims)
	}

	return 0, nil
}
