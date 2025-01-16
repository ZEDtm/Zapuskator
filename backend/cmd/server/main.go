package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"project/backend/config"
	"project/backend/core"
	"project/backend/server"
	"project/backend/utils"
)

var (
	version    string = "2.0.0"
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/server.toml", "path to config file")
}

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	path, err := os.Getwd()
	if err != nil {
		logger.Fatalf("Ошибка при получении текущей дериктории: %v", err)
	}
	logger.Printf("Текущая директория: %s", path)
	logger.Printf("Текущая версия: %s", version)

	newConfig := config.NewConfig(path)

	_, err = toml.DecodeFile(configPath, newConfig)
	if err != nil {
		logger.Fatalf("Ошибка конфигурации: %v", err)
	}

	processManager := core.NewProcessManager()

	var (
		upgrader = &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
	)

	sServer := server.NewServer(upgrader, processManager, logger)

	go sServer.Start(newConfig)

	done := make(chan struct{})

	utils.OpenBrowser(logger, newConfig)
	utils.StopCmd(logger, processManager, done)

	<-done
	logger.Printf("Завершено")
	os.Exit(0)
}
