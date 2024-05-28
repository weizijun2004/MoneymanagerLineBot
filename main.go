package main

import (
	"moneyLineBot/model/create"
)

func main() {
	create.CreateGroup("test", []string{"test1", "test2"})
	create.CreateEvent("test", []string{"test1", "test2"}, "testEvent")

}
