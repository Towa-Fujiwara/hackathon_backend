package model

import "time"

type Follow struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	FollowUserId string `json:"followUserId"`
	CreatedAt time.Time `json:"createdAt"`
}