package controller

import (
	"context"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)

// userContextKey は、リクエストコンテキスト内でFirebase UIDを安全に受け渡すための一意なキーです。
type contextKey string
const userContextKey = contextKey("firebase_uid")

// AuthMiddleware は、リクエストの Authorization ヘッダーからIDトークンを検証するミドルウェアを生成します。
//
// 使い方:
//  authClient, _ := app.Auth(ctx)
//  firebaseAuthMiddleware := AuthMiddleware(authClient)
//  http.Handle("/your/protected/route", firebaseAuthMiddleware(yourHandler))
func AuthMiddleware(authClient *auth.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// http.HandlerFunc は、関数をhttp.Handlerに変換するアダプターです。
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// リクエストヘッダーから "Authorization" の値を取得します。
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// ヘッダーが "Bearer <token>" の形式であることを確認します。
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format. Expected 'Bearer <token>'", http.StatusUnauthorized)
				return
			}
			idToken := parts[1]

			// Firebase Admin SDK を使ってIDトークンを検証します。
			token, err := authClient.VerifyIDToken(r.Context(), idToken)
			if err != nil {
				// トークンが無効な場合（期限切れ、不正な形式など）
				log.Printf("error verifying ID token: %v\n", err)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// 検証成功後、トークンからUID（ユーザーID）を取得します。
			// このUIDをリクエストコンテキストに保存します。
			ctx := context.WithValue(r.Context(), userContextKey, token.UID)

			// 元のリクエスト(r)の代わりに、UIDを含む新しいコンテキストを持つリクエストを
			// 次のハンドラに渡します。
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}