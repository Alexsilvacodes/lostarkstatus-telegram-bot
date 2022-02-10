package main

import (
	"log"
	"os"
	"time"

	"github.com/Alexsilvacodes/LostArkStatus/lostarkstatus"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func sendMessage(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func buildMessage(update tgbotapi.Update, f func(m tgbotapi.MessageConfig)) {
	lostarkstatus.GetStatus(func(s []lostarkstatus.Server) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		statusStr := "📊 Servers Status\n\n" +
			"_Time: " + time.Now().Format("15:04:05") + " (GMT)_ \n\n"
		for _, server := range s {
			statusStr += "▪️ " + server.Name + ": " + server.Status + "\n"
		}
		msg.Text = statusStr
		msg.ParseMode = tgbotapi.ModeMarkdown

		f(msg)
	})
}

func main() {
	err := godotenv.Load(".env")
	token := os.Getenv("TELEGRAM_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
					ChatID:    update.Message.Chat.ID,
					MessageID: update.Message.MessageID,
				})

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.Text = "To retrieve Lost Ark server status execute /status"

				sendMessage(msg, bot)
			case "/status":
				bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
					ChatID:    update.Message.Chat.ID,
					MessageID: update.Message.MessageID,
				})

				buildMessage(update, func(msg tgbotapi.MessageConfig) {
					sendMessage(msg, bot)
				})
			}
		}
	}
}
