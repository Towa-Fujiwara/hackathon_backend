package main

import (
	"net/http"
	"hackathon/controller"
	"fmt"
	"hackathon/dao"
	"hackathon/usecase"
)



const (
	createPostTable = `
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT,
			text TEXT,
			image TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`

	// 投稿の作成を行うSQL文
	insertPost = "INSERT INTO posts (user_id, text, image, created_at) VALUES (?, ?, ?, ?)"

	// 投稿の取得を行うSQL文
	selectPosts = "SELECT * FROM posts ORDER BY created_at DESC"
)

// main関数は、プログラムのエントリーポイント、init()関数の実行後に実行される
func main() {
	dao.InitDB()
	postDao := dao.NewPostDao(dao.DB())
	postUsecase := usecase.NewPostUsecase(postDao)
	postController := controller.NewPostController(postUsecase)
	registerUserDao := dao.NewUserDao(dao.DB())
	registerUserUsecase := usecase.NewUserUsecase(registerUserDao)
	registerUserController := controller.NewRegisterUserController(registerUserUsecase)
	searchUserDao := dao.NewUserDao(dao.DB())
	searchUserUsecase := usecase.NewUserUsecase(searchUserDao)
	searchUserController := controller.NewSearchUserController(searchUserUsecase)

	// ルーティングの設定
	http.HandleFunc("/api/users/", registerUserController.RegisterUserHandler)
	http.HandleFunc("/api/search/", searchUserController.SearchUsersHandler)
	http.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			postController.GetAllPostsHandler(w, r)
		case http.MethodPost:
			postController.CreatePostHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// サーバーの起動、ポート番号は8080
	fmt.Println("http://localhost:8080 でサーバーを起動します")
	http.ListenAndServe(":8080", nil)
}