package main

import (
	"example.com/app/database"
	"example.com/app/events"
	"example.com/app/repo"
	"example.com/app/router"
	"example.com/app/util"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func init() {
	// create database connection instance for first time
	_ = database.GetInstance()
	err := repo.TagRepoImpl{}.CreateMany(util.GenerateTags())
	if err != nil {
		return
	}
	go events.ConnectToKafkaAsConsumer()
}

func main() {
	app := router.Setup()

	// graceful shutdown on signal interrupts
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		_ = <- c
		fmt.Println("Shutting down...")
		_ = app.Shutdown()
	}()

	if err := app.Listen(":8081"); err != nil {
		log.Panic(err)
	}
}
