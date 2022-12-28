package main

import (
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/third-place/community-service/internal/model"
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
	session := model.CreateSession(uuid.MustParse("4e1dce1a-6041-4d63-ad00-91715274f643"))
	postUuid := uuid.MustParse(os.Args[1])
	postService := service.CreatePostService()
	post, err := postService.GetPost(session, postUuid)
	if err != nil {
		log.Fatal("post not found")
	}
	err = postService.PublishPostToKafka(post)
	time.Sleep(2 * time.Second)
	if err != nil {
		log.Fatal("message not sent")
	}
	log.Print("post re-published")
}
