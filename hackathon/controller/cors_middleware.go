package controller

import (
	"net/http"
	"log"
	)

func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request Origin: %s", r.Header.Get("Origin"))
        log.Printf("Request Method: %s", r.Method)
        allowedOrigins := []string{
            "http://localhost:5173",
            "http://localhost:3000",
			"https://hackathon-frontend-one-khaki.vercel.app",
			"https://hackathon-frontend-dusky.vercel.app",
			"https://hackathon-frontend-cu4r-2fjz7p3tv-tfujis-projects.vercel.app",
			"https://hackathon-frontend-cu4r.vercel.app",
			"https://hackathon-frontend-o0ty9vtbr-towa-fujiwaras-projects.vercel.app",
			"https://hackathon-frontend-cu4r-git-firebase-tfujis-projects.vercel.app",
			"https://hackathon-frontend-cu4r-myxph6rbv-tfujis-projects.vercel.app",
        }
        
        origin := r.Header.Get("Origin")
        for _, allowedOrigin := range allowedOrigins {
            if origin == allowedOrigin {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                break
            }
        }
        
        w.Header().Set("Access-Control-Allow-Methods", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
