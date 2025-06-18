package model

import (
	"time"
)


type Post struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	UserName  string    `json:"name"`
	IconUrl  string `json:"iconUrl"`
	Text string `json:"text"`
	Image string `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
	LikeCount    int       `json:"likeCount"`
	CommentCount int       `json:"commentCount"`
}

type Comment struct {
    Id        string    `json:"id"`
    UserId    string    `json:"userId"`
    PostId    string    `json:"postId"`
    Text      string    `json:"text"`
    CreatedAt time.Time `json:"createdAt"`
}