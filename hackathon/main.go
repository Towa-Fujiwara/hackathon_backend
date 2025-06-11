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
	"os"
	"github.com/go-chi/chi"
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

	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")
	dao.InitDB(dbUser, dbPwd, dbName, instanceConnectionName)
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
	r := chi.NewRouter()

	firebaseAuthMiddleware := controller.AuthMiddleware(authClient)

	r.Get("/api/users/", registerUserController.RegisterUserHandler) 
	r.Get("/api/search/", searchUserController.SearchUsersHandler)
	r.Get("/api/posts/", postController.GetAllPostsHandler)

	// 認証が必要なエンドポイント
	r.Group(func(r chi.Router) {	
		r.Use(firebaseAuthMiddleware)
		r.Post("/api/posts/", postController.CreatePostHandler)
		r.Post("/api/posts/{postId}/comments", commentController.CreateCommentHandler)
		r.Post("/api/users/{userId}/follow", followUserController.FollowUserHandler)
		r.Get("/api/users/{userId}/followers", followUserController.GetFollowersHandler)
		r.Get("/api/users/{userId}/following", followUserController.GetFollowingHandler)
		r.Get("/api/users/{userId}/is-following", followUserController.IsFollowingHandler)
	})

	// サーバーの起動、ポート番号は8080
	fmt.Println("http://localhost:8080 でサーバーを起動します")
	http.ListenAndServe(":8080", r)
}