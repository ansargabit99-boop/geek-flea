package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const UserIDKey contextKey = "userID"

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

func Middleware(next http.Handler) http.Handler {
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
		idfloat,ok := claims["id"].(float64)
		if !ok {
			http.Error(w,"something went wrong",http.StatusInternalServerError)
			return 
		}
		ctx:=context.WithValue(r.Context(),UserIDKey,int(idfloat))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func AdminMiddleware(next http.Handler,pool *pgxpool.Pool) http.Handler{
	return  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId,ok:= r.Context().Value(UserIDKey).(int)
		if !ok {
			http.Error(w,"Unauthorised",http.StatusUnauthorized)
			return 
		}
		type userRole struct{
			Role string `json:"role"`
		}
		var role userRole
		err := pool.QueryRow(r.Context(),"SELECT role FROM users WHERE id = $1",userId).Scan(&role.Role)
		if err != nil {
			http.Error(w,"something went wrong",http.StatusInternalServerError)
			return 
		}
		if role.Role != "admin" {
			http.Error(w,"dont have permission",http.StatusUnauthorized)
			return 
		}
		next.ServeHTTP(w,r)
	})
}
