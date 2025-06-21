package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"hackathon/usecase"
	"hackathon/model"
)

type RegisterUserController struct {
	registerUserUsecase usecase.UserUsecase
}

func NewRegisterUserController(ru usecase.UserUsecase) *RegisterUserController {
	return &RegisterUserController{registerUserUsecase: ru}
}

func (c *RegisterUserController) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("fail: HTTP Method is %s\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 認証ミドルウェアからユーザーのUIDを取得
	userUID, ok := r.Context().Value(userContextKey).(string)
	if !ok || userUID == "" {
		log.Printf("fail: user UID not found in context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("fail: json.Decode, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// FirebaseUIDを自動設定
	user.FirebaseUID = userUID

	row, err := c.registerUserUsecase.RegisterUser(user.UserId, user.FirebaseUID, user.Name, user.Bio, user.IconUrl); 
	if err != nil {
		log.Printf("fail: registerUserUsecase.RegisterUser, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusCreated, map[string]string{"userId": row.UserId})
}

func (c *RegisterUserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.Printf("fail: HTTP Method is %s\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 認証ミドルウェアからユーザーのUIDを取得
	userUID, ok := r.Context().Value(userContextKey).(string)
	if !ok || userUID == "" {
		log.Printf("fail: user UID not found in context")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("fail: json.Decode, %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ユーザー情報を更新
	updatedUser, err := c.registerUserUsecase.UpdateUser(userUID, user.Name, user.Bio, user.IconUrl)
	if err != nil {
		log.Printf("fail: registerUserUsecase.UpdateUser, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, updatedUser)
}


