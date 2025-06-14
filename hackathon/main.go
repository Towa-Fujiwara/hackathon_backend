package main

import (
	"net/http"
	"hackathon/controller"
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
			userId TEXT,
			text TEXT,
			image TEXT,
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`

	insertPost = "INSERT INTO posts (userId, text, image, createdAt) VALUES (?, ?, ?, ?)"


	selectPosts = "SELECT * FROM posts ORDER BY createdAt DESC"
)

func main() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("serviceAccountKey.json")
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

	r.Use(controller.CorsMiddleware)
    firebaseAuthMiddleware := controller.AuthMiddleware(authClient)

	r.Get("/api/search", searchUserController.SearchUsersHandler)
	r.Get("/api/posts", postController.GetAllPostsHandler)


	r.Group(func(r chi.Router) {	
		r.Use(firebaseAuthMiddleware)
		r.Post("/api/users", registerUserController.RegisterUserHandler) 
		r.Get("/api/users/me", searchUserController.GetUserProfileHandler)     
		 //r.Put("/api/users/me", registerUserController.UpdateUserHandler)      
		r.Get("/api/posts/me", postController.GetAllPostsByUserIdHandler) 
		r.Post("/api/posts", postController.CreatePostHandler)
		r.Post("/api/posts/{postId}/comments", commentController.CreateCommentHandler)
		r.Post("/api/users/{userId}/follow", followUserController.FollowUserHandler)
		r.Get("/api/users/{userId}/followers", followUserController.GetFollowersHandler)
		r.Get("/api/users/{userId}/following", followUserController.GetFollowingHandler)
		r.Get("/api/users/{userId}/is-following", followUserController.IsFollowingHandler)
	})


	http.ListenAndServe(":8080", r)
}