package controller

import (
	"net/http"
	"hackathon/model"
	"hackathon/usecase"
	"github.com/go-chi/chi"
	"log"
)

type PostController struct {
	postUsecase usecase.PostUsecase
	userUsecase usecase.UserUsecase
}

func NewPostController(pu usecase.PostUsecase, uu usecase.UserUsecase) *PostController {
	return &PostController{
	postUsecase: pu,
    userUsecase: uu, 
	}
}

func (c *PostController) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	post, err := c.postUsecase.FindById(chi.URLParam(r, "postId"))
	if err != nil {
		log.Printf("ERROR: FindById failed: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, post)
}


func (c *PostController) GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {

	posts, err := c.postUsecase.FindAllPosts()
	if err != nil {
		log.Printf("ERROR: FindAllPosts failed: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, posts)
}


func (c *PostController) GetAllPostsByUserIdHandler(w http.ResponseWriter, r *http.Request) {
	firebaseUID, ok := r.Context().Value(userContextKey).(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "User ID not found in context. This endpoint requires authentication.", http.StatusInternalServerError)
		return
	}


	appUser, err := c.userUsecase.GetUserByFirebaseUID(firebaseUID)
    if err != nil || appUser == nil {
        log.Printf("ERROR: Failed to find user by firebaseUID %s: %v\n", firebaseUID, err)
        http.Error(w, "Authenticated user not found in application database.", http.StatusInternalServerError)
        return
    }

	posts, err := c.postUsecase.FindAllPostsByUserId(appUser.UserId)
	if err != nil {
		log.Printf("ERROR: FindAllPostsByUserId failed: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, posts)
}


func (c *PostController) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	firebaseUID, ok := r.Context().Value(userContextKey).(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "User ID not found in context. This endpoint requires authentication.", http.StatusInternalServerError)
		return
	}

	appUser, err := c.userUsecase.GetUserByFirebaseUID(firebaseUID)
    if err != nil || appUser == nil {
        log.Printf("ERROR: Failed to find user by firebaseUID %s: %v\n", firebaseUID, err)
        http.Error(w, "Authenticated user not found in application database.", http.StatusInternalServerError)
        return
    }

	var post model.Post
	if err := decodeBody(r, &post); err != nil {
		respondJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	post.UserId = appUser.UserId
	result, err := c.postUsecase.CreatePost(&post)
	if err != nil {
		log.Printf("ERROR: CreatePost failed: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusCreated, result)
}