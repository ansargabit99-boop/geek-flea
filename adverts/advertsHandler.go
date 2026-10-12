package adverts

import (
	"encoding/json"
	"fmt"
	"geekflea/auth"
	"net/http"
	"strconv"
	"strings"
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
	advert:=&advertThing{}
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

func deleteAdverts(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
    advertID, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || advertID <= 0 {
        http.Error(w, `{"message":"invalid advert ID"}`, http.StatusBadRequest)
        return
    }

    tag, err := pool.Exec(
        r.Context(),
        "DELETE FROM adverts WHERE id = $1",
        advertID,
    )

    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(errorRe{Message: "something went wrong"})
        return
    }

    if tag.RowsAffected() == 0 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(errorRe{Message: "advert not found"})
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "advert deleted successfully",
    })
}

func changeAdverts(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	type changes struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Price       *float64 `json:"price"`
	}

	advertId := r.PathValue("id")
	numAdvertId, err := strconv.Atoi(advertId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorRe{Message: "invalid id"})
		return
	}

	data := &changes{}

	err = json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorRe{Message: "invalid body"})
		return
	}

	query := "UPDATE adverts SET "
	var sets []string
	var values []any
	i := 1

	if data.Name != "" {
		sets = append(sets, fmt.Sprintf("name=$%d", i))
		values = append(values, data.Name)
		i++
	}

	if data.Description != "" {
		sets = append(sets, fmt.Sprintf("description=$%d", i))
		values = append(values, data.Description)
		i++
	}

	if data.Price != nil {
		sets = append(sets, fmt.Sprintf("price=$%d", i))
		values = append(values, *data.Price)
		i++
	}

	if len(values) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorRe{Message: "nothing to change"})
		return
	}

	query += strings.Join(sets, ", ")
	query += fmt.Sprintf(" WHERE id=$%d", i)
	values = append(values, numAdvertId)

	tag, err := pool.Exec(r.Context(), query, values...)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorRe{Message: "something went wrong"})
		return
	}

	if tag.RowsAffected() == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorRe{Message: "advert not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(errorRe{Message: "successfully changed"})
}
