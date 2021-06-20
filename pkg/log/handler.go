package log

import (
	"batch-logger/pkg/utils"
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Log struct {
	UserId int     `json:"user_id"`
	Total  float64 `json:"total"`
	Title  string  `json:"title"`
	Meta   struct {
		Logins []struct {
			Time time.Time `json:"time"`
			Ip   string    `json:"ip"`
		} `json:"logins"`
		PhoneNumbers struct {
			Home   string `json:"home"`
			Mobile string `json:"mobile"`
		} `json:"phone_numbers"`
	} `json:"meta"`
	Completed bool `json:"completed"`
}

type Handler struct {
	BatchSize      int    // Maximum batch size
	BatchInterval  int    // Interval in seconds at which syncing should happen
	TargetEndpoint string // Target endpoint where JSON arrays should we eventually written
	payloads       []Log
}

func CreateHandler(batchSize, batchInterval int, targetEndpoint string) Handler {
	return Handler{
		BatchSize:      batchSize,
		BatchInterval:  batchInterval,
		TargetEndpoint: targetEndpoint,
	}
}

func (h *Handler) AddPayload(payload Log) []Log {
	h.payloads = append(h.payloads, payload)

	// Check if batch size is exceeded or not
	if len(h.payloads) >= h.BatchSize {
		// Call sync payloads concurrently
		go h.syncPayloadsToPostEndpoint()
	}

	return h.payloads
}

// SyncAtIntervals Sync Payload Data in intervals to the provided endpoint
func (h *Handler) SyncAtIntervals() {
	for range time.Tick(time.Second * time.Duration(h.BatchInterval)) {
		if len(h.payloads) > 0 {
			// Call sync payloads concurrently
			go h.syncPayloadsToPostEndpoint()
		}
	}
}

func (h *Handler) HandleLogRequest(w http.ResponseWriter, r *http.Request) {
	// Fetch Payload from context and store in memory
	ctx := r.Context()
	payload := ctx.Value(utils.GetPayloadRequestId(ctx)).(Log)
	h.AddPayload(payload)

	// Write response
	w.WriteHeader(200)
	_, err := w.Write([]byte("Successfully Logged"))
	if err != nil {
		log.Errorf("Response failed with error: %s", err)
		return
	}
}

func (h *Handler) syncPayloadsToPostEndpoint() {
	log.Info("Syncing logs to POST Endpoint")

	payloadsLength := len(h.payloads)                           // Number of elements in slice
	postBody, err := json.Marshal(h.payloads[0:payloadsLength]) // Prepare data for request
	h.payloads = h.payloads[payloadsLength:]                    // remove elements
	if err != nil {
		// Exit the application if JSON can not be serialized
		panic("JSON could not be marshalled")
	}

	start := time.Now()

	for count := 0; count < 3; count++ {
		res, err := http.Post(h.TargetEndpoint, "application/json", bytes.NewBuffer(postBody))

		// If err is nil then request was successful
		if err == nil {
			log.Infof("Synced Batch of Size %d which returned status code %d and took %s", payloadsLength, res.StatusCode, time.Since(start))
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
