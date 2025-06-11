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
}

func NewPostLikeController(pl usecase.PostLikeUsecase) *PostLikeController {
	return &PostLikeController{postLikeUsecase: pl}
}

func (c *PostLikeController) LikePostHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(userContextKey).(string)
	if !ok || uid == "" {
		// このエラーは、ミドルウェアが正しく設定されていれば通常は発生しません。
		http.Error(w, "User ID not found in context. This endpoint requires authentication.", http.StatusInternalServerError)
		return
	}
	postId := chi.URLParam(r, "postId")
	if postId == "" {
        http.Error(w, "Bad Request: Post ID is required", http.StatusBadRequest)
        return
    }
	toggledLike, err := c.postLikeUsecase.ToggleLike(uid, postId)
    if err != nil {
        log.Printf("ERROR: ToggleLike failed: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // ToggleLike の戻り値で処理を分岐
    if toggledLike != nil {
        // オブジェクトが返ってきたら、いいねが作成された
        respondJSON(w, http.StatusCreated, toggledLike)
    } else {
        // nil が返ってきたら、いいねが削除された
        respondJSON(w, http.StatusOK, map[string]string{"message": "like removed"})
    }
}


