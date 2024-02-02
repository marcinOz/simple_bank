package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/techschool/simplebank/api"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	log.Printf("### config.DBSource:%v", config.DBSource)

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	log.Println("### DB cunnected successfully!")

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	log.Println("### Server Created successfully!")

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
	log.Println("### Server started successfully!")
}
