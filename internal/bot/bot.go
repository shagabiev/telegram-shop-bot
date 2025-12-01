package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shagabiev/telegram-shop-bot/internal/handlers"
	"github.com/shagabiev/telegram-shop-bot/internal/keyboard"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	handler *handlers.Handler
}

func NewBot(token string, adminID int64) *Bot {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic("Ошибка создания бота: " + err.Error())
	}
	handler := handlers.NewHandler(api, adminID)
	return &Bot{api: api, handler: handler}
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // обычные сообщения
			chatID := update.Message.Chat.ID
			text := update.Message.Text

			if chatID == b.handler.AdminID {
				if strings.HasPrefix(text, "add ") {
					b.handler.AddProduct(text[4:], chatID)
					continue
				}
				if strings.HasPrefix(text, "del ") {
					b.handler.DeleteProduct(text[4:], chatID)
					continue
				}
			}

			switch text {
			case "/start":
				msg := tgbotapi.NewMessage(chatID, "Добро пожаловать!")
				msg.ReplyMarkup = keyboard.MainMenu()
				b.api.Send(msg)

			case "📦 Каталог":
				b.handler.Catalog(chatID)
			}
		}

		// Обработка Inline кнопок
		if update.CallbackQuery != nil {
			b.handler.HandleBuy(update.CallbackQuery)
		}
	}
}
