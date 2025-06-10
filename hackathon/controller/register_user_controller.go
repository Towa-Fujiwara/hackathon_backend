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

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("fail: json.Decode, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	row, err := c.registerUserUsecase.RegisterUser(user.Id, user.Name, user.Password, user.Profile.DisplayName, user.Profile.Bio, user.Profile.IconUrl, user.Profile.BackgroundImageUrl, user.Age); 
	if err != nil {
		log.Printf("fail: registerUserUsecase.RegisterUser, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusCreated, map[string]string{"id": row.Id})
}


