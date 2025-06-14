package model

import (
	"fmt"
	"time"
)

type User struct {
	Id      string      `json:"id"`
	Name    string      `json:"name"`
	Profile UserProfile `json:"profile"`
	CreatedAt time.Time `json:"created_at"`
}
type UserProfile struct {
	IconUrl string `json:"icon_url"`
	Bio string `json:"bio"`
}


func NewUser(id, name, bio, iconURL string) (*User, error) {
    
    if id == "" {
        return nil, fmt.Errorf("id is empty")
    }
    if name == "" || len(name) > 50 {
        return nil, fmt.Errorf("invalid name: %s", name)
    }

    return &User{
        Id:       id,
        Name:     name,
        Profile: UserProfile{
            Bio:                bio,
            IconUrl:            iconURL,
        },
		CreatedAt: time.Now(), 
    }, nil
}
