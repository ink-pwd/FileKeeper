package telegram

import (
	"mime/multipart"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramStorage struct {
	chat int
	bot  *tgbotapi.BotAPI
}

func NewTelegramStorage(chatID int, bot *tgbotapi.BotAPI) *TelegramStorage {
	return &TelegramStorage{
		chat: chatID,
		bot:  bot,
	}
}

func (t *TelegramStorage) SendFile(file multipart.File, name string) (string, error) {
	var (
		tgFile tgbotapi.FileReader
		msg    tgbotapi.DocumentConfig
		result tgbotapi.Message
		err    error
	)
	/*
		Создаем файл телеграмма из полученного файла
	*/
	tgFile = tgbotapi.FileReader{
		Name:   name,
		Reader: file,
	}
	/*
		Формируем сообщение
	*/
	msg = tgbotapi.NewDocument(int64(t.chat), tgFile)
	/*
		Отправляем сообщение
	*/
	result, err = t.bot.Send(msg)
	if err != nil {
		return "", err
	}
	return result.Document.FileID, err
}
func (t *TelegramStorage) GetFile(id string) (tgbotapi.File, error) {
	var (
		fileConfig tgbotapi.FileConfig
	)
	fileConfig = tgbotapi.FileConfig{FileID: id}
	/*
		Возвращаем данные о файле для дальнейшего получения его
	*/
	return t.bot.GetFile(fileConfig)
}
