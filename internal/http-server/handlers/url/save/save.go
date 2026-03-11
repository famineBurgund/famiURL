package save

import (
	"net/http"

	response "github.com/famineBurgund/famiURL/internal/lib/api/response"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL   string `json:"url" validate"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("fail decode request body", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))
	}
}
