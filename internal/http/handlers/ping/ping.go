package handlers

import (
	"log/slog"
	"net/http"
)

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.ping.New"
		log := log.With(slog.String("op", op))
		log.Info("ping request received")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("pong"))
		if err != nil {
			log.Error("failed to write response", slog.String("err", err.Error()))
		}
	}
}
