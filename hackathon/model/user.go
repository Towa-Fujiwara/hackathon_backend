package model

import (
	"fmt"
	"time"
)

type User struct {
	Id      string      `json:"id"`
	Name    string      `json:"name"`
	Age     int         `json:"age"`
	Password string `json:"password"`
	Profile UserProfile `json:"profile"`
	CreatedAt time.Time `json:"created_at"`
}
type UserProfile struct {
	IconUrl string `json:"icon_url"`
	DisplayName string `json:"display_name"`
	Bio string `json:"bio"`
}


func NewUser(id, name, password, displayName, bio, iconURL string, age int) (*User, error) {
    
    if id == "" {
        return nil, fmt.Errorf("id is empty")
    }
    if name == "" || len(name) > 50 {
        return nil, fmt.Errorf("invalid name: %s", name)
    }
    if age < 10 {
        return nil, fmt.Errorf("you cannot register under 10 years old: %d", age)
    }
    if password == "" {
        return nil, fmt.Errorf("password is empty")
    }

    return &User{
        Id:       id,
        Name:     name,
        Age:      age,
        Password: password,
        Profile: UserProfile{
            DisplayName:        displayName,
            Bio:                bio,
            IconUrl:            iconURL,
        },
		CreatedAt: time.Now(), 
    }, nil
}
