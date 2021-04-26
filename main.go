package main

import (
	"example.com/app/database"
	"example.com/app/events"
)

func init() {
	// create database connection instance for first time
	_ = database.GetInstance()
	events.ConnectToKafkaAsConsumer()
}

func main() {

}
