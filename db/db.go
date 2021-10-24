package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/meddler-io/watchdog/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"meddler.io/result-publisher/helper"
	"meddler.io/result-publisher/structs"
)

func toDoc(v structs.TaskResult) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func UpdateTaslResult(taskResult structs.TaskResult) {

	/*
	   Connect to my cluster
	*/

	MONGO_HOST := helper.GetenvStr("MONGO_HOST", "localhost")
	MONGO_PORT := helper.GetenvStr("MONGO_PORT", "27017")
	MONGO_DB := helper.GetenvStr("MONGO_DB", "fastapi")

	MONGO_COLLECTION := "builds_executor"

	MONGO_URL := fmt.Sprintf("mongodb://%s:%s/%s", MONGO_HOST, MONGO_PORT, MONGO_DB)

	log.Println("MONGO_URL", MONGO_URL)

	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO_URL))
	if err != nil {
		log.Println("Mongo Connection Failed")
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	_id, err := primitive.ObjectIDFromHex(taskResult.Identifier)

	// bootstrap.CONSTANTS.System.INPUTDIR

	log.Println(_id, err)
	collection := client.Database(MONGO_DB).Collection(MONGO_COLLECTION)

	responseStruct := bson.M{

		"exec_status":      taskResult.Status,
		"message":          taskResult.Message,
		"watchdog_version": taskResult.WatchdogVersion,
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(taskResult.Response), &response)

	if err == nil {
		responseStruct["result"] = response
	}

	result, err := collection.UpdateByID(context.Background(),
		_id,
		bson.M{
			"$set": responseStruct,
		},
	)

	logger.Println("Task:", "Updated Results", taskResult.Identifier, result, err)

}
