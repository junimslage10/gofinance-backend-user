package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"

	"github.com/junimslage10/gofinance-backend-user/api"
	db "github.com/junimslage10/gofinance-backend-user/db/sqlc"
	env "github.com/junimslage10/gofinance-backend-user/util"
)

func main() {
	var dbDriver, filledStringDbSource, serverAddress = env.LoadEnv()

	conn, err := sql.Open(dbDriver, filledStringDbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start api: ", err)
	}

}
