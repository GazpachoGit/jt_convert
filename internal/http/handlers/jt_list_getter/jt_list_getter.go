package handlers

import (
	response "jt_converter/internal/http/model"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type JTGetter interface {
	GetJTList() ([]string, error)
}

func New(log *slog.Logger, getterJT JTGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.jt_list_getter"
		log = log.With(slog.String("op", op))

		list, err := getterJT.GetJTList()
		if err != nil {
			log.Error("failed to get JT list", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		render.JSON(w, r, map[string][]string{
			"jts": list,
		})
	}
}
