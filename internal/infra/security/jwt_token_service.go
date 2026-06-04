package security

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	appauth "github.com/viictormotorhead/financial-app/internal/application/auth"
)

const jwtIssuer = "financial-app"

type JWTTokenService struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTTokenService(secret string, ttl time.Duration) *JWTTokenService {
	return &JWTTokenService{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

var (
	_ appauth.TokenIssuer = (*JWTTokenService)(nil)
	_ appauth.TokenParser = (*JWTTokenService)(nil)
)

func (s *JWTTokenService) Issue(ctx context.Context, claims appauth.TokenClaims) (string, error) {
	userID := strings.TrimSpace(claims.UserID)
	username := strings.TrimSpace(claims.Username)
	if userID == "" || username == "" {
		return "", fmt.Errorf("user id and username are required to issue token")
	}

	now := time.Now().UTC()
	mapClaims := jwt.MapClaims{
		"iss":      jwtIssuer,
		"sub":      userID,
		"username": username,
		"iat":      now.Unix(),
		"exp":      now.Add(s.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

func (s *JWTTokenService) Parse(tokenString string) (appauth.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return appauth.TokenClaims{}, appauth.ErrUnauthorized
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return appauth.TokenClaims{}, appauth.ErrUnauthorized
	}

	userID, _ := mapClaims["sub"].(string)
	username, _ := mapClaims["username"].(string)
	userID = strings.TrimSpace(userID)
	username = strings.TrimSpace(username)
	if userID == "" || username == "" {
		return appauth.TokenClaims{}, appauth.ErrUnauthorized
	}

	return appauth.TokenClaims{
		UserID:   userID,
		Username: username,
	}, nil
}
