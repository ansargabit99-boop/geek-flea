package main

import (
	"geekflea/auth"
	data "geekflea/database"
	"geekflea/users"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err !=nil {
		log.Fatal("something went wrong ")
	}
	pool,err := data.ConnectDb(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("something faiked")
	}
	
	defer pool.Close()
	mux:=http.NewServeMux()
	auth.AuthHandlers(mux,pool)
	port:= os.Getenv("PORT")
	users.UsersHandler(mux,pool)
	log.Fatal(http.ListenAndServe(":"+port,mux))

}