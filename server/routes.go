package main

import (
	"database/sql"
	"log"
	"os"
	"ymb-cloz/internal/handler"
	"ymb-cloz/internal/service"
	"ymb-cloz/internal/store"

	"ymb-cloz/internal/bot"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func setupRoutes(r *gin.Engine, db *sql.DB) {
	// Initialize dependencies
	gameStore := store.NewGameStore(db)
	gameService := service.NewGameService(gameStore)
	gameHandler := handler.NewGameHandler(gameService)

	playerStore := store.NewPlayerStore(db)
	playerService := service.NewPlayerService(playerStore)
	playerHandler := handler.NewPlayerHandler(playerService)

	// Initialize Telegram bot
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if botToken == "" {
		log.Println("No Telegram bot token provided")
	}

	if botToken != "" {
		tgBot, err := tgbotapi.NewBotAPI(botToken)
		if err != nil {
			log.Printf("Error initializing Telegram bot: %v", err)
		} else {
			bot := bot.NewBot(tgBot, playerService)
			go bot.Start()
		}
	}

	api := r.Group("/api")
	{
		api.POST("/games", gameHandler.CreateGame)
		api.GET("/players", playerHandler.GetAllPlayers)
	}
}
