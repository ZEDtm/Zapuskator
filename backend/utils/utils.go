package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"project/backend/config"
	"project/backend/core"
	"syscall"
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

func StopCmd(logger *log.Logger, pm *core.ProcessManager, done chan struct{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Printf("Сигнал: %v. Завершение работы...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		processes := pm.GetAllProcesses()
		logger.Printf("Найдено процессов для остановки: %d", len(processes))

		for _, process := range processes {
			logger.Printf("Остановка процесса %d...", process.PID)
			err := pm.StopApp(process.PID)
			if err != nil {
				logger.Printf("Ошибка при завершении процесса %d: %v", process.PID, err)
			} else {
				logger.Printf("Процесс %d успешно остановлен", process.PID)
			}

			select {
			case <-ctx.Done():
				logger.Printf("Таймаут при остановке процессов")
				break
			default:
				continue
			}
		}

		close(done)
	}()
}
