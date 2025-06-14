package model

import (
	"fmt"
	"time"
)

type User struct {
    UserId   string `json:"userId"`      // ユーザーが設定するID（Xの@usernameのようなもの）
    FirebaseUID string `json:"firebaseUid"` // FirebaseのUID（内部的に使用）
    Name string `json:"name"`
    Bio      string `json:"bio"`
    IconUrl  string `json:"iconUrl"`
	CreatedAt time.Time `json:"createdAt"`
}


func NewUser(userId, firebaseUID, name, bio, iconURL string) (*User, error) {
    
    if userId == "" {
        return nil, fmt.Errorf("userId is empty")
    }
    if firebaseUID == "" {
        return nil, fmt.Errorf("firebaseUID is empty")
    }
    if name == "" || len(name) > 50 {
        return nil, fmt.Errorf("invalid name: %s", name)
    }

    return &User{
		UserId:   userId,
		FirebaseUID: firebaseUID,
		Name:     name,
		Bio:      bio,
		IconUrl:  iconURL,
		CreatedAt: time.Now(), 
    }, nil
}
