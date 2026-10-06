package auth

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthHandlers(mux *http.ServeMux,pool *pgxpool.Pool) {
	mux.HandleFunc("POST /login",func (w http.ResponseWriter,r *http.Request)  {
		Login(w,r,pool)
	})
	mux.HandleFunc("POST /register",func(w http.ResponseWriter,r*http.Request){
		Register(w,r,pool)
	})
	mux.Handle("GET /me",Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GetSelfInformation(w,r,pool)
	})))
}