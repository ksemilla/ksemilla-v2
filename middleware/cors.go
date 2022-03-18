package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
	"github.com/ksemilla/ksemilla-v2/config"
)

func GetCorsHandler() func(http.Handler) http.Handler {

	config := config.Config()

	return cors.Handler(cors.Options{
		AllowedOrigins: config.CORS_ALLOWED_HOSTS,
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}
