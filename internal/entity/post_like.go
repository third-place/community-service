package entity

type PostLike struct {
	ID     uint `gorm:"primary_key"`
	UserID uint `gorm:"unique_index:unique_user_post"`
	User   *User
	PostID uint `gorm:"unique_index:unique_user_post"`
	Post   *Post
}
