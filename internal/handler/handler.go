package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ink-pwd/FileKeeper/internal/telegram"
	"github.com/ink-pwd/FileKeeper/logger"
)

type FileHandler struct {
	log         *logger.StdLogger
	bot         *telegram.TelegramStorage
	maxFileSize int
	maxRamSize  int
	token       string
}

func NewFileHandler(log *logger.StdLogger, bot *telegram.TelegramStorage,
	maxFileSize, maxRamSize int, token string) *FileHandler {
	return &FileHandler{
		log:         log,
		bot:         bot,
		maxFileSize: maxFileSize,
		maxRamSize:  maxRamSize,
		token:       token,
	}
}

func (f *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		file      multipart.File
		header    *multipart.FileHeader
		messageID string
	)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	/*
		Ограничиваем максимальный размер файла.
		<<20 значит, что размер указан в магабайтайх.
	*/
	r.Body = http.MaxBytesReader(w, r.Body, int64(f.maxFileSize)<<20)
	/*
		Выделяем оперативную память для обработки запроса
	*/
	err = r.ParseMultipartForm(int64(f.maxRamSize) << 20)
	if err != nil {
		f.log.Error("parse multi form: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*
		Удаляем все временные файлы после завершения
	*/
	defer r.MultipartForm.RemoveAll()

	/*
		Получаем файл из формы.
		Прекращаем чтение по завершению
	*/
	file, header, err = r.FormFile("file")
	if err != nil {
		f.log.Error("get file err: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	/*
		Отправляем файл на сервер телеграм.
		Получаем ID сообщения.
	*/
	messageID, err = f.bot.SendFile(file, header.Filename)
	if err != nil {
		f.log.Error("send message telegram: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	f.log.Info("add file \"%s\" size: %d byte, ID: %s", header.Filename, header.Size, messageID)

	/*

	 */
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{"id": messageID})
}

func (f *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	var (
		fileID  string
		err     error
		file    tgbotapi.File
		fileURL string
		resp    *http.Response
	)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileID = r.PathValue("id")
	if fileID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	f.log.Info("get file: %s", fileID)
	/*
		Получаем файл через бота
	*/
	file, err = f.bot.GetFile(fileID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*
		Формируем ссылку для получения файла
	*/
	fileURL = fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", f.token, file.FilePath)
	/*
		Получаем файл
	*/
	resp, err = http.Get(fileURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		f.log.Error("get file: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	/*
		Отдаем клиенту
	*/
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		f.log.Error("sending file: %s", err.Error())
		return
	}
}
