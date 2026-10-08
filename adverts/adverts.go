package adverts

import (
	"net/http"
	"geekflea/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)

func advertsHandler(mux *http.ServeMux,pool *pgxpool.Pool) {
	mux.Handle("GET /adverts",auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	    getAdverts(w,r,pool)
	})))
	mux.Handle("GET /adverts/{id}",auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getAdvert(w,r,pool)
	})))
}	