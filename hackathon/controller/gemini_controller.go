package controller

import (
	"encoding/json"
	"net/http"
	"hackathon/usecase"
)

type GeminiController struct {
	geminiUsecase *usecase.GeminiUsecase
	userUsecase   usecase.UserUsecase
}

func NewGeminiController(geminiUsecase *usecase.GeminiUsecase, userUsecase usecase.UserUsecase) *GeminiController {
	return &GeminiController{
		geminiUsecase: geminiUsecase,
		userUsecase:   userUsecase, 
	}
}

// ユーザーの投稿からサマリーを生成するエンドポイント（POST /api/users/{userId}/summary）
func (gc *GeminiController) GenerateUserSummaryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	firebaseUID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "userID not found in context", http.StatusUnauthorized)
		return
	}
	userProfile, err := gc.userUsecase.GetUserByFirebaseUID(firebaseUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if userProfile == nil {
		http.Error(w, "User profile not found for the authenticated user", http.StatusNotFound)
		return
	}
	dbUserId := userProfile.UserId 

	summary, err := gc.geminiUsecase.GenerateUserSummary(r.Context(), dbUserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

func (gc *GeminiController) GenerateMySummaryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
} 