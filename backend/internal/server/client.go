package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"project/backend/core"
	"project/backend/internal/handler"
)

type Client struct {
	hub            *Hub
	conn           *websocket.Conn
	send           chan []byte
	logger         *log.Logger
	processManager *core.ProcessManager
	handlers       *handler.MessageHandler
}

func (c *Client) Send(message []byte) {
	c.send <- message
}

func (c *Client) Broadcast(message []byte) {
	c.hub.broadcast <- message
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
		action := msg["action"]

		messageHandler, ok := c.handlers.Get(action)
		if !ok {
			c.logger.Printf("Нет обработчика для: %s", action)
			continue
		}

		if err = messageHandler(msg, c, c.processManager); err != nil {
			c.logger.Printf("Ошибка в обработчике %s: %v", action, err)
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
