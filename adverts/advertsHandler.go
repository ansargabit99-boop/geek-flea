package adverts

import (
	"encoding/json"
	"net/http"
	"strconv"
	"geekflea/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)
type  errorRe struct {
			Message string
		}
type advertsSkelet struct {
	Id int `json:"id"`
	User_id int `json:"user_id"`
	Name string `JSON:"name"`
	Desc string `JSON:"desc"`
	Price int `JSON:"price"`
}
func getAdverts(w http.ResponseWriter,r *http.Request,pool *pgxpool.Pool) {
	type advertsStruct struct {
		Id int `json:"id"`
		User_id int `json:"user_id"`
		Name string `json:"name"`
		Desc string `json:"desc"`
		Price int `json:"price"`
	}
	var adverts []advertsStruct
	rows,err := pool.Query(r.Context(),"SELECT * FROM adverts")
	if err != nil {
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorRe{Message:"something wetn wrong"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u advertsStruct
		err:= rows.Scan(&u.Id,&u.User_id,&u.Name,&u.Desc,&u.Price)

		if err != nil {
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorRe{Message:"something went wrong"})
			return
		}
		adverts = append(adverts, u)
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(adverts)
}
func getAdvert(w http.ResponseWriter,r *http.Request,pool *pgxpool.Pool) {
	var id = r.PathValue("id")
	idNum,err := strconv.Atoi(id)
	if err != nil {
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorRe{Message: "invalid id"})
		return
	}
	var advert advertsSkelet
	err = pool.QueryRow(r.Context(),"SELECT * FROM advets WHERE id=$1",idNum).Scan(&advert.Id,&advert.Name,&advert.User_id,&advert.Desc,&advert.Price)
	if err != nil {
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorRe{Message: "something went wrong"})
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(advert)
}
func postAdverts(w http.ResponseWriter, r *http.Request,pool *pgxpool.Pool) {
	userId:=r.Context().Value(auth.UserIDKey)
	type advertThing struct {
		Name string `json:"Name"`
		Desc string `json:"Desc"`
		Price int `json:"price"`
	}
	var advert advertThing
	err:= json.NewDecoder(r.Body).Decode(advert)
	if err != nil {
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorRe{Message: "something went wrong"})
		return
	}
	var resData advertsSkelet
	err =  pool.QueryRow(r.Context(),"INSERT INTO adverts (name,user_id,description,price)",advert.Name,userId,advert.Desc,advert.Price).Scan(&resData.Id,&resData.User_id,&resData.Name,&resData.Desc,&resData.Price)
	if err != nil {
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(errorRe{Message: "something went wrong"})
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resData)

}