package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"project/backend/config"
	"project/backend/core"
	"project/backend/internal/lifecycle"
	"project/backend/internal/server"
	"project/backend/internal/utils"
)

var (
	version    = "2.0.0"
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

	sServer := server.NewServer(logger, processManager)

	go sServer.Start(newConfig)

	utils.OpenBrowser(logger, newConfig)

	done := lifecycle.OnShutdown(logger, processManager)
	<-done
	logger.Printf("Завершено")
	os.Exit(0)
}
