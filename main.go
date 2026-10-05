package main

import (
	data "geekflea/database"
	"log"
	"net/http"
	"geekflea/auth"
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
	log.Fatal(http.ListenAndServe(":"+port,mux))

}