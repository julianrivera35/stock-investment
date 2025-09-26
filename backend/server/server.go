package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	s := &Server{
		router: mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	//Health check
	s.router.HandleFunc("/health", s.healthCheck).Methods("GET")

	//API Routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	//Companies
	api.HandleFunc("/companies", s.getCompanies).Methods("GET")
	api.HandleFunc("/companies/{ticker}", s.getCompanyByTicker).Methods("GET")

	//Brokerages
	api.HandleFunc("/brokerages", s.getBrokerages).Methods("GET")

	//Recommendations
	api.HandleFunc("/recommendations", s.getRecommendations).Methods("GET")
	api.HandleFunc("/recommendations/company/{ticker}", s.getRecommendationsByTicker).Methods("GET")
	api.HandleFunc("/recommendations/brokerage/{brokerage}", s.getRecommendationsByBrokerage).Methods("GET")

	// CORS Middleware
	s.router.Use(corsMiddleware)
}

// CORS Middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) Start(port string) {
	fmt.Printf("🚀 Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, s.router))
}

// Health check handler
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "ok",
		"message":   "Server is running",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "stock-investment-backend",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
