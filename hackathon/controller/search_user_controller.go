package controller

import (
	"net/http"
	"hackathon/usecase"
	"log"
	"strings"
	"github.com/go-chi/chi/v5"
)

type SearchUserController struct {
	searchUserUsecase usecase.UserUsecase
}

func NewSearchUserController(su usecase.UserUsecase) *SearchUserController {
	return &SearchUserController{searchUserUsecase: su}
}

func (c *SearchUserController) GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/id/") {
		// /api/users/id/{userId} の場合の処理
		userId := chi.URLParam(r, "userId")
		if userId == "" {
			http.Error(w, "User ID is required in the URL path", http.StatusBadRequest)
			return
		}

		user, err := c.searchUserUsecase.GetUserByUserId(userId)
		if err != nil {
			log.Printf("ERROR: GetUserByUserId failed: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		respondJSON(w, http.StatusOK, user)

	} else {
	uid, ok := r.Context().Value(userContextKey).(string)
	if !ok || uid == "" {
		http.Error(w, "User ID not found in context. This endpoint requires authentication.", http.StatusInternalServerError)
		return
	}
	user, err := c.searchUserUsecase.GetUserByFirebaseUID(uid)
	if err != nil {
		log.Printf("ERROR: SearchUserExist failed: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if user == nil{
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	respondJSON(w, http.StatusOK, user)
}
}

func (c *SearchUserController) SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondJSON(w, http.StatusBadRequest, "Bad Request")
		return
	}
	users, err := c.searchUserUsecase.SearchUsers(query)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, "Server Error")
		return
	}
	respondJSON(w, http.StatusOK, users)
}
