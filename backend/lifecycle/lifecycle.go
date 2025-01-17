package lifecycle

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ShutDowner interface {
	Shutdown(logger *log.Logger, ctx context.Context) error
}

func OnShutdown(logger *log.Logger, shutDowners ...ShutDowner) chan struct{} {
	done := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Printf("Сигнал: %v. Завершение работы...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		for _, shutDowner := range shutDowners {
			logger.Printf("Завершение работы компонента...")
			err := shutDowner.Shutdown(logger, ctx)
			if err != nil {
				logger.Printf("Ошибка при завершении работы компонента: %v", err)
			} else {
				logger.Printf("Компонент успешно завершён")
			}
		}

		close(done)
	}()

	return done
}
