package server

import (
	"context"
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

type Server struct {
	httpServer *http.Server
	metrics    *Metrics
}

type Metrics struct {
	BotRequestsCount int64
	BotIsConnected   bool
	AMLRequestsCount int64
	AMLIsConnected   bool
}

type StatusResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Bot       struct {
		RequestsCount int64 `json:"requests_count"`
		IsConnected   bool  `json:"is_connected"`
	} `json:"bot"`
	AML struct {
		RequestsCount int64 `json:"requests_count"`
		IsConnected   bool  `json:"is_connected"`
	} `json:"aml"`
}

func New(addr string) *Server {
	mux := http.NewServeMux()
	srv := &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		metrics: &Metrics{},
	}

	mux.HandleFunc("/status", srv.handleStatus)

	return srv
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) IncrementBotRequests() {
	atomic.AddInt64(&s.metrics.BotRequestsCount, 1)
}

func (s *Server) SetBotConnected(connected bool) {
	s.metrics.BotIsConnected = connected
}

func (s *Server) IncrementAMLRequests() {
	atomic.AddInt64(&s.metrics.AMLRequestsCount, 1)
}

func (s *Server) SetAMLConnected(connected bool) {
	s.metrics.AMLIsConnected = connected
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := StatusResponse{
		Status:    "ok",
		Timestamp: time.Now(),
	}
	response.Bot.RequestsCount = atomic.LoadInt64(&s.metrics.BotRequestsCount)
	response.Bot.IsConnected = s.metrics.BotIsConnected
	response.AML.RequestsCount = atomic.LoadInt64(&s.metrics.AMLRequestsCount)
	response.AML.IsConnected = s.metrics.AMLIsConnected

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
