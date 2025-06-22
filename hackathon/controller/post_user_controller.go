package controller

import (
	"net/http"
	"hackathon/model"
	"hackathon/usecase"
	"hackathon/dao"
	"github.com/go-chi/chi/v5"
	"log"
)

type PostController struct {
	postUsecase usecase.PostUsecase
	userUsecase usecase.UserUsecase
	followUserDao dao.FollowUserDao
}

func NewPostController(pu usecase.PostUsecase, uu usecase.UserUsecase, fud dao.FollowUserDao) *PostController {
	return &PostController{
		postUsecase: pu,
		userUsecase: uu,
		followUserDao: fud,
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

// フォロー中のユーザーの投稿を取得するハンドラー
func (c *PostController) GetFollowingPostsHandler(w http.ResponseWriter, r *http.Request) {
	// 認証されたユーザーのUIDを取得
	firebaseUID, ok := r.Context().Value(userContextKey).(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "User ID not found in context. This endpoint requires authentication.", http.StatusUnauthorized)
		return
	}

	// 認証されたユーザーの情報を取得
	appUser, err := c.userUsecase.GetUserByFirebaseUID(firebaseUID)
	if err != nil || appUser == nil {
		log.Printf("ERROR: Failed to find user by firebaseUID %s: %v\n", firebaseUID, err)
		http.Error(w, "Authenticated user not found in application database.", http.StatusInternalServerError)
		return
	}

	// フォローしているユーザーのリストを取得
	following, err := c.followUserDao.GetFollowing(appUser.UserId)
	if err != nil {
		log.Printf("ERROR: Failed to get following users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// フォローしているユーザーがいない場合は空の配列を返す
	if len(following) == 0 {
		respondJSON(w, http.StatusOK, []model.Post{})
		return
	}

	// フォローしているユーザーのIDリストを作成
	followingUserIds := make([]string, len(following))
	for i, follow := range following {
		followingUserIds[i] = follow.FollowUserId
	}

	// フォローしているユーザーの投稿のみを取得
	var posts []model.Post
	allPosts, err := c.postUsecase.FindAllPosts()
	if err != nil {
		log.Printf("ERROR: FindAllPosts failed: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// フォローしているユーザーの投稿のみをフィルタリング
	for _, post := range allPosts {
		for _, followingUserId := range followingUserIds {
			if post.UserId == followingUserId {
				posts = append(posts, post)
				break
			}
		}
	}

	log.Printf("Found %d posts from %d following users", len(posts), len(followingUserIds))
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

// 他のユーザーの投稿を取得するハンドラー
func (c *PostController) GetPostsByUserIdHandler(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	posts, err := c.postUsecase.FindAllPostsByUserId(userId)
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

func (c *PostController) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// 認証されたユーザーのUIDを取得
	firebaseUID, ok := r.Context().Value(userContextKey).(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "User ID not found in context. This endpoint requires authentication.", http.StatusUnauthorized)
		return
	}

	// 投稿IDを取得
	postId := chi.URLParam(r, "postId")
	if postId == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	// 投稿を取得して所有者を確認
	post, err := c.postUsecase.FindById(postId)
	if err != nil {
		log.Printf("ERROR: Failed to find post: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// 認証されたユーザーの情報を取得
	appUser, err := c.userUsecase.GetUserByFirebaseUID(firebaseUID)
	if err != nil || appUser == nil {
		log.Printf("ERROR: Failed to find user by firebaseUID %s: %v\n", firebaseUID, err)
		http.Error(w, "Authenticated user not found in application database.", http.StatusInternalServerError)
		return
	}

	// 投稿の所有者かどうかを確認
	if post.UserId != appUser.UserId {
		http.Error(w, "You can only delete your own posts", http.StatusForbidden)
		return
	}

	// 投稿を削除
	err = c.postUsecase.DeletePost(postId)
	if err != nil {
		log.Printf("ERROR: Failed to delete post: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}