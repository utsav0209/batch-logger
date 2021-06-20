package middleware

import (
	logHandler "batch-logger/pkg/log"
	"batch-logger/pkg/utils"
	"context"
	"encoding/json"
	"net/http"
)

func AddLogToCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var l logHandler.Log
		err := json.NewDecoder(r.Body).Decode(&l)

		if err != nil {
			http.Error(w, "JSON could not be deserialized", 400)
			return
		}

		ctx = context.WithValue(r.Context(), utils.GetPayloadRequestId(ctx), l)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
