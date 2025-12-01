package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/shagabiev/telegram-shop-bot/internal/bot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	token := os.Getenv("BOT_TOKEN")
	adminIDStr := os.Getenv("ADMIN_ID")
	adminID, _ := strconv.ParseInt(adminIDStr, 10, 64)

	b := bot.NewBot(token, adminID)
	log.Println("Bot is running...")
	b.Start()
}
