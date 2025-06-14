package model

import (
	"fmt"
	"time"
)

type User struct {
    UserId   string `json:"userId"`
    Name string `json:"name"`
    Bio      string `json:"bio"`
    IconUrl  string `json:"iconUrl"`
	CreatedAt time.Time `json:"createdAt"`
}


func NewUser(userId, name, bio, iconURL string) (*User, error) {
    
    if userId == "" {
        return nil, fmt.Errorf("ID is empty")
    }
    if name == "" || len(name) > 50 {
        return nil, fmt.Errorf("invalid name: %s", name)
    }

    return &User{
		UserId:   userId,
		Name:     name,
		Bio:      bio,
		IconUrl:  iconURL,
		CreatedAt: time.Now(), 
    }, nil
}
