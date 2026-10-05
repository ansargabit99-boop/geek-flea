package data

import (
	"context"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDb(db_url string) (*pgxpool.Pool,error){
	ctx:= context.Background()
	pool,err := pgxpool.New(ctx,db_url)
	if err!= nil {
		log.Fatal("failed to connect to database")
	}
	err = pool.Ping(ctx)
	if err != nil {
		return pool,err
	}
	return pool,nil
}