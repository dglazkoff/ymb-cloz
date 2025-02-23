package bot

import (
	"fmt"
	"log"
	"strings"

	"ymb-cloz/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot           *tgbotapi.BotAPI
	playerService *service.PlayerService
}

func NewBot(bot *tgbotapi.BotAPI, playerService *service.PlayerService) *Bot {
	return &Bot{
		bot:           bot,
		playerService: playerService,
	}
}

func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}

func (b *Bot) handleHelp(c *tgbotapi.Update) error {
	helpText := `🎮 *YMB Cloz Bot* 🎮

Available commands:
/help \- Show this help message
/top\_winrate \- Show players sorted by win rate
/top\_games \- Show players sorted by games played
/top\_captains \- Show top captains by win rate
/top\_role \<role\> \- Show top players by role \(carry/mid/offlane/pos4/pos5\)

Example:
/top\_role carry \- Show top carry players`

	return b.sendMessage(c.Message.Chat.ID, helpText)
}

func (b *Bot) handleTopWinRate(c *tgbotapi.Update) error {
	stats, err := b.playerService.GetTopByWinRate()
	if err != nil {
		log.Printf("Error getting top win rates: %v", err)
		return b.sendMessage(c.Message.Chat.ID, "Error fetching statistics")
	}

	if len(stats) == 0 {
		return b.sendMessage(c.Message.Chat.ID, "No statistics available")
	}

	response := "*Top players by win rate:*\n\n"
	for i, stat := range stats {
		response += fmt.Sprintf("%d\\. *%s* \\- %s\n",
			i+1,
			escapeMarkdown(stat.Nickname),
			escapeMarkdown(stat.Stats))
	}

	return b.sendMessage(c.Message.Chat.ID, response)
}

func (b *Bot) handleTopGames(c *tgbotapi.Update) error {
	stats, err := b.playerService.GetTopByGames()
	if err != nil {
		log.Printf("Error getting top games: %v", err)
		return b.sendMessage(c.Message.Chat.ID, "Error fetching statistics")
	}

	if len(stats) == 0 {
		return b.sendMessage(c.Message.Chat.ID, "No statistics available")
	}

	response := "*Top players by games played:*\n\n"
	for i, stat := range stats {
		response += fmt.Sprintf("%d\\. *%s* \\- %s\n",
			i+1,
			escapeMarkdown(stat.Nickname),
			escapeMarkdown(stat.Stats))
	}

	return b.sendMessage(c.Message.Chat.ID, response)
}

func (b *Bot) handleTopCaptains(c *tgbotapi.Update) error {
	stats, err := b.playerService.GetTopCaptains()
	if err != nil {
		log.Printf("Error getting top captains: %v", err)
		return b.sendMessage(c.Message.Chat.ID, "Error fetching statistics")
	}

	if len(stats) == 0 {
		return b.sendMessage(c.Message.Chat.ID, "No captain statistics available")
	}

	response := "*Top captains by win rate:*\n\n"
	for i, stat := range stats {
		response += fmt.Sprintf("%d\\. *%s* \\- %s\n",
			i+1,
			escapeMarkdown(stat.Nickname),
			escapeMarkdown(stat.Stats))
	}

	return b.sendMessage(c.Message.Chat.ID, response)
}

func (b *Bot) handleTopRole(c *tgbotapi.Update) error {
	args := strings.Fields(c.Message.Text)[1:]
	if len(args) < 1 {
		return b.sendMessage(c.Message.Chat.ID, "Please specify a role: carry/mid/offlane/pos4/pos5\nExample: /top\\_role carry")
	}

	roleStr := strings.ToLower(args[0])
	stats, err := b.playerService.GetTopByRole(roleStr)
	if err != nil {
		log.Printf("Error getting top by role %s: %v", roleStr, err)
		return b.sendMessage(c.Message.Chat.ID, "Error fetching statistics")
	}

	if len(stats) == 0 {
		return b.sendMessage(c.Message.Chat.ID, fmt.Sprintf("No statistics available for role: %s", escapeMarkdown(roleStr)))
	}

	response := fmt.Sprintf("*Top %s players by win rate:*\n\n", escapeMarkdown(roleStr))
	for i, stat := range stats {
		response += fmt.Sprintf("%d\\. *%s* \\- %s\n",
			i+1,
			escapeMarkdown(stat.Nickname),
			escapeMarkdown(stat.Stats))
	}

	return b.sendMessage(c.Message.Chat.ID, response)
}

func (b *Bot) sendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		var err error
		switch update.Message.Command() {
		case "help":
			err = b.handleHelp(&update)
		case "top_winrate":
			err = b.handleTopWinRate(&update)
		case "top_games":
			err = b.handleTopGames(&update)
		case "top_captains":
			err = b.handleTopCaptains(&update)
		case "top_role":
			err = b.handleTopRole(&update)
		}

		if err != nil {
			log.Printf("Error handling command: %v", err)
		}
	}

	return nil
}
