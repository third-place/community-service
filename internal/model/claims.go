package model

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserUuid string `json:"userUuid"`
	jwt.RegisteredClaims
}
