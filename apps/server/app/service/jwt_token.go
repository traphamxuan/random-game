package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/traphamxuan/random-game/package/config"
)

type JwtToken struct {
	config    *config.Configuration
	secretKey string
}

var _ IService = (*JwtToken)(nil)

type JwtTokenConfig struct {
	SecretKey string `env:"JWT_SECRET_KEY" envDefault:"secret-key"`
}

func NewJwtToken(cfg *config.Configuration) JwtToken {
	return JwtToken{
		config: cfg,
	}
}

func (a *JwtToken) Start() error {
	var cfg JwtTokenConfig
	if err := a.config.ParseConfig(&cfg); err != nil {
		return err
	}
	a.secretKey = cfg.SecretKey
	return nil
}

func (a *JwtToken) ComposeToken(userID string) (string, time.Time, error) {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(24 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &jwt.RegisteredClaims{
		ID: userID,
		// In JWT, the expiry time is expressed as unix milliseconds
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenStr, err := token.SignedString(a.secretKey)
	return tokenStr, expirationTime, err
}

func (a *JwtToken) VerifyToken(tokenStr string) (claims jwt.RegisteredClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
		return a.secretKey, nil
	})
	if err != nil {
		return claims, err
	}
	if !token.Valid {
		return claims, jwt.ErrTokenSignatureInvalid
	}
	return claims, nil
}

func (a *JwtToken) RefreshToken(tokenStr string) (string, time.Time, error) {
	claims, err := a.VerifyToken(tokenStr)
	if err != nil {
		return "", time.Time{}, err
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return tokenStr, claims.ExpiresAt.Time, nil
	}

	return a.ComposeToken(claims.ID)
}