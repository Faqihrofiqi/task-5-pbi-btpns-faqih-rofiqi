package makan

import "github.com/golang-jwt/jwt/v5"

var JWT_KEY = []byte("dvdfvmkfvmfkvmfk3029390")

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
