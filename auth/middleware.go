package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken(tokenStr string) (jwt.MapClaims, error) { //what does things also do here
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) { //what dpes the (t) (Anby,err mesan here)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok { ////why _ here and di ont understand your structure fo go ypui know
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	return claims, nil

}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "no token", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "bad token", http.StatusUnauthorized)
			return
		}
		fmt.Println("userId", claims["id"])
		next.ServeHTTP(w, r)
	})
}
