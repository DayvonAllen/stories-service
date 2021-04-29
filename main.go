package main

import (
	"example.com/app/database"
	"example.com/app/events"
	"example.com/app/router"
	"log"
)

func init() {
	// create database connection instance for first time
	_ = database.GetInstance()
	//err := repo.TagRepoImpl{}.CreateMany(util.GenerateTags())
	//if err != nil {
	//	return
	//}
	go events.ConnectToKafkaAsConsumer()
}

func main() {
	app := router.Setup()
	log.Fatal(app.Listen(":8081"))
}
