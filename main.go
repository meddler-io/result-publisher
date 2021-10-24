package main

import (
	"encoding/json"
	"os"

	"github.com/meddler-io/watchdog/bootstrap"
	"github.com/meddler-io/watchdog/consumer"
	"github.com/meddler-io/watchdog/logger"
	"meddler.io/result-publisher/db"
)

func getenvStr(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func main() {

	forever := make(chan bool)

	username := getenvStr("RMQ_USERNAME", "user")
	password := getenvStr("RMQ_PASSWORD", "bitnami")
	host := getenvStr("RMQ_HOST", "192.168.29.9")
	logger.Println("username", username)
	logger.Println("password", password)
	logger.Println("host", host)
	logger.Println("MESSAGEQUEUE", bootstrap.CONSTANTS.Reserved.MESSAGEQUEUE)
	// password := getenvStr("PORt", "bitnami")

	queue := consumer.NewQueue("amqp://"+username+":"+password+"@"+host, bootstrap.CONSTANTS.Reserved.MESSAGEQUEUE)
	defer queue.Close()
	// db.Test()

	queue.Consume(func(i string) {
		logger.Println("Received message with first consumer: %s", i)

		taskResult := bootstrap.TaskResult{}
		err := json.Unmarshal([]byte(i), &taskResult)
		logger.Println(i)
		logger.Println(taskResult)
		logger.Println(err)

		db.Test(taskResult)
	})

	<-forever

}
