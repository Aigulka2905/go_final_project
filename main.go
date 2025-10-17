package main

import (
	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
	"log"
)

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatal(err)
	}
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
