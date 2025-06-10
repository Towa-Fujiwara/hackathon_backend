package model

import (
	"time"
)


type Post struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Text string `json:"text"`
	Image string `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
}


type Comment struct {
    Id        string    `json:"id"`
    UserId    string    `json:"user_id"`
    PostId    string    `json:"post_id"`
    Text      string    `json:"text"`
    CreatedAt time.Time `json:"created_at"`
}