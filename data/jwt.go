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

// TemplateData return the token key and the actual data to execute the template with
type TemplateData struct {
	Key  string
	Data interface{}
}
