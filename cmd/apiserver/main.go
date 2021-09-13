package main

import (
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"techtask/internal/app/apiserver"
	"time"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.json", "path to config file")
}

// @title TechTask App api
// @version 1.0
// @description API server for TechTask Application

// @host localhost:8080
// @BasePath /

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	config := apiserver.NewConfig()
	s := apiserver.New(config)
	go func() {
		if err := s.Run(); err != nil {
			logrus.Fatal(err)
		}
	}()
    // Graceful shutdown {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	logrus.Info("Graceful shutdown: Приложение заканчивает свое выполнение")

	if err := s.Shutdown(ctx); err != nil {
		logrus.Errorf("Graceful shutdown: Произошла ошибка при отключении сервера: %s", err.Error())
	}else{
		logrus.Info("Graceful shutdown: Приложение завершило выполнение")
	}
	// }
}
