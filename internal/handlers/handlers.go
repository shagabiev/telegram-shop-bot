package handlers

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// --- структура товара ---
type Product struct {
	Name        string
	Description string
	Price       float64
	PhotoURL    string
}

// --- встроенный каталог ---
var Catalog = []Product{
	{
		Name:        "GTMBAR Spark 8000 BLUEBERRY ICE",
		Description: "Голубика с холодком 2%",
		Price:       750,
		PhotoURL:    "https://nsk.ilfumoshop.ru/image/cache/import_files/ce/ce83f0a2-b04d-11ee-aee7-00155dcf0b04-600x600.jpeg",
	},
	{
		Name:        "GTMBAR Spark 8000 JUICY PEACH",
		Description: "Сочный персик 2%",
		Price:       750,
		PhotoURL:    "https://nsk.ilfumoshop.ru/image/cache/import_files/ce/ce83f0ac-b04d-11ee-aee7-00155dcf0b04-600x600.jpeg",
	},
	{
		Name:        "GTMBAR Spark 8000 FANTA ORANGE",
		Description: "Апельсиновая содовая 2%",
		Price:       750,
		PhotoURL:    "https://nsk.ilfumoshop.ru/image/cache/import_files/ce/ce83f0a8-b04d-11ee-aee7-00155dcf0b04-600x600.jpeg",
	},
	{
		Name:        "GTMBAR Spark 8000 BLUE RAZZ LEMONADE",
		Description: "Лимонад с голубой малиной 2%",
		Price:       750,
		PhotoURL:    "https://nsk.ilfumoshop.ru/image/cache/import_files/ce/ce83f0a0-b04d-11ee-aee7-00155dcf0b04-600x600.jpeg",
	},
	{
		Name:        "GTMBAR Spark 8000 CHERRY RASPBERRY",
		Description: "Вишня малина 2%",
		Price:       750,
		PhotoURL:    "https://nsk.ilfumoshop.ru/image/cache/import_files/ce/ce83f0a4-b04d-11ee-aee7-00155dcf0b04-600x600.jpeg",
	},
	{
		Name:        "GTMBAR Spark 8000 SOUR APPLE",
		Description: "Кислое яблоко 2%",
		Price:       750,
		PhotoURL:    "https://nsk.ilfumoshop.ru/image/cache/import_files/86/86e048ad-b04c-11ee-aee7-00155dcf0b04-600x600.jpeg",
	},
	{
		Name:        "GTMBAR Spark 8000 GREEN TEA",
		Description: "Зеленый чай 2%",
		Price:       750,
		PhotoURL:    "https://nsk.ilfumoshop.ru/image/cache/import_files/ce/ce83f0aa-b04d-11ee-aee7-00155dcf0b04-600x600.jpeg",
	},
	{
		Name:        "GTMBAR Spark 8000 MINT",
		Description: "Мята 2%",
		Price:       750,
		PhotoURL:    "https://nsk.ilfumoshop.ru/image/cache/import_files/ce/ce83f0b0-b04d-11ee-aee7-00155dcf0b04-600x600.jpeg",
	},
	{
		Name:        "GTMBAR Spark 8000 RASPBERRY WATERMELON",
		Description: "Малина арбуз 2%",
		Price:       750,
		PhotoURL:    "https://nsk.ilfumoshop.ru/image/cache/import_files/ce/ce83f0b0-b04d-11ee-aee7-00155dcf0b04-600x600.jpeg",
	},
}

// --- админский обработчик ---
type AdminHandler struct {
	Bot     *tgbotapi.BotAPI
	AdminID int64
}

func NewAdminHandler(bot *tgbotapi.BotAPI, adminID int64) *AdminHandler {
	return &AdminHandler{Bot: bot, AdminID: adminID}
}

// Добавление товара через бота
func (h *AdminHandler) AddProduct(text string, chatID int64) {
	parts := strings.Split(text, "|")
	if len(parts) != 4 {
		h.Bot.Send(tgbotapi.NewMessage(chatID, "❌ Формат: Название | Описание | Цена | ФотоURL"))
		return
	}

	price, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil {
		h.Bot.Send(tgbotapi.NewMessage(chatID, "❌ Неверная цена"))
		return
	}

	p := Product{
		Name:        strings.TrimSpace(parts[0]),
		Description: strings.TrimSpace(parts[1]),
		Price:       price,
		PhotoURL:    strings.TrimSpace(parts[3]),
	}

	Catalog = append(Catalog, p)
	h.Bot.Send(tgbotapi.NewMessage(chatID, "✔ Товар добавлен!"))
}

// Удаление товара по индексу
func (h *AdminHandler) DeleteProduct(indexStr string, chatID int64) {
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 || index >= len(Catalog) {
		h.Bot.Send(tgbotapi.NewMessage(chatID, "❌ Неверный индекс товара"))
		return
	}

	Catalog = append(Catalog[:index], Catalog[index+1:]...)
	h.Bot.Send(tgbotapi.NewMessage(chatID, "🗑 Товар удалён"))
}

// --- пользовательский обработчик ---
type UserHandler struct {
	Bot     *tgbotapi.BotAPI
	AdminID int64
}

func NewUserHandler(bot *tgbotapi.BotAPI, adminID int64) *UserHandler {
	return &UserHandler{Bot: bot, AdminID: adminID}
}

// Показ каталога
func (h *UserHandler) Catalog(chatID int64) {
	if len(Catalog) == 0 {
		h.Bot.Send(tgbotapi.NewMessage(chatID, "Каталог пуст!"))
		return
	}

	for idx, p := range Catalog {
		msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(p.PhotoURL))
		msg.Caption = fmt.Sprintf("%d. %s\n%s\nЦена: %.2f₽", idx, p.Name, p.Description, p.Price)

		// Добавляем кнопку "Купить"
		btn := tgbotapi.NewInlineKeyboardButtonData("Купить", strconv.Itoa(idx))
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btn))

		h.Bot.Send(msg)
	}
}

// Обработка покупки
func (h *UserHandler) HandleBuy(callback *tgbotapi.CallbackQuery) {
	idx, err := strconv.Atoi(callback.Data)
	if err != nil || idx < 0 || idx >= len(Catalog) {
		h.Bot.Send(tgbotapi.NewMessage(callback.From.ID, "❌ Ошибка покупки"))
		return
	}

	product := Catalog[idx]

	// уведомление администратору
	userLink := fmt.Sprintf("tg://user?id=%d", callback.From.ID)
	btn := tgbotapi.NewInlineKeyboardButtonURL("Написать пользователю", userLink)
	markup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btn))

	msg := tgbotapi.NewMessage(h.AdminID, fmt.Sprintf("Покупка!\nТовар: %s", product.Name))
	msg.ReplyMarkup = markup
	h.Bot.Send(msg)

	// подтверждение покупателю
	h.Bot.Send(tgbotapi.NewMessage(callback.From.ID, fmt.Sprintf("Вы выбрали: %s\nАдминистратор свяжется с вами.", product.Name)))
}
