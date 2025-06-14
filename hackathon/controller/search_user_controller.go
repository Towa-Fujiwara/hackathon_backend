package controller

import (
	"net/http"
	"hackathon/usecase"
	"log"
)

type SearchUserController struct {
	searchUserUsecase usecase.UserUsecase
}

func NewSearchUserController(su usecase.UserUsecase) *SearchUserController {
	return &SearchUserController{searchUserUsecase: su}
}

func (c *SearchUserController) GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(userContextKey).(string)
	if !ok || uid == "" {
		http.Error(w, "User ID not found in context. This endpoint requires authentication.", http.StatusInternalServerError)
		return
	}
	user, err := c.searchUserUsecase.SearchUserExist(uid)
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






















	/*parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		respondJSON(w, http.StatusBadRequest, "User ID is missing")
		return
	}
	userID := parts[3]

	var responseData model.ProfilePageResponse

	// 1. ユーザープロフィールを取得
	var userProfile model.PublicUser
	err := db.QueryRow(selectPublicUserByID, userID).Scan(
		&userProfile.Id,
		&userProfile.Name,
		&userProfile.Profile.IconUrl,
		&userProfile.Profile.DisplayName,
		&userProfile.Profile.Bio,
		&userProfile.Profile.BackgroundImageUrl,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("User not found for ID: %s", userID)
			respondJSON(w, http.StatusNotFound, "User not found")
			return
		}
		log.Printf("DB error fetching profile for user %s: %v", userID, err)
		respondJSON(w, http.StatusInternalServerError, "Database error on fetching profile")
		return
	}
	responseData.UserProfile = userProfile

	// 2. ユーザーの投稿一覧を取得 (ロジックは前回から変更なし)
	rows, err := db.Query(selectPostsWithCountsByUserID, userID)
	if err != nil {
		log.Printf("DB error fetching posts for user %s: %v", userID, err)
		respondJSON(w, http.StatusInternalServerError, "Database error on fetching posts")
		return
	}
	defer rows.Close()

	var userPosts []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.Id, &post.UserId, &post.Text, &post.Image, &post.CreatedAt,
			&post.LikeCount, &post.CommentCount,
		)
		if err != nil {
			log.Printf("Error scanning post row: %v", err)
			continue
		}
		userPosts = append(userPosts, post)
	}
	responseData.UserPosts = userPosts

	// 3. 取得したデータをまとめてJSONでレスポンス
	respondJSON(w, http.StatusOK, responseData)
}*/

