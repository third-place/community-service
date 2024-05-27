package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/third-place/community-service/internal/service"
)

func main() {
	service.CreateConsumerService().Loop()
}
