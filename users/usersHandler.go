package users

import (
	"net/http"
	"geekflea/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UsersHandler(mux *http.ServeMux,pool *pgxpool.Pool) {
	mux.Handle("/GET",auth.Middleware(http.HandlerFunc( func (w http.ResponseWriter, r * http.Request)  {
		getUsers(w,r,pool)
	})))
	mux.Handle("GET /users/{id}",auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getUser(w,r,pool)
	})))
}