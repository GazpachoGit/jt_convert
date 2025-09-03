package handlers

import (
	response "jt_converter/internal/http/model"
	model "jt_converter/internal/storage/model/pmis"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type Request struct {
	JtFileName string `json:"jt_file_name"`
}

type PMIGetter interface {
	GetPMIs(jtFileName string) (*model.Model, error)
}

func New(log *slog.Logger, getterPMI PMIGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pmi_getter"
		log = log.With(slog.String("op", op))

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("can't decode request body", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("can't decode request body"))
			return
		}

		if req.JtFileName == "" {
			log.Error("empty JT file name")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty JT file name"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		model, err := getterPMI.GetPMIs(req.JtFileName)
		if err != nil {
			log.Error("error during GetPMIs", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		render.JSON(w, r, model)
	}
}
