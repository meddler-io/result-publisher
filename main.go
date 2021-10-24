package main

import (
	"encoding/json"
	"log"

	"github.com/meddler-io/watchdog/bootstrap"
	"github.com/meddler-io/watchdog/consumer"
	"github.com/meddler-io/watchdog/logger"
	"meddler.io/result-publisher/db"
	"meddler.io/result-publisher/helper"

	"meddler.io/result-publisher/structs"
)

func main() {

	forever := make(chan bool)

	username := helper.GetenvStr("RMQ_USERNAME", "user")
	password := helper.GetenvStr("RMQ_PASSWORD", "bitnami")
	host := helper.GetenvStr("RMQ_HOST", "192.168.29.9")
	logger.Println("username", username)
	logger.Println("password", password)
	logger.Println("host", host)
	logger.Println("MESSAGEQUEUE", bootstrap.CONSTANTS.Reserved.MESSAGEQUEUE)
	// password := getenvStr("PORt", "bitnami")

	queue := consumer.NewQueue("amqp://"+username+":"+password+"@"+host, bootstrap.CONSTANTS.Reserved.MESSAGEQUEUE)
	defer queue.Close()
	// db.Test()

	queue.Consume(func(i string) {

		taskResult := structs.TaskResult{}
		err := json.Unmarshal([]byte(i), &taskResult)
		if err != nil {
			log.Println("Invalid Message", i)
			return
		}
		logger.Println("Task:", taskResult.Identifier, taskResult.Status, taskResult.Message)
		db.UpdateTaslResult(taskResult)
	})

	<-forever

}
