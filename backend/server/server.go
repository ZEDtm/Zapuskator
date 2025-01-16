package server

import (
	"encoding/xml"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"project/backend/config"
	"project/backend/core"
	"strconv"
	"time"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Server struct {
	hub            *Hub
	upgrader       *websocket.Upgrader
	processManager *core.ProcessManager
	logger         *log.Logger
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func NewServer(upgrader *websocket.Upgrader, processManager *core.ProcessManager, logger *log.Logger) *Server {
	s := &Server{
		hub:            NewHub(),
		upgrader:       upgrader,
		processManager: processManager,
		logger:         logger,
	}
	go s.hub.Run()
	go s.StartProcessChecker()
	return s
}

func (s *Server) Start(config *config.Config) {
	static := http.FileServer(http.Dir(filepath.Join(config.Path, "static")))
	idx := http.FileServer(http.Dir(config.Path))

	http.Handle("/assets/", http.StripPrefix("/assets/", static))
	http.Handle("/", idx)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.ServeWs(w, r)
	})

	s.logger.Printf("Сервер запущен на порту %s", config.Port)
	s.logger.Fatal(http.ListenAndServe(config.Port, nil))
}

func (s *Server) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Println(err)
		return
	}
	client := &Client{hub: s.hub, conn: conn, send: make(chan []byte, 256), logger: s.logger, processManager: s.processManager}
	client.hub.register <- client
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
