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
	"github.com/go-chi/chi/v5"
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
	
	projectID := os.Getenv("PROJECT_ID")
	location := os.Getenv("LOCATION")

	postDao := dao.NewPostDao(dao.DB())
	registerUserDao := dao.NewUserDao(dao.DB())
	searchUserDao := dao.NewUserDao(dao.DB())
	commentDao := dao.NewCommentDao(dao.DB())
	followUserDao := dao.NewFollowUserDao(dao.DB())
	postLikeDao := dao.NewLikeDao(dao.DB())
	var geminiUsecase *usecase.GeminiUsecase
	if projectID != "" && location != "" {
		geminiUsecase, err = usecase.NewGeminiUsecase(postDao, projectID, location)
		if err != nil {
			log.Printf("Warning: failed to initialize Gemini usecase: %v", err)
		}
	}

	postUsecase := usecase.NewPostUsecase(postDao)
	registerUserUsecase := usecase.NewUserUsecase(registerUserDao)
	searchUserUsecase := usecase.NewUserUsecase(searchUserDao)
	commentUsecase := usecase.NewCommentUsecase(commentDao)
	followUserUsecase := usecase.NewFollowUserUsecase(followUserDao)
	postLikeUsecase := usecase.NewPostLikeUsecase(postLikeDao)
	userUsecase := usecase.NewUserUsecase(registerUserDao)
	
	postController := controller.NewPostController(postUsecase, registerUserUsecase)
	registerUserController := controller.NewRegisterUserController(registerUserUsecase)
	searchUserController := controller.NewSearchUserController(searchUserUsecase)
	commentController := controller.NewPostCommentController(commentUsecase, userUsecase)
	followUserController := controller.NewFollowUserController(followUserUsecase)
	postLikeController := controller.NewPostLikeController(postLikeUsecase, userUsecase)
	
	// Geminiコントローラーの初期化
	var geminiController *controller.GeminiController
	if geminiUsecase != nil {
		geminiController = controller.NewGeminiController(geminiUsecase)
	}

	// ルーティングの設定
	r := chi.NewRouter()

	r.Use(controller.CorsMiddleware)
    firebaseAuthMiddleware := controller.AuthMiddleware(authClient)

	r.Get("/api/search", searchUserController.SearchUsersHandler)
	r.Get("/api/posts", postController.GetAllPostsHandler)
	r.Get("/api/posts/{postId}", postController.GetPostHandler)
	r.Get("/api/users/id/{userId}", searchUserController.GetUserProfileHandler)  

	r.Group(func(r chi.Router) {	
		r.Use(firebaseAuthMiddleware)
		r.Post("/api/users", registerUserController.RegisterUserHandler) 
		r.Get("/api/users/me", searchUserController.GetUserProfileHandler)     
		 //r.Put("/api/users/me", registerUserController.UpdateUserHandler) 
		r.Get("/api/posts/me", postController.GetAllPostsByUserIdHandler) 
		r.Post("/api/posts", postController.CreatePostHandler)
		r.Post("/api/posts/{postId}/like", postLikeController.LikePostHandler)
		r.Get("/api/posts/{postId}/comments", commentController.GetCommentsHandler)
		r.Post("/api/posts/{postId}/comments", commentController.CreateCommentHandler)
		r.Post("/api/users/{userId}/follow", followUserController.FollowUserHandler)
		r.Get("/api/users/{userId}/followers", followUserController.GetFollowersHandler)
		r.Get("/api/users/{userId}/following", followUserController.GetFollowingHandler)
		r.Get("/api/users/{userId}/is-following", followUserController.IsFollowingHandler)
		
		// Gemini関連のエンドポイント
		if geminiController != nil {
			r.Post("/api/users/me/summary", geminiController.GenerateMySummaryHandler)
		}
	})
	
	// 認証不要のGeminiエンドポイント
	if geminiController != nil {
		r.Get("/api/users/{userId}/summary", geminiController.GenerateUserSummaryHandler)
	}

	http.ListenAndServe(":8080", r)
}