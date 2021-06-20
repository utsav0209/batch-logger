package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func syncPayloadsToTargetEndpoint(payloads []Log, targetEndpoint string) {

	start := time.Now()

	postBody, err := json.Marshal(payloads) // Prepare logs for request

	if err != nil {
		// Exit the application if JSON can not be serialized
		panic("JSON could not be marshalled")
	}

	for count := 0; count < 3; count++ {
		res, err := http.Post(targetEndpoint, "application/json", bytes.NewBuffer(postBody))

		// If err is nil then request was successful
		if err == nil {
			log.Infof("Synced Batch of Size %d which returned status code %d and took %s", len(payloads), res.StatusCode, time.Since(start))
			break
		}

		// If syncing failed for 3 times in a row exit the application
		if count == 2 {
			panic(fmt.Sprintf("Log syncing failed with error: %s", err))
		}

		// Wait for 2 seconds before retrying
		time.Sleep(time.Second * time.Duration(2))
	}
}
