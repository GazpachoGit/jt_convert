package loadfile

import (
	response "jt_converter/internal/http/model"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type LoadFileRequest struct {
	Uid      string `json:"uid"`
	TypeName string `json:"type"`
	Name     string `json:"name"`
}

type Service interface {
	LoadFile(uid, typeName, name string) error
}

func New(log *slog.Logger, tc_service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.getdataset.New"
		log := log.With(slog.String("op", op))

		var req LoadFileRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		if req.Uid == "" || req.TypeName == "" || req.Name == "" {
			log.Error("invalid body", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid body"))
			return
		}
		log.Info("received request to load file from TC", slog.String("uid", req.Uid), slog.String("type", req.TypeName))

		err = tc_service.LoadFile(req.Uid, req.TypeName, req.Name)
		if err != nil {
			log.Error("failed to load file from TC", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		log.Info("successfully loaded file from TC")
		render.JSON(w, r, response.Response{Status: "OK"})
	}
}
