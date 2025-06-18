package controller

import (
	"net/http"
	"hackathon/usecase"
	"hackathon/model"
	"github.com/go-chi/chi/v5"
)

type FollowUserController struct {
	followUserUsecase usecase.FollowUserUsecase
}

func NewFollowUserController(followUserUsecase usecase.FollowUserUsecase) *FollowUserController {
	return &FollowUserController{followUserUsecase: followUserUsecase}
}

func (c *FollowUserController) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(userContextKey).(string)
	if !ok || uid == "" {
		http.Error(w, "User ID not found in context. This endpoint requires authentication.", http.StatusInternalServerError)
		return
	}
	var follow model.Follow
	if err := decodeBody(r, &follow); err != nil {
		respondJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	follow.UserId = uid
	followUser, err := c.followUserUsecase.FollowUser(&follow)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, followUser)
}

func (c *FollowUserController) GetFollowersHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(userContextKey).(string)
	if !ok || uid == "" {
		http.Error(w, "User ID not found in context. This endpoint requires authentication.", http.StatusInternalServerError)
		return
	}
	followers, err := c.followUserUsecase.GetFollowers(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, followers)
}
func (c *FollowUserController) GetFollowingHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(userContextKey).(string)
	if !ok || uid == "" {
		http.Error(w, "User ID not found in context. This endpoint requires authentication.", http.StatusInternalServerError)
		return
	}
	following, err := c.followUserUsecase.GetFollowing(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, following)
}

func (c *FollowUserController) IsFollowingHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(userContextKey).(string)
	if !ok || uid == "" {
		http.Error(w, "User ID not found in context. This endpoint requires authentication.", http.StatusInternalServerError)
		return
	}
	followUserId := chi.URLParam(r, "followUserId")
	isFollowing, err := c.followUserUsecase.IsFollowing(uid, followUserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, isFollowing)
}
