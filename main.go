package main

import (
	"encoding/json"
	"log"
	"maps"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/olegrom32/warehouse/internal/packager"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{}))

	// SIZES="[250, 500, 1000, 2000, 5000]"
	boxSizesEnv := os.Getenv("SIZES")
	var sizes []int

	if err := json.Unmarshal([]byte(boxSizesEnv), &sizes); err != nil {
		// Default value
		sizes = []int{250, 500, 1000, 2000, 5000}
	}

	log.Printf("Starting app with box sizes: %v", sizes)

	packagerSvc := packager.NewService(sizes)

	r.Post("/order", func(w http.ResponseWriter, r *http.Request) {
		// Always respond with json
		w.Header().Set("content-type", "application/json")

		// req is the request body
		var req struct {
			Items int `json:"items"`
		}

		// Unmarshal the request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		// Call packager service
		res := packagerSvc.Package(req.Items)

		log.Printf("[POST /order] Packaging the total of %d items: %v", req.Items, res)

		// response is the response payload
		response := struct {
			Result map[int]int `json:"result"`
		}{
			Result: map[int]int{},
		}

		maps.Copy(response.Result, res)

		// Marshal the response
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		// Send back the response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(responseBytes)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
