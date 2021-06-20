package log

import (
	"batch-logger/pkg/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Handler struct {
	batchInterval  int    // Interval in seconds at which syncing should happen
	targetEndpoint string // Target endpoint where JSON arrays should we eventually written
	payloads       Payloads
}

func CreateHandler(batchSize, batchInterval int, targetEndpoint string) Handler {
	return Handler{
		batchInterval:  batchInterval,
		targetEndpoint: targetEndpoint,
		payloads: Payloads{
			batchSize: batchSize,
			logs:      make(chan Log),
			forceSync: make(chan string),
		},
	}
}

func (h *Handler) SyncAtIntervals() {
	// Sync Payload Data in intervals to the provided endpoint
	for range time.Tick(time.Second * time.Duration(h.batchInterval)) {
		h.payloads.sync()
	}
}

func (h *Handler) CollectLogs() {
	// Collect logs from Handler's payload channel and forceSync channel to sync data concurrently
	payloads := make([]Log, 0, h.payloads.batchSize)

	for {
		select {
		case p := <-h.payloads.logs: // Collect logs from payload channel
			payloads = append(payloads, p)
			if len(payloads) >= h.payloads.batchSize {
				// Concurrently call to sync data to post endpoint
				go syncPayloadsToTargetEndpoint(payloads, h.targetEndpoint)
				payloads = make([]Log, 0, h.payloads.batchSize)
			}
		case <-h.payloads.forceSync: // Forcefully sync on intervals
			if len(payloads) > 0 {
				// Concurrently call to sync data to post endpoint
				go syncPayloadsToTargetEndpoint(payloads, h.targetEndpoint)
				payloads = make([]Log, 0, h.payloads.batchSize)
			}
		}
	}
}

func (h *Handler) HandleLogRequest(w http.ResponseWriter, r *http.Request) {
	// Fetch Payload from context and stream it
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
