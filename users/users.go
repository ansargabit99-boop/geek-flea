package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func errorMessageHandler(w http.ResponseWriter,head int,message string) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(head)
	json.NewEncoder(w).Encode(map[string]string{message:message})
}
func getUsers(w http.ResponseWriter ,r *http.Request,pool *pgxpool.Pool) {
	type users struct {
		Id int `json:"oddd"`
		Name string `json:"name"`
		Gmail string `json:"gmail"`
		Number string `json:"number"`
	}
	var respondData  = []users{}
	rows,err := pool.Query(r.Context(),"SELECT id,name,gmail,number FROM users")
	if err != nil {
		errorMessageHandler(w,http.StatusInternalServerError,"something went wrong")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u users
		if err := rows.Scan(&u.Id,&u.Name,&u.Gmail,&u.Number);err!=nil {
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		respondData = append(respondData, u)

	}
	if err  = rows.Err();err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respondData)
}
func getUser(w http.ResponseWriter,r *http.Request,pool *pgxpool.Pool) {
	strId := r.PathValue("id")
	id,err := strconv.Atoi(strId)
	if err != nil {
		errorMessageHandler(w,http.StatusBadRequest,"invalid Id")
	}
	type user struct {
		Id int `json:"id"`
		Name string `json:"name"`
		Gmail string `json:"gmail"`
		Number string `json:"Number"`
	}
	var respondData = user{}
	err = pool.QueryRow(r.Context(),"SELECT * FROM users WHERE id=$1",id).Scan(&respondData.Id,&respondData.Name,&respondData.Gmail,&respondData.Number)
	if err != nil {
		errorMessageHandler(w, http.StatusInternalServerError,"something went wrong")
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respondData)
}

