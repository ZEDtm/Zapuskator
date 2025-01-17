package server

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"project/backend/config"
	"project/backend/core"
	"project/backend/internal/handler"
	"strconv"
	"time"
)

type Server struct {
	hub            *Hub
	upgrader       *websocket.Upgrader
	processManager *core.ProcessManager
	logger         *log.Logger
	handlers       *handler.MessageHandler
}

func NewServer(logger *log.Logger, processManager *core.ProcessManager) *Server {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return &Server{
		hub:            NewHub(),
		upgrader:       upgrader,
		processManager: processManager,
		logger:         logger,
		handlers:       CreateHandlers(),
	}
}

func (s *Server) Start(config *config.Config) {
	mux := http.NewServeMux()

	static := http.FileServer(http.Dir(filepath.Join(config.Path, "static")))
	mux.Handle("/assets/", http.StripPrefix("/assets/", static))

	public := http.FileServer(http.Dir(filepath.Join(config.Path, "public")))
	mux.Handle("/", public)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.ServeWs(w, r)
	})

	go s.hub.Run()
	go s.StartProcessChecker()

	server := &http.Server{
		Addr:         config.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	s.logger.Printf("Сервер запущен на порту %s", config.Port)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Printf("Ошибка при запуске сервера: %v", err)
	}
}

func (s *Server) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Println(err)
		return
	}
	client := &Client{
		logger:         s.logger,
		hub:            s.hub,
		conn:           conn,
		send:           make(chan []byte, 256),
		processManager: s.processManager,
		handlers:       s.handlers,
	}
	s.hub.register <- client
	go client.writePump()
	go client.readPump()
}

func (s *Server) StartProcessChecker() {
	ticker := time.NewTicker(5 * time.Second) // Проверка каждые 5 секунд
	defer ticker.Stop()

	for range ticker.C {
		processes := s.processManager.GetAllProcesses()
		for _, process := range processes {
			if process.Status == "stopped" {
				message := []byte(`{"action": "stop", "pid": "` + strconv.Itoa(process.PID) + `", "status": "stopped"}`)
				s.hub.broadcast <- message
				if err := s.processManager.DeleteProcess(process.PID); err != nil {
					s.logger.Printf("Ошибка при попытке удалить процесс: %v", err)
				}
			}
		}
	}
}

func getServerInfo(urlParam string) (string, string, string, error) {
	endpoint := "/resto/get_server_info.jsp?encoding=UTF-8"
	fullURL := urlParam + endpoint

	resp, err := http.Get(fullURL)
	if err != nil {
		return "", "", "", fmt.Errorf("error fetching server info: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("error reading response: %v", err)
	}

	// Парсим XML
	type ServerInfo struct {
		Edition string `xml:"edition"`
		Version string `xml:"version"`
		State   string `xml:"serverState"`
	}

	var info ServerInfo
	err = xml.Unmarshal(body, &info)
	if err != nil {
		return "", "", "", fmt.Errorf("error parsing XML: %v", err)
	}

	return info.Edition, info.Version, info.State, nil
}
