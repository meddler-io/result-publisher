package db

import (
	"context"
	"log"
	"time"

	"github.com/meddler-io/watchdog/bootstrap"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// type TaskResultBson struct {
// 	TaskResult
// 	Response string `json:"response" ` // success_endpoint
// }

func toDoc(v bootstrap.TaskResult) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func Test(taskResult bootstrap.TaskResult) {

	/*
	   Connect to my cluster
	*/
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/fastapi"))
	if err != nil {
		log.Fatal(err)
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
	collection := client.Database("fastapi").Collection("builds_executor")
	result, err := collection.UpdateByID(context.Background(),
		_id,
		bson.M{
			"$set": bson.M{
				"result":           taskResult.Response,
				"exec_status":      taskResult.Status,
				"message":          taskResult.Message,
				"watchdog_version": taskResult.WatchdogVersion,
			},
		},
	)

	log.Println("result", result.MatchedCount, result.ModifiedCount, result.UpsertedCount)
	log.Println("err", err)

	raw, err := collection.FindOne(context.Background(), bson.M{
		"_id": _id,
	}).DecodeBytes()

	log.Println("raw", raw)

}
