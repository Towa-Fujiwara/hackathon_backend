package controller

import "net/http"

func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        allowedOrigins := []string{
            "http://localhost:5173",
            "http://localhost:3000",
			"https://hackathon-frontend-one-khaki.vercel.app",
			"https://hackathon-frontend-4ws2075hc-towa-fujiwaras-projects.vercel.app",
        }
        
        origin := r.Header.Get("Origin")
        for _, allowedOrigin := range allowedOrigins {
            if origin == allowedOrigin {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                break
            }
        }
        
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
