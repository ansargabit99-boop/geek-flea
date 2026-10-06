package auth

import (
	"encoding/json"
	
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)
type errorResponse struct {
	message string
}
func generateToken(id int) (string,error){
	claims := jwt.MapClaims{
		"id":id,
		"exp": time.Now().Add(30*24*time.Hour).Unix(),
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodES256,claims)

	secret := os.Getenv("JWT_SECRET")
	return  token.SignedString([]byte(secret))
}
func errorMessageHandler(w http.ResponseWriter,head int,message string) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(head)
	json.NewEncoder(w).Encode(errorResponse{message:message})
}
func Register(w http.ResponseWriter,r *http.Request,pool *pgxpool.Pool) {
	type requestedData struct {
		Name string `json:"name"`
		Gmail string `json:"gmail"`
		Number string `json:"number"`
		Password string `json:"password"`
	}
	data:= requestedData{}
	err:= json.NewDecoder(r.Body).Decode(&data)
	if data.Name == "" || data.Gmail == "" || data.Number == "" ||data.Password==""{
		errorMessageHandler(w,http.StatusBadRequest,"please fill in the form corretly")
		return
	}
	if len(data.Password) < 8 {
		errorMessageHandler(w,http.StatusBadRequest,"password should contain at least 8 characters")
		return
	}
	if err!=nil {
		errorMessageHandler(w,http.StatusBadRequest,"pleas fill the fields correctly")
		return
	}
	var id int
	hashed,err := bcrypt.GenerateFromPassword([]byte(data.Password),bcrypt.DefaultCost)
	if err != nil {
		errorMessageHandler(w,http.StatusInternalServerError,"something wentt wrong")
		return
	}

	err = pool.QueryRow(r.Context(),"INSERT INTO users (name,gmail,number,password) VALUES ($1,$2,$3,$4) RETURNING id",data.Name,data.Gmail,data.Number,string(hashed)).Scan(&id)
	if err != nil {
		errorMessageHandler(w,http.StatusInternalServerError,"something went wrong")
		return
	}
	token,err:= generateToken(id)
	if err != nil {
		errorMessageHandler(w,http.StatusInternalServerError,"something went wrong")
		return
	}
	type responseData struct {
		Token string `json:"token"`
		Id int `json:"id"`
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData{Token:token,Id:id})

}
func Login(w http.ResponseWriter,r*http.Request,pool *pgxpool.Pool) {
	type helkData struct {
		Token string `json:"token"`
		Id int `json:"id"`
	}
	type responseData struct {
		GmailOrNum string `json:"gmailOrNum"`
		Password string `json:"password"`
	}
	var userCredentials = responseData{}
	json.NewDecoder(r.Body).Decode(&userCredentials)
	var id int
	var password string
	if strings.Contains(userCredentials.GmailOrNum,"@") {
		err:= pool.QueryRow(r.Context(),"SELECT id FROM users WHERE gmail=$1",userCredentials.GmailOrNum).Scan(&id)
		if err!=nil {
			errorMessageHandler(w,http.StatusInternalServerError,"something went wron")
			return
		}
	} else {
		err:= pool.QueryRow(r.Context(),"SELECT id,passwordFROM users WHERE number=$1",userCredentials.GmailOrNum).Scan(&id,&password)
		if err != nil {
			errorMessageHandler(w,http.StatusInternalServerError,"something went wrong")
			return
		}
	}
	err := bcrypt.CompareHashAndPassword(
		[]byte(userCredentials.Password),
		[]byte(password),
	)
	if err!=nil {
		errorMessageHandler(w,http.StatusForbidden,"incorrect credentials")
		return
	}
	token,err := generateToken(id)
	if err !=nil {
		errorMessageHandler(w,http.StatusInternalServerError,"something went wrong")
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(helkData{Token:token,Id:id})

}
func GetSelfInformation(w http.ResponseWriter,r *http.Request,pool *pgxpool.Pool){
	userId,ok :=r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w,"unauthorised",http.StatusUnauthorized)
		return
	}
	type userData struct {
		Id int `json:"id"`
		Name string `json:"name"`
		Gmail string `json:"gmail"`
		Number string `json:"number"`
	}
	var user = userData{}
	err := pool.QueryRow(r.Context(),"SELECT id,name,gmail,number FROM users WHERE id=$1",userId).Scan(&user.Id,&user.Name,&user.Gmail,&user.Number)
	if err != nil {
		errorMessageHandler(w,http.StatusInternalServerError,"something went wrong")
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}