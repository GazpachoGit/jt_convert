package main

import (
	tcclient "jt_converter/internal/clients/tc_client"
	"jt_converter/internal/config"
	jt_list_getter "jt_converter/internal/http/handlers/jt_list_getter"
	loadfile "jt_converter/internal/http/handlers/load_file"
	ping "jt_converter/internal/http/handlers/ping"
	pmi_getter "jt_converter/internal/http/handlers/pmi_getter"
	pmi_list_getter "jt_converter/internal/http/handlers/pmi_list_getter"
	jtmng "jt_converter/internal/service/jt_manager"
	tc "jt_converter/internal/service/tc_service"
	xml "jt_converter/internal/service/xml_manager"
	"jt_converter/internal/storage/bbolt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("logger ready")

	cfg := config.MustLoad()

	storage := bbolt.New(cfg.JT.DBPath, log)
	defer storage.Close()

	xmlMngr := xml.NewXMLManager(log)

	tcClient := tcclient.NewTCClient(cfg.TC.TCURL, cfg.TC.User, cfg.TC.Password, log)
	tcService := tc.NewTCService(tcClient, log, cfg.JT.JtStoragePath)

	jt_manager := jtmng.New(
		cfg.JT.VisualizerPath,
		cfg.JT.JtStoragePath,
		cfg.JT.XmlStoragePath,
		storage,
		xmlMngr,
		tcService,
		log,
	)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/ping", ping.New(log))
	router.Route("/v1", func(r chi.Router) {
		r.Post("/jts/getPMIs", pmi_getter.New(log, jt_manager))
		r.Post("/jts/loadJT", loadfile.New(log, jt_manager))
		r.Get("/pmis", pmi_list_getter.New(log, jt_manager))
		r.Get("/jts", jt_list_getter.New(log, jt_manager))
	})

	log.Info("starting server", slog.String("address", cfg.HTTPSever.Address))
	srv := &http.Server{
		Addr:         cfg.HTTPSever.Address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
}
