package main

import (
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/third-place/community-service/internal/service"
	"log"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	postUuid := uuid.MustParse(os.Args[1])
	replyService := service.CreateReplyService()
	reply, err := replyService.GetReply(postUuid)
	if err != nil {
		log.Fatal("reply not found")
	}
	err = replyService.PublishReplyToKafka(reply)
	time.Sleep(2 * time.Second)
	if err != nil {
		log.Fatal("message not sent")
	}
	log.Print("reply re-published")
}
