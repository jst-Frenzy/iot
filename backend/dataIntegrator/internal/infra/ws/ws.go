package ws

import (
	"context"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"sync"
)

type DataForSend struct {
	Temperature int32
	Wet         int32
}

type Server struct {
	addr     string
	upgrader websocket.Upgrader

	clients map[*websocket.Conn]struct{}
	mu      sync.Mutex

	logger *slog.Logger
}

func New(addr string) *Server {
	return &Server{
		addr: addr,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients: make(map[*websocket.Conn]struct{}),
		logger:  slog.With("service", "websocket"),
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", s.handleWS)

	server := &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()

		s.logger.Info("shutting down websocket server")

		if err := server.Shutdown(context.Background()); err != nil {
			s.logger.Error("shutdown error", "error", err)
		}

		s.closeAllClients()
	}()

	s.logger.Info("websocket server started", "address", s.addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) handleWS(
	w http.ResponseWriter,
	r *http.Request,
) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("upgrade error", "error", err)
		return
	}

	s.mu.Lock()
	s.clients[conn] = struct{}{}
	s.mu.Unlock()

	s.logger.Info("client connected")

	go s.listenClient(conn)
}

func (s *Server) listenClient(conn *websocket.Conn) {
	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()

		_ = conn.Close()

		s.logger.Info("client disconnected")
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func (s *Server) Send(v DataForSend) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.clients {
		if err := conn.WriteJSON(v); err != nil {
			s.logger.Error("broadcast error", "error", err)

			_ = conn.Close()
			delete(s.clients, conn)
		}
	}
}

func (s *Server) closeAllClients() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.clients {
		_ = conn.Close()
		delete(s.clients, conn)
	}
}
