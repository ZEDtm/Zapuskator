package server

import (
	"encoding/json"
	"project/backend/internal/handler"
	"strconv"
)

type MessageProcessInfo struct {
	Action           string `json:"action"`
	PID              int    `json:"pid"`
	Status           string `json:"status"`
	Edition          string `json:"edition"`
	Version          string `json:"version"`
	URLParam         string `json:"urlParam"`
	AdditionalFolder string `json:"additionalFolder"`
	Path             string `json:"path"`
}

type MessageServerInfo struct {
	Action  string `json:"action"`
	Edition string `json:"edition"`
	Version string `json:"version"`
	State   string `json:"state"`
}

func handleStart(msg map[string]string, client handler.Client, processManager handler.ProcessManager) error {
	path, edition, urlParam, version := msg["path"], msg["edition"], msg["urlParam"], msg["version"]

	pid, additionalFolder, err := processManager.StartApp(path, edition, urlParam)
	if err != nil {
		return err
	}

	message := &MessageProcessInfo{
		Action:           "start",
		PID:              pid,
		Status:           "running",
		Edition:          edition,
		Version:          version,
		URLParam:         urlParam,
		AdditionalFolder: additionalFolder,
		Path:             path,
	}

	byteMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	client.Broadcast(byteMessage)

	return nil
}

func handleStop(msg map[string]string, _ handler.Client, processManager handler.ProcessManager) error {
	pid, err := strconv.Atoi(msg["pid"])
	if err != nil {
		return err
	}

	err = processManager.StopApp(pid)
	if err != nil {
		return err
	}

	return nil
}

func handleGetServerInfo(msg map[string]string, client handler.Client, _ handler.ProcessManager) error {
	urlParam := msg["url"]

	edition, version, state, err := getServerInfo(urlParam)
	if err != nil {
		return err
	}

	message := &MessageServerInfo{
		Action:  "get_server_info",
		Edition: edition,
		Version: version,
		State:   state,
	}

	byteMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	client.Broadcast(byteMessage)

	return nil
}

func handleGetProcesses(_ map[string]string, client handler.Client, processManager handler.ProcessManager) error {
	processes := processManager.GetAllProcesses()
	for _, process := range processes {
		message := &MessageProcessInfo{
			Action:           "start",
			PID:              process.PID,
			Status:           process.Status,
			Edition:          process.Info.Edition,
			Version:          "unknown",
			URLParam:         process.Info.UrlParam,
			AdditionalFolder: process.Info.AdditionalFolder,
			Path:             "unknown",
		}

		byteMessage, err := json.Marshal(message)
		if err != nil {
			return err
		}

		client.Send(byteMessage)
	}

	return nil
}

func CreateHandlers() *handler.MessageHandler {
	handlers := handler.NewMessageHandlers()

	handlers.Handle("start", handleStart)
	handlers.Handle("stop", handleStop)
	handlers.Handle("get_server_info", handleGetServerInfo)
	handlers.Handle("get_processes", handleGetProcesses)

	return handlers
}
