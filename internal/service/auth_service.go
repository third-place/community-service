package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/antihax/optional"
	"github.com/third-place/community-service/internal/auth"
	"github.com/third-place/community-service/internal/auth/model"
	"github.com/third-place/community-service/internal/db"
	"github.com/third-place/community-service/internal/repository"
	"github.com/third-place/community-service/internal/util"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

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

func (a *AuthService) GetSessionFromRequest(r *http.Request) *model.Session {
	sessionToken := a.getSessionToken(r)
	if sessionToken != "" {
		session, err := a.getSession(sessionToken)
		if err != nil {
			log.Print("error getting session :: ", err)
		}
		return session
	}
	return nil
}

func (a *AuthService) DoWithValidSession(w http.ResponseWriter, r *http.Request, doAction func(session *model.Session) (interface{}, error)) {
	sessionToken := a.getSessionToken(r)
	if sessionToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("missing required header: x-session-token"))
		return
	}
	session, err := a.getSession(sessionToken)
	if err != nil {
		log.Print("FAILED! Either error, or Uuid mismatch :: ", err)
		err := errors.New("invalid session")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	object, err := doAction(session)
	util.WriteResponse(w, object, err)
}

func (a *AuthService) getSession(sessionId string) (*model.Session, error) {
	ctx := context.TODO()
	response, _ := a.client.DefaultApi.GetSession(ctx, &auth.GetSessionOpts{
		Token: optional.NewString(sessionId),
	})
	if response == nil || response.StatusCode != http.StatusOK {
		return nil, errors.New("no session found")
	}
	session, _ := DecodeRequestToNewSession(response)
	return session, nil
}

func (a *AuthService) getSessionToken(r *http.Request) string {
	return r.Header.Get("x-session-token")
}

func DecodeRequestToNewSession(r *http.Response) (*model.Session, error) {
	decoder := json.NewDecoder(r.Body)
	var session *model.Session
	err := decoder.Decode(&session)
	if err != nil {
		return nil, err
	}
	return session, nil
}
