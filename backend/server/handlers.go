package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-investment-backend/service"

	"github.com/gorilla/mux"
)

// Response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func (s *Server) getCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := service.GetAllCompanies()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponse(w, companies, nil)
}

func (s *Server) getCompanyByTicker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticker := vars["ticker"]

	company, err := service.GetCompanyByTicker(ticker)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	sendSuccessResponse(w, company, nil)
}

func (s *Server) getBrokerages(w http.ResponseWriter, r *http.Request) {
	brokerages, err := service.GetAllBrokerages()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponse(w, brokerages, nil)
}

func (s *Server) getRecommendations(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	ticker := r.URL.Query().Get("ticker")
	brokerageID := r.URL.Query().Get("borkerage_id")

	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	recommendations, total, err := service.GetRecommendations(limit, offset, ticker, brokerageID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	meta := &Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	sendSuccessResponse(w, recommendations, meta)
}

func (s *Server) getRecommendationsByTicker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticker := vars["ticker"]

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	recommendations, total, err := service.GetRecommendations(limit, offset, ticker, "")
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	meta := &Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	sendSuccessResponse(w, recommendations, meta)
}

func (s *Server) getRecommendationsByBrokerage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	brokerageID := vars["id"]

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	recommendations, total, err := service.GetRecommendations(limit, offset, "", brokerageID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	meta := &Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	sendSuccessResponse(w, recommendations, meta)
}

func sendSuccessResponse(w http.ResponseWriter, data interface{}, meta *Meta) {
	w.Header().Set("Content-Type", "application/json")
	response := APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := APIResponse{
		Success: false,
		Error:   message,
	}
	json.NewEncoder(w).Encode(response)
}
