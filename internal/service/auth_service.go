package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/third-place/community-service/internal/auth"
	"github.com/third-place/community-service/internal/auth/model"
	"github.com/third-place/community-service/internal/db"
	model2 "github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/repository"
	"net/http"
	"os"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

type AuthService struct {
	client         *auth.APIClient
	cookieStore    *sessions.CookieStore
	userRepository *repository.UserRepository
}

func CreateDefaultAuthService() *AuthService {
	return &AuthService{
		client:         auth.NewAPIClient(auth.NewConfiguration()),
		userRepository: repository.CreateUserRepository(db.CreateDefaultConnection()),
	}
}

func (a *AuthService) GetSessionFromRequest(r *http.Request) (*model.Session, error) {
	sessionToken := a.getSessionToken(r)
	if sessionToken == "" {
		return nil, errors.New("no session token found")
	}
	session, err := a.getSession(sessionToken)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (a *AuthService) getSession(sessionToken string) (*model.Session, error) {
	claims := &model2.Claims{}
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

func (a *AuthService) getSessionToken(r *http.Request) string {
	return r.Header.Get("x-session-token")
}
