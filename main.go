package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/neetusuthar90/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		fmt.Println("Debugging enabled")
	} else {
		fmt.Println("Debugging disabled")
	}

	router := chi.NewRouter()

	// Mount the file server for the /app/* path
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	router.Mount("/app", apiCfg.middlewareMetricsInc(middleware.Logger(fsHandler)))

	// API router endpoints
	apiRouter := chi.NewRouter()
	router.Mount("/api", apiCfg.middlewareMetricsInc(middleware.Logger(apiRouter)))
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Get("/chirps", apiCfg.handlerChirpsRetrieve)
	apiRouter.Get("/chirps/{id}", apiCfg.handlerChirpyIDRetrieve)
	apiRouter.Post("/users", apiCfg.handlerUsersCreate)

	// Admin endpoints
	router.Get("/metrics", apiCfg.handlerMetrics)

	corsMux := middlewareCors(router)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
