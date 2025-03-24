package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/Te8va/shortURL/internal/app/domain"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(secretKey string) (string, int, error) {
	id, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", 0, fmt.Errorf("GenerateToken: %w", err)
	}

	claims := Claims{
		UserID: int(id.Int64()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", 0, fmt.Errorf("GenerateToken: %w", err)
	}

	return tokenString, claims.UserID, nil
}

func GetUserIDFromCookie(r *http.Request, secretKey string) (int, bool) {
	cookie, err := r.Cookie("auth")
	if err != nil {
		return 0, false
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Println("JWT parsing error:", err)
		return 0, false
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, false
	}

	return claims.UserID, true
}

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, valid := GetUserIDFromCookie(r, secretKey)

			if !valid {
				tokenString, newUserID, err := GenerateToken(secretKey)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:  "auth",
					Value: tokenString,
					Path:  "/",
				})

				userID = newUserID
			}

			ctx := context.WithValue(r.Context(), domain.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
