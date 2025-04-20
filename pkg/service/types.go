package service

import "github.com/dgrijalva/jwt-go"

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id`
}
