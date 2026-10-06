package users

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func usersHandler(mux *http.ServeMux,pool *pgxpool.Pool) {
	mux.HandleFunc("/GET",func (w http.ResponseWriter, r * http.Request)  {
		getUsers(w,r,pool)
	})
}