package test

import (
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/model"
	"math/rand"
	"strconv"
	"time"
)

func CreateTestUser() *model.User {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Int()
	return &model.User{
		Uuid:     uuid.New().String(),
		Username: "user" + strconv.Itoa(randomInt),
	}
}
