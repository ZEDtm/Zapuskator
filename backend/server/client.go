package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"project/backend/core"
	"strconv"
)

type Client struct {
	hub            *Hub
	conn           *websocket.Conn
	send           chan []byte
	logger         *log.Logger
	processManager *core.ProcessManager
}

func (c *Client) readPump() {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Printf("Паника в readPump: %v", r)
		}
		c.hub.unregister <- c
		err := c.conn.Close()
		if err != nil {
			c.logger.Printf("Ошибка при закрытии соединения с клиентом: %v", err)
		}
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Printf("WebSocket ошибка при чтении сообщения: %v", err)
			}
			break
		}

		var msg map[string]string
		if err = json.Unmarshal(message, &msg); err != nil {
			c.logger.Printf("Ошибка JSON дешифрования: %v", err)
			continue
		}

		switch msg["action"] {
		case "start": // Необходимо передать path edition urlParam
			path, edition, urlParam, version := msg["path"], msg["edition"], msg["urlParam"], msg["version"]
			pid, additionalFolder, err := c.processManager.StartApp(path, edition, urlParam)

			if err != nil {
				c.logger.Printf("Ошибка при запуске приложения: %v", err)
				continue
			}
			c.hub.broadcast <- []byte(`{"action": "start", "pid": "` + strconv.Itoa(pid) + `", "edition": "` + edition + `",  "urlParam": "` + urlParam + `", "additionalFolder": "` + additionalFolder + `", "version": "` + version + `", "status": "running"}`)

		case "stop":
			pid, err := strconv.Atoi(msg["pid"])
			if err != nil {
				c.logger.Printf("Неверный PID: %v", err)
				continue
			}
			err = c.processManager.StopApp(pid)
			if err != nil {
				c.logger.Printf("Ошибка при остановке приложения: %v", err)
				continue
			}
		case "get_server_info":
			urlParam := msg["url"]
			edition, version, state, err := getServerInfo(urlParam)
			if err != nil {
				c.logger.Printf("Ошибка при получении информации о сервере: %v", err)
				continue
			}
			c.hub.broadcast <- []byte(`{"action": "server_info", "edition": "` + edition + `", "version": "` + version + `", "state": "` + state + `"}`)
		case "get_processes":
			processes := c.processManager.GetAllProcesses()
			for _, process := range processes {
				c.send <- []byte(`{"action": "start", "pid": "` + strconv.Itoa(process.PID) + `", "edition": "` + process.Info.Edition + `",  "urlParam": "` + process.Info.UrlParam + `", "additionalFolder": "` + process.Info.AdditionalFolder + `", "version": "unknown", "status": "` + process.Status + `"}`)
			}

		}
	}
}

func (c *Client) writePump() {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Printf("Паника в writePump: %v", r)
		}
		if c.conn != nil {
			err := c.conn.Close()
			if err != nil {
				c.logger.Printf("Соединение было неожиданно закрыто на стороне клиента: %v", err)
			}
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.logger.Println("Канал отправки закрыт")
				return
			}

			if c.conn != nil {
				err := c.conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					c.logger.Printf("Ошибка WebSocket при записи сообщения: %v", err)
					return
				}
			}
		}
	}
}
