package controller

import (
    "net/http"
    "hackathon/usecase"
    "github.com/go-chi/chi/v5"
	//"hackathon/model"
	"log"
)

type PostLikeController struct {
	postLikeUsecase usecase.PostLikeUsecase
    userUsecase     usecase.UserUsecase
}

func NewPostLikeController(pl usecase.PostLikeUsecase, uu usecase.UserUsecase) *PostLikeController {
	return &PostLikeController{
        postLikeUsecase: pl,
        userUsecase:     uu,
    }
}

func (c *PostLikeController) LikePostHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(userContextKey).(string)
	if !ok || uid == "" {
		http.Error(w, "User ID not found in context. This endpoint requires authentication.", http.StatusInternalServerError)
		return
	}
    appUser, err := c.userUsecase.GetUserByFirebaseUID(uid)
	if err != nil || appUser == nil {
		log.Printf("ERROR: Failed to find user by firebaseUID %s: %v\n", uid, err)
		http.Error(w, "Authenticated user not found in application database.", http.StatusInternalServerError)
		return
	}
	postId := chi.URLParam(r, "postId")
	if postId == "" {
        http.Error(w, "Bad Request: Post ID is required", http.StatusBadRequest)
        return
    }
	toggledLike, err := c.postLikeUsecase.ToggleLike(appUser.UserId, postId)
    if err != nil {
        log.Printf("ERROR: ToggleLike failed: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    if toggledLike != nil {
        respondJSON(w, http.StatusCreated, toggledLike)
    } else {
        respondJSON(w, http.StatusOK, map[string]string{"message": "like removed"})
    }
}


