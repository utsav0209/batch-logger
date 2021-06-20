package log

import (
	"batch-logger/pkg/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Handler struct {
	BatchInterval  int    // Interval in seconds at which syncing should happen
	TargetEndpoint string // Target endpoint where JSON arrays should we eventually written
	payloads       Payloads
}

func CreateHandler(batchSize, batchInterval int, targetEndpoint string) Handler {
	return Handler{
		BatchInterval:  batchInterval,
		TargetEndpoint: targetEndpoint,
		payloads: Payloads{
			BatchSize: batchSize,
			logs:      make(chan Log),
		},
	}
}

func (h *Handler) SyncAtIntervals() {
	// Sync Payload Data in intervals to the provided endpoint
	for range time.Tick(time.Second * time.Duration(h.BatchInterval)) {
		if len(h.payloads.logs) > 0 {
			// Concurrently call to sync data to post endpoint
			go syncPayloadsToTargetEndpoint(nil, h.TargetEndpoint)
		}
	}
}

func (h *Handler) CollectLogs() {
	// Collect logs from Handler's payload channel and sync when data is full
	payloads := make([]Log, 0, h.payloads.BatchSize)
	for {
		payload := <-h.payloads.logs
		payloads = append(payloads, payload)

		if len(payloads) >= h.payloads.BatchSize {
			// Concurrently call to sync data to post endpoint
			go syncPayloadsToTargetEndpoint(payloads, h.TargetEndpoint)
			payloads = make([]Log, 0, h.payloads.BatchSize)
		}
	}
}

func (h *Handler) HandleLogRequest(w http.ResponseWriter, r *http.Request) {
	// Fetch Payload from context and store in memory
	ctx := r.Context()
	payload := ctx.Value(utils.GetPayloadRequestId(ctx)).(Log)
	h.payloads.Add(payload)

	// Write response
	w.WriteHeader(200)
	_, err := w.Write([]byte("Successfully Logged"))
	if err != nil {
		log.Errorf("Response failed with error: %s", err)
		return
	}
}
