package middleware

import (
	"context"
	"net/http"
	"strings"
	"Devenir_dev/internal/api/models"
    "os"
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
        claims, ok := r.Context().Value("user").(*models.Claims)
        if !ok {
            http.Error(w, "Erreur de récupération des claims", http.StatusInternalServerError)
            return
        }

        // Vérifie si l'utilisateur est un admin
        if claims.Role != "admin" {
            http.Error(w, "Accès refusé : Vous n'êtes pas autorisé", http.StatusForbidden)
            return
        }

        // Si l'utilisateur est admin, appelle le handler suivant
        next.ServeHTTP(w, r)
    })
}
func RequireRole(role models.Role, permission string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Retrieve the user's role from the context (assume it's injected by JWT middleware)
            userRole, ok := r.Context().Value("role").(models.Role)
            if !ok || userRole != role {
                http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
                return
            }

            // Check if the role has the required permission
            if !models.RolePermissions[userRole][permission] {
                http.Error(w, "Forbidden: Permission denied", http.StatusForbidden)
                return
            }

            // If all checks pass, proceed to the next handler
            next.ServeHTTP(w, r)
        })
    }
}