package main

import (
	"jungle-proj/api"
	"jungle-proj/db"
	"log"
)

const addr = ":3000"

func main() {

	store := db.NewMemoryStorage()

	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("connot create server", err)
	}

	err = server.Start(addr)
	if err != nil {
		log.Fatal(err)
	}
}
