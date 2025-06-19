package controller

import (
	"encoding/json"
	"net/http"
	"hackathon/usecase"
	"github.com/go-chi/chi/v5"
)

type GeminiController struct {
	geminiUsecase *usecase.GeminiUsecase
}

func NewGeminiController(geminiUsecase *usecase.GeminiUsecase) *GeminiController {
	return &GeminiController{
		geminiUsecase: geminiUsecase,
	}
}

// ユーザーの投稿からサマリーを生成するエンドポイント
func (gc *GeminiController) GenerateUserSummaryHandler(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	// ユーザーサマリーを生成
	summary, err := gc.geminiUsecase.GenerateUserSummary(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JSONレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// 自分の投稿からサマリーを生成するエンドポイント
func (gc *GeminiController) GenerateMySummaryHandler(w http.ResponseWriter, r *http.Request) {
	// 認証ミドルウェアからユーザーIDを取得
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "userID not found in context", http.StatusUnauthorized)
		return
	}

	summary, err := gc.geminiUsecase.GenerateUserSummary(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JSONレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
} 