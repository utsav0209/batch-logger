package log

import "time"

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

type Payloads struct {
	batchSize int         // Maximum batch size
	logs      chan Log    // Channel to stream logs
	forceSync chan string // Channel to sync logs on intervals
}

func (p Payloads) Add(log Log) {
	p.logs <- log
}

func (p Payloads) sync() {
	p.forceSync <- "sync"
}
