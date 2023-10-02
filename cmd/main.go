package main

import (
	"os"
	"time"

	"github.com/hud0shnik/githubstatsapi/controllers"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// config - структура конфигов
type config struct {
	Server controllers.Config `yaml:"server"`
}

// configure получает конфиги из файла
func configure(fileName string) (*config, error) {

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var config config

	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {

	// Настройка логгера
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
		PrettyPrint:     true,
	})

	// Получение конфигов
	config, err := configure("config.yaml")
	if err != nil {
		logrus.WithError(err).Fatal("can't read config")
	}

	// Вывод времени начала работы
	logrus.Printf("API Starts at %s port", config.Server.ServerPort)

	// Создание сервера
	s := controllers.NewServer(&config.Server)

	// Запуск API
	logrus.Fatal(s.Server.ListenAndServe())

}
