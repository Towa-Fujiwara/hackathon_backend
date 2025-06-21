package controller

import (
	"net/http"
	"strings"
	"github.com/go-chi/chi/v5"
	"hackathon/usecase"
	"log"
)

type FollowUserController struct {
	followUserUsecase usecase.FollowUserUsecase
}

type IsFollowingResponse struct {
    IsFollowing bool `json:"isFollowing"`
}

type UserFollowCountsResponse struct {
    FollowingCount int `json:"followingCount"`
    FollowersCount int `json:"followersCount"`
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

	targetUserId := chi.URLParam(r, "userId") // URLパスから対象ユーザーIDを取得
    if targetUserId == "" {
        http.Error(w, "対象ユーザーIDが指定されていません", http.StatusBadRequest)
        return
    }

	err := c.followUserUsecase.FollowUser(r.Context(), uid, targetUserId)
    if err != nil {
        log.Printf("フォロー失敗: %v", err)
        if strings.Contains(err.Error(), "自分自身") {
            http.Error(w, "自分自身をフォローすることはできません", http.StatusBadRequest)
            return
        }
        if strings.Contains(err.Error(), "既にフォロー") {
            http.Error(w, "既にフォローしています", http.StatusConflict)
            return
        }
        http.Error(w, "フォローに失敗しました", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (c *FollowUserController) UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(userContextKey).(string)
	if !ok || uid == "" {
		http.Error(w, "User ID not found in context. This endpoint requires authentication.", http.StatusInternalServerError)
		return
	}
	targetUserId := chi.URLParam(r, "userId")
    if targetUserId == "" {
        http.Error(w, "対象ユーザーIDが指定されていません", http.StatusBadRequest)
        return
    }

	err := c.followUserUsecase.UnfollowUser(r.Context(), uid, targetUserId)
    if err != nil {
        log.Printf("フォロー解除失敗: %v", err)
        if strings.Contains(err.Error(), "自分自身") {
            http.Error(w, "自分自身をフォローすることはできません", http.StatusBadRequest)
            return
        }
        if strings.Contains(err.Error(), "フォローしていません") {
            http.Error(w, "フォローしていません", http.StatusConflict)
            return
        }
        http.Error(w, "フォロー解除に失敗しました", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (c *FollowUserController) IsFollowingHandler(w http.ResponseWriter, r *http.Request) {
	targetUserId := chi.URLParam(r, "userId")
	if targetUserId == "" {
		http.Error(w, "対象ユーザーIDが指定されていません", http.StatusBadRequest)
		return
	}

	uid, ok := r.Context().Value(userContextKey).(string)
	if !ok || uid == "" {
		respondJSON(w, http.StatusOK, IsFollowingResponse{IsFollowing: false})
		return
	}

	isFollowing, err := c.followUserUsecase.IsFollowing(r.Context(), uid, targetUserId)
	if err != nil {
		log.Printf("フォロー状態の確認に失敗: %v", err)
		http.Error(w, "フォロー状態の確認に失敗しました", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, IsFollowingResponse{IsFollowing: isFollowing})
}

func (c *FollowUserController) GetFollowersHandler(w http.ResponseWriter, r *http.Request) {
	targetUserId := chi.URLParam(r, "userId")
	if targetUserId == "" {
		http.Error(w, "対象ユーザーIDが指定されていません", http.StatusBadRequest)
		return
	}

	followers, err := c.followUserUsecase.GetFollowers(r.Context(), targetUserId)
	if err != nil {
		log.Printf("フォロワーリストの取得に失敗: %v", err)
		http.Error(w, "フォロワーリストの取得に失敗しました", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, followers)
}

func (c *FollowUserController) GetFollowingHandler(w http.ResponseWriter, r *http.Request) {
	targetUserId := chi.URLParam(r, "userId")
	if targetUserId == "" {
		http.Error(w, "対象ユーザーIDが指定されていません", http.StatusBadRequest)
		return
	}

	following, err := c.followUserUsecase.GetFollowing(r.Context(), targetUserId)
	if err != nil {
		log.Printf("フォローリストの取得に失敗: %v", err)
		http.Error(w, "フォローリストの取得に失敗しました", http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, following)
}

func (c *FollowUserController) GetUserFollowCountsHandler(w http.ResponseWriter, r *http.Request) {
    targetUserId := chi.URLParam(r, "userId")
    if targetUserId == "" {
        http.Error(w, "対象ユーザーIDが指定されていません", http.StatusBadRequest)
        return
    }

    followingCount, followersCount, err := c.followUserUsecase.GetUserFollowCounts(r.Context(), targetUserId)
    if err != nil {
        log.Printf("ユーザーのフォローカウントの取得に失敗: %v", err)
        http.Error(w, "ユーザーのフォローカウントの取得に失敗しました", http.StatusInternalServerError)
        return
    }

    respondJSON(w, http.StatusOK, UserFollowCountsResponse{
        FollowingCount: followingCount,
        FollowersCount: followersCount,
    })
}