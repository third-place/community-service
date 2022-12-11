package util

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/model"
	"os"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func GetSession(sessionToken string) (*model.Session, error) {
	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(sessionToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token not valid")
	}
	_, err = uuid.Parse(claims.UserUuid)
	if err != nil {
		return nil, err
	}
	return &model.Session{
		User: model.User{
			Uuid: claims.UserUuid,
		},
		Token: sessionToken,
	}, nil
}
