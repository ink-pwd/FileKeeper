package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ink-pwd/FileKeeper/internal/consts"
	"github.com/ink-pwd/FileKeeper/internal/handler"
	"github.com/ink-pwd/FileKeeper/internal/telegram"
	"github.com/ink-pwd/FileKeeper/logger"
	"github.com/joho/godotenv"
)

func main() {
	var (
		mux          *http.ServeMux
		server       *http.Server
		log          *logger.StdLogger
		tg           *tgbotapi.BotAPI
		bot          *telegram.TelegramStorage
		handl        *handler.FileHandler
		err          error
		readEnv      string
		host         string
		port         string
		token        string
		chatID       int
		maxFileSize  int
		maxRamSize   int
		readTimeOut  int
		writeTimeOut int
	)

	/*
		Получаем env конфигурацию
	*/
	err = godotenv.Load()
	if err != nil {
		log.Fatal("load env : %s", err.Error())
		return
	}

	readEnv = os.Getenv("READTIMEOUT")
	readTimeOut, err = strconv.Atoi(readEnv)
	if err != nil {
		log.Fatal("READTIMEOUT was entered incorrectly: %s", err.Error())
		return
	}

	readEnv = os.Getenv("WRITETIMEOUT")
	writeTimeOut, err = strconv.Atoi(readEnv)
	if err != nil {
		log.Fatal("WRITETIMEOUT was entered incorrectly: %s", err.Error())
		return
	}

	readEnv = os.Getenv("MAXFILESIZEMEGABYTE")
	maxFileSize, err = strconv.Atoi(readEnv)
	if err != nil {
		log.Fatal("max file size was entered incorrectly: %s", err.Error())
		return
	}

	readEnv = os.Getenv("MAXRAMSIZEMEGABYTE")
	maxRamSize, err = strconv.Atoi(readEnv)
	if err != nil {
		log.Fatal("max ram size was entered incorrectly: %s", err.Error())
		return
	}

	readEnv = os.Getenv("CHATID")
	chatID, err = strconv.Atoi(readEnv)
	if err != nil {
		log.Fatal("chat ID was entered incorrectly: %s", err.Error())
		return
	}

	host = os.Getenv("HOST")
	port = os.Getenv("PORT")
	token = os.Getenv("TOKEN")

	log = logger.NewStdLogger()

	/*
		Подключаемся к телеграм боту и передаем в нашу структуру для более удобной работы
	*/
	tg, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("bot connect: %s", err.Error())
		return
	}
	bot = telegram.NewTelegramStorage(chatID, tg)
	log.Info("successful connection to telegram bot")

	/*
		Создаем файловый обработчик
	*/
	handl = handler.NewFileHandler(log, bot, maxFileSize, maxRamSize, token)

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
		Addr:         host + port,
		Handler:      mux,
		ReadTimeout:  time.Duration(readTimeOut) * time.Second,
		WriteTimeout: time.Duration(writeTimeOut) * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Error("server error: %s", err.Error())
		return
	}

	log.Info("the server is running on: %s%s", host, port)
}
