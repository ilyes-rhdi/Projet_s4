package middleware

import (
	"Devenir_dev/internal/api/models"
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" || r.URL.Path == "/submit" {
			next.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Injecter les claims dans le contexte de la requête
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Récupère les claims depuis le contexte
		claims, ok := r.Context().Value("user").(jwt.MapClaims)
		if !ok {
			http.Error(w, "Erreur de récupération des claims", http.StatusInternalServerError)
			return
		}

		// Vérifie si l'utilisateur est un admin
		role, ok := claims["role"].(string)
		if !ok || role != string(models.Admin) {
			http.Error(w, "Accès refusé : Vous n'êtes pas autorisé", http.StatusForbidden)
			return
		}

		// Si l'utilisateur est admin, appelle le handler suivant
		next.ServeHTTP(w, r)
	})
}

func IsResponsable(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Récupère les claims depuis le contexte
		claims, ok := r.Context().Value("user").(jwt.MapClaims)
		if !ok {
			http.Error(w, "Erreur de récupération des claims", http.StatusInternalServerError)
			return
		}

		// Vérifie si l'utilisateur est un responsable
		role, ok := claims["role"].(string)
		if !ok || role != string(models.Responsable) {
			http.Error(w, "Accès refusé : Vous n'êtes pas autorisé", http.StatusForbidden)
			return
		}

		// Si l'utilisateur est responsable, appelle le handler suivant
		next.ServeHTTP(w, r)
	})
}
