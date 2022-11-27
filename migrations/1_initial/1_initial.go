package main

import (
	"github.com/third-place/community-service/internal/db"
	"github.com/third-place/community-service/internal/entity"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	conn := db.CreateDefaultConnection()
	conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\" WITH SCHEMA public;")
	conn.AutoMigrate(
		&entity.User{},
		&entity.Post{},
		&entity.Follow{},
		&entity.Report{},
		&entity.Image{},
		&entity.PostLike{},
	)
	conn.Model(&entity.Post{}).
		AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	conn.Model(&entity.Follow{}).
		AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("following_id", "users(id)", "RESTRICT", "RESTRICT")
	conn.Model(&entity.Report{}).
		AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("reported_id", "posts(id)", "RESTRICT", "RESTRICT")
	conn.Model(&entity.Image{}).
		AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("post_id", "posts(id)", "RESTRICT", "RESTRICT")
	conn.Model(&entity.PostLike{}).
		AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").
		AddForeignKey("post_id", "posts(id)", "RESTRICT", "RESTRICT")
}
