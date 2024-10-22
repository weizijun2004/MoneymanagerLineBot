package main

import (
	"log"

	"moneyLineBot/model/connect"
	"moneyLineBot/model/create"
)

func main() {
	err := connect.InitializeMongoDatabaseClient()
	if err != nil {
		log.Fatal(err)
	}
	create.CreateGroup("test", []string{"test1", "test2"})
	create.CreateEvent("test", []string{"test1", "test2"}, "testEvent")
	connect.CloseMongoDatabaseClient()
}
