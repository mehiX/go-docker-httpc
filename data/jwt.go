package data

import (
	jwt "github.com/dgrijalva/jwt-go"
)

// UserClaims include custom and standard claims for JWT
type UserClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
