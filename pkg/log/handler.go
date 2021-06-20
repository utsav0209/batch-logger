package log

import (
	"batch-logger/pkg/utils"
	"log"
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
	payload []Log
}

func CreateHandler() Handler {
	return Handler{}
}

func (h Handler) HandleLogRequest(w http.ResponseWriter, r *http.Request) {
	// Fetch Payload from context and store in memory
	ctx := r.Context()
	payload := ctx.Value(utils.GetPayloadRequestId(ctx)).(Log)
	h.payload = append(h.payload, payload)

	// Write response
	w.WriteHeader(200)
	_, err := w.Write([]byte("Successfully Logged"))
	if err != nil {
		log.Printf("Response failed with error: %s", err)
		return
	}
}
