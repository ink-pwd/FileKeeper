package config

import (
	"os"
	"strconv"
)

type Config struct {
	Host         string
	Port         string
	Token        string
	ChatID       int
	ReadTimeout  int
	WriteTimeout int
	MaxFileSize  int
	MaxRamSize   int
}

func Load() (*Config, error) {
	var cfg Config
	var err error

	cfg.Host = os.Getenv("HOST")
	cfg.Port = os.Getenv("PORT")
	cfg.Token = os.Getenv("TOKEN")

	cfg.ChatID, err = strconv.Atoi(os.Getenv("CHATID"))
	if err != nil {
		return nil, err
	}

	cfg.ReadTimeout, err = strconv.Atoi(os.Getenv("READTIMEOUT"))
	if err != nil {
		return nil, err
	}

	cfg.WriteTimeout, err = strconv.Atoi(os.Getenv("WRITETIMEOUT"))
	if err != nil {
		return nil, err
	}

	cfg.MaxFileSize, err = strconv.Atoi(os.Getenv("MAXFILESIZEMEGABYTE"))
	if err != nil {
		return nil, err
	}

	cfg.MaxRamSize, err = strconv.Atoi(os.Getenv("MAXRAMSIZEMEGABYTE"))
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
