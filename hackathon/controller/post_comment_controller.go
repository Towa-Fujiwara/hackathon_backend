package controller

import (
	"net/http"
	"github.com/go-chi/chi"
	"hackathon/usecase"
	"hackathon/model"
	"encoding/json"
	"log"
)

type PostCommentController struct {
	commentUsecase usecase.CommentUsecase
}

func NewPostCommentController(cu usecase.CommentUsecase) *PostCommentController {
	return &PostCommentController{commentUsecase: cu}
}

func (c *PostCommentController) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(userContextKey).(string)
	if !ok || uid == "" {
		http.Error(w, "User ID not found in context. This endpoint requires authentication.", http.StatusInternalServerError)
		return
	}
	postId := chi.URLParam(r, "postId")
	if postId == "" {
		http.Error(w, "Bad Request: Post ID is required", http.StatusBadRequest)
		return
	}
	var comment model.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Bad Request: Comment is required", http.StatusBadRequest)
		return
	}
	comment.UserId = uid
	createdComment, err := c.commentUsecase.CreateComment(&comment)
	if err != nil {
		log.Printf("failed to create comment: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		return
	}
	respondJSON(w, http.StatusCreated, createdComment)
}

func (c *PostCommentController) GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	comments, err := c.commentUsecase.GetCommentsByPostId(postId)
	if err != nil {
		log.Printf("failed to get comments: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		return
	}
	respondJSON(w, http.StatusOK, comments)
}
	
