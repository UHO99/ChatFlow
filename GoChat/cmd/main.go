package main

import (
	"GoChat/api"
	"GoChat/config"
	"GoChat/db/sqlc"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot Load Config : ", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Cannot create database connection : ", err)
	}

	store := sqlc.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server : ", err)
	}

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("Cannot start server : ", err)
	}
}
