package model

import "time"

type Like struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	PostId    string    `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}
