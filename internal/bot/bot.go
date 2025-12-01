package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shagabiev/telegram-shop-bot/internal/handlers"
	"github.com/shagabiev/telegram-shop-bot/internal/keyboard"
)

type Bot struct {
	api   *tgbotapi.BotAPI
	user  *handlers.UserHandler
	admin *handlers.AdminHandler
}

func NewBot(token string, adminID int64) *Bot {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic("Ошибка создания бота: " + err.Error())
	}

	return &Bot{
		api:   api,
		user:  handlers.NewUserHandler(api, adminID),
		admin: handlers.NewAdminHandler(api, adminID),
	}
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			text := update.Message.Text

			// Админские команды
			if chatID == b.admin.AdminID {
				if strings.HasPrefix(text, "add ") {
					b.admin.AddProduct(text[4:], chatID)
					continue
				}
				if strings.HasPrefix(text, "del ") {
					b.admin.DeleteProduct(text[4:], chatID)
					continue
				}
			}

			// Пользовательские команды
			switch text {
			case "/start":
				msg := tgbotapi.NewMessage(chatID, "Добро пожаловать!")
				msg.ReplyMarkup = keyboard.MainMenu()
				b.api.Send(msg)

			case "📦 Каталог":
				b.user.Catalog(chatID)

			case "📖 Информация":
				contactInfo := "Розничная продажа в г.Казань (личная встреча) - 750₽\n" +
					"Оптовая продажа 450₽ (от 20 шт, личная встреча в г.Казань)"
				b.api.Send(tgbotapi.NewMessage(chatID, contactInfo))
			}
		}

		// Inline кнопки (покупка)
		if update.CallbackQuery != nil {
			b.user.HandleBuy(update.CallbackQuery)
			b.api.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, "Спасибо!"))
		}
	}
}
