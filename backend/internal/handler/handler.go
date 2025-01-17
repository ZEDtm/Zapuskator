package handler

import "project/backend/core"

type Client interface {
	Send(message []byte)
	Broadcast(message []byte)
}

type ProcessManager interface {
	StartApp(path, edition, urlParam string) (int, string, error)
	StopApp(pid int) error
	GetAllProcesses() []*core.Process
}

type Handler func(message map[string]string, client Client, manager ProcessManager) error
type MessageHandler struct {
	handlers map[string]Handler
}

func NewMessageHandlers() *MessageHandler {
	return &MessageHandler{
		handlers: make(map[string]Handler),
	}
}

func (mh *MessageHandler) Handle(action string, handler Handler) {
	mh.handlers[action] = handler
}

func (mh *MessageHandler) Get(action string) (Handler, bool) {
	handler, ok := mh.handlers[action]
	return handler, ok
}
