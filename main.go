package main

import (
	"log"
)

func main() {
	store, err := NewPostgresStorage()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":9696", store)
	server.Run()

}
