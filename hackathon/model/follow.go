package model

import "time"

type Follow struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	FollowedUserId string `json:"followed_user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type FollowRequest struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	FollowedUserId string `json:"followed_user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type FollowRequestResponse struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	FollowedUserId string `json:"followed_user_id"`
	CreatedAt time.Time `json:"created_at"`
}