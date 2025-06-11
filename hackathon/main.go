package main

import (
	"net/http"
	"hackathon/controller"
	"fmt"
	"hackathon/dao"
	"hackathon/usecase"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"context"
	"log"
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
	ctx := context.Background()
	opt := option.WithCredentialsFile("hackathon/serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

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
	commentDao := dao.NewCommentDao(dao.DB())
	commentUsecase := usecase.NewCommentUsecase(commentDao)
	commentController := controller.NewPostCommentController(commentUsecase)
	followUserDao := dao.NewFollowUserDao(dao.DB())
	followUserUsecase := usecase.NewFollowUserUsecase(followUserDao)
	followUserController := controller.NewFollowUserController(followUserUsecase)

	// ルーティングの設定

	firebaseAuthMiddleware := controller.AuthMiddleware(authClient)

	http.HandleFunc("/api/users/", registerUserController.RegisterUserHandler) 
	http.HandleFunc("/api/search/", searchUserController.SearchUsersHandler)
	http.HandleFunc("GET /api/posts/", postController.GetAllPostsHandler)

	// 認証が必要なエンドポイント
	// ミドルウェアでラップするために http.Handle を使用します
	http.Handle("POST /api/posts/", firebaseAuthMiddleware(http.HandlerFunc(postController.CreatePostHandler)))
	http.Handle("POST /api/posts/{postId}/comments", firebaseAuthMiddleware(http.HandlerFunc(commentController.CreateCommentHandler)))
	http.Handle("POST /api/users/{userId}/follow", firebaseAuthMiddleware(http.HandlerFunc(followUserController.FollowUserHandler)))
	http.Handle("GET /api/users/{userId}/followers", firebaseAuthMiddleware(http.HandlerFunc(followUserController.GetFollowersHandler)))
	http.Handle("GET /api/users/{userId}/following", firebaseAuthMiddleware(http.HandlerFunc(followUserController.GetFollowingHandler)))
	http.Handle("GET /api/users/{userId}/is-following", firebaseAuthMiddleware(http.HandlerFunc(followUserController.IsFollowingHandler)))
	// サーバーの起動、ポート番号は8080
	fmt.Println("http://localhost:8080 でサーバーを起動します")
	http.ListenAndServe(":8080", nil)
}