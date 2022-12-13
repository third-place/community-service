package util

import (
	"context"
	kafka2 "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/third-place/community-service/internal/db"
	"github.com/third-place/community-service/internal/entity"
	"github.com/third-place/community-service/internal/kafka"
	"github.com/third-place/community-service/internal/model"
	"math/rand"
	"strconv"
	"time"
)

var dbConn *gorm.DB

func CreateTestUser() *model.User {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Int()
	return &model.User{
		Uuid:     uuid.New().String(),
		Username: "user" + strconv.Itoa(randomInt),
	}
}

func SetupTestDatabase() *gorm.DB {
	if dbConn != nil {
		return dbConn
	}
	// 1. Create PostgreSQL container request
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}

	// 2. Start PostgreSQL container
	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})

	// 3.1 Get host and port of PostgreSQL container
	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	// 3.2 Create db connection string and connect
	dbConn = db.CreateConnection(
		host,
		port.Port(),
		"testdb",
		"postgres",
		"postgres",
	)

	migrateDb(dbConn)

	return dbConn
}

func migrateDb(dbConn *gorm.DB) {
	dbConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	dbConn.AutoMigrate(
		&entity.Follow{},
		&entity.Image{},
		&entity.Post{},
		&entity.PostLike{},
		&entity.Report{},
		&entity.User{},
	)
}

func CreateTestProducer() kafka.Producer {
	return &TestProducer{}
}

type TestProducer struct{}

func (t *TestProducer) Produce(msg *kafka2.Message, deliveryChan chan kafka2.Event) error {
	return nil
}
