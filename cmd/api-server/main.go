package main

import (
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ink-pwd/FileKeeper/internal/config"
	"github.com/ink-pwd/FileKeeper/internal/consts"
	"github.com/ink-pwd/FileKeeper/internal/handler"
	"github.com/ink-pwd/FileKeeper/internal/telegram"
	"github.com/ink-pwd/FileKeeper/logger"
)

func main() {
	var (
		mux    *http.ServeMux
		server *http.Server
		log    *logger.StdLogger
		tg     *tgbotapi.BotAPI
		bot    *telegram.TelegramStorage
		handl  *handler.FileHandler
		cfg    *config.Config
		err    error
	)

	log = logger.NewStdLogger()

	/*
		Получаем env конфигурацию
	*/
	cfg, err = config.Load()
	if err != nil {
		log.Fatal("load env: %s", err.Error())
		return
	}

	/*
		Подключаемся к телеграм боту и передаем в нашу структуру для более удобной работы
	*/
	tg, err = tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatal("bot connect: %s", err.Error())
		return
	}
	bot = telegram.NewTelegramStorage(cfg.ChatID, tg)
	log.Info("successful connection to telegram bot")

	/*
		Создаем файловый обработчик
	*/
	handl = handler.NewFileHandler(log, bot, cfg.MaxFileSize, cfg.MaxRamSize, cfg.Token)

	/*
		Настройка запросов и обработчиков сервера
	*/
	mux = http.NewServeMux()
	mux.HandleFunc(consts.UPLOAD, handl.UploadFile)
	mux.HandleFunc(consts.FILES, handl.GetFile)

	/*
		Создаем кастомный сервер ограничивая время записи и чтения.
		Значения лучше ограничить, что бы не создавались "бесконечные" соединения.
	*/
	server = &http.Server{
		Addr:         cfg.Host + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Error("server error: %s", err.Error())
		return
	}

	log.Info("the server is running on: %s%s", cfg.Host, cfg.Port)
}
