// jwtutil helper 
// pkg/jwtutil/jwtutil.go
package jwtutil

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims defines custom claims, embedding StandardClaims.
// You can add more fields here if you want (e.g. Role, Email, etc.).
type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

var (
	// ErrTokenExpired is returned when JWT is expired
	ErrTokenExpired = errors.New("token has expired")

	// ErrTokenMalformed is returned when JWT is not well‐formed
	ErrTokenMalformed = errors.New("token is malformed")

	// ErrTokenInvalid is returned when JWT is invalid for any other reason
	ErrTokenInvalid = errors.New("token is invalid")
)

// getSigningKey reads the JWT_SECRET from env or panics if not set.
func getSigningKey() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("environment variable JWT_SECRET is not set")
	}
	return []byte(secret), nil
}

// getExpiryDuration reads JWT_EXPIRES_IN (minutes) from env, defaults to 60 if missing/invalid.
func getExpiryDuration() time.Duration {
	s := os.Getenv("JWT_EXPIRES_IN")
	if s == "" {
		return time.Minute * 60 // default: 60 minutes
	}

	minutes, err := strconv.Atoi(s)
	if err != nil || minutes <= 0 {
		return time.Minute * 60
	}
	return time.Minute * time.Duration(minutes)
}

// GenerateToken creates a signed JWT string for the given user ID.
// It uses HS256 algorithm.
func GenerateToken(userID int) (string, error) {
	key, err := getSigningKey()
	if err != nil {
		return "", err
	}

	expiry := time.Now().Add(getExpiryDuration())
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			// You can add Audience, Issuer, Subject here if desired:
			// Issuer:   "your-app-name",
			// Subject:  strconv.Itoa(userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}

	return signedStr, nil
}

// ValidateToken parses and validates the JWT string. If valid, it returns the claims.
func ValidateToken(tokenStr string) (*JWTClaims, error) {
	key, err := getSigningKey()
	if err != nil {
		return nil, err
	}

	parsedToken, err := jwt.ParseWithClaims(
		tokenStr,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC and specifically HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return key, nil
		},
	)

	if err != nil {
		// Distinguish between different JWT errors
		ve, ok := err.(*jwt.ValidationError)
		if ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
			return nil, ErrTokenMalformed
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok || !parsedToken.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// ExtractUserID extracts the user_id claim from a validated token string.
// Returns an error if the token is invalid or the claim is missing.
func ExtractUserID(tokenStr string) (int, error) {
	claims, err := ValidateToken(tokenStr)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
