package structs

import "github.com/meddler-io/watchdog/bootstrap"

type TaskResult struct {
	bootstrap.TaskResult
	Response string `json:"response" ` // success_endpoint
}
