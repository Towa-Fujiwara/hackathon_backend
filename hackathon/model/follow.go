package model

import "time"

type Follow struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	FollowUserId string `json:"follow_user_id"`
	CreatedAt time.Time `json:"created_at"`
}