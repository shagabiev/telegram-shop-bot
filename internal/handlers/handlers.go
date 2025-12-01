package handlers

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shagabiev/telegram-shop-bot/internal/models"
)

type Handler struct {
	Bot      *tgbotapi.BotAPI
	AdminID  int64
	Products []models.Product
}

func NewHandler(bot *tgbotapi.BotAPI, adminID int64) *Handler {
	return &Handler{Bot: bot, AdminID: adminID, Products: []models.Product{}}
}

// Добавление товара админом
func (h *Handler) AddProduct(text string, chatID int64) {
	parts := strings.Split(text, "|")
	if len(parts) != 4 {
		h.Bot.Send(tgbotapi.NewMessage(chatID, "❌ Формат: Название | Описание | Цена | ФотоURL"))
		return
	}

	price, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil {
		h.Bot.Send(tgbotapi.NewMessage(chatID, "❌ Некорректная цена"))
		return
	}

	p := models.Product{
		Name:        strings.TrimSpace(parts[0]),
		Description: strings.TrimSpace(parts[1]),
		Price:       price,
		PhotoURL:    strings.TrimSpace(parts[3]),
	}

	h.Products = append(h.Products, p)
	h.Bot.Send(tgbotapi.NewMessage(chatID, "✔ Товар добавлен!"))
}

// Удаление товара по индексу
func (h *Handler) DeleteProduct(indexStr string, chatID int64) {
	idx, err := strconv.Atoi(indexStr)
	if err != nil || idx < 1 || idx > len(h.Products) {
		h.Bot.Send(tgbotapi.NewMessage(chatID, "❌ Некорректный номер товара"))
		return
	}

	h.Products = append(h.Products[:idx-1], h.Products[idx:]...)
	h.Bot.Send(tgbotapi.NewMessage(chatID, "🗑 Товар удалён"))
}

// Показ каталога с кнопкой "Купить"
func (h *Handler) Catalog(chatID int64) {
	if len(h.Products) == 0 {
		h.Bot.Send(tgbotapi.NewMessage(chatID, "Каталог пуст"))
		return
	}

	for i, p := range h.Products {
		text := fmt.Sprintf("%d. %s\n%s\nЦена: %.2f₽", i+1, p.Name, p.Description, p.Price)

		msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(p.PhotoURL))
		msg.Caption = text

		// Inline-кнопка "Купить"
		buyBtn := tgbotapi.NewInlineKeyboardButtonData("Купить", fmt.Sprintf("buy_%d", i))
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buyBtn),
		)

		h.Bot.Send(msg)
	}
}

// Обработка нажатия "Купить"
func (h *Handler) HandleBuy(callback *tgbotapi.CallbackQuery) {
	data := callback.Data
	if !strings.HasPrefix(data, "buy_") {
		return
	}

	indexStr := strings.TrimPrefix(data, "buy_")
	idx, err := strconv.Atoi(indexStr)
	if err != nil || idx < 0 || idx >= len(h.Products) {
		h.Bot.Send(tgbotapi.NewMessage(callback.Message.Chat.ID, "❌ Ошибка покупки"))
		return
	}

	product := h.Products[idx]

	// Ответ пользователю
	h.Bot.Send(tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("Вы выбрали товар: %s. Администратор свяжется с вами.", product.Name)))

	// Уведомление администратору с кнопкой "Написать пользователю"
	userLink := fmt.Sprintf("tg://user?id=%d", callback.From.ID)
	btn := tgbotapi.NewInlineKeyboardButtonURL("Написать пользователю", userLink)
	markup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btn))

	msg := tgbotapi.NewMessage(h.AdminID, fmt.Sprintf("Покупка!\nТовар: %s", product.Name))
	msg.ReplyMarkup = markup // <- присваиваем напрямую
	h.Bot.Send(msg)

}
