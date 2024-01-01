package main

import (
	"database/sql"
	"jungle-proj/api"
	sqlc "jungle-proj/db/sqlc"
	"jungle-proj/util"

	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := sqlc.NewStore(conn)

	// ctx, cancle := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer cancle()

	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("connot create server", err)
	}

	err = server.Start(config.Port)
	if err != nil {
		log.Fatal(err)
	}
}
