package utils

import (
	"fmt"
	"log"
	"os/exec"
	"project/backend/config"
	"time"
)

func OpenBrowser(logger *log.Logger, config *config.Config) {
	if !config.OpenBrowser {
		return
	}

	go func() {
		time.Sleep(2 * time.Second)
		url := fmt.Sprintf("http://localhost%s/", config.Port)

		err := exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		if err != nil {
			logger.Printf("Ошибка при запуске браузера: %v", err)
		}

		logger.Printf("Запущен браузер по адресу: %s", url)
	}()
}
