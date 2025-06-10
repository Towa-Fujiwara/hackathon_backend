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
	BackgroundImageUrl string `json:"background_image_url"`
}

func NewUser(id, name, password, displayName, bio, iconURL, backgroundURL string, age int) (*User, error) {
    
    // バリデーション
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

    // バリデーションが通ったら、受け取ったIDを使ってインスタンスを生成
    return &User{
        Id:       id, // 引数のidを設定
        Name:     name,
        Age:      age,
        Password: password, // 注意: パスワードはハッシュ化すべき
        Profile: UserProfile{
            DisplayName:        displayName,
            Bio:                bio,
            IconUrl:            iconURL,
            BackgroundImageUrl: backgroundURL,
        },
		CreatedAt: time.Now(), 
    }, nil
}
