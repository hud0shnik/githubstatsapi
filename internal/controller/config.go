package controller

import "time"

// Config - структура конфига сервера
type Config struct {
	ServerPort     string        `yaml:"server_port"`
	BasePath       string        `yaml:"base_path"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
}
