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
/prokuror \- Show prokuror stats
/happy\_birthday \- Special birthday wishes for Даник

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

func (b *Bot) handleProkuror(c *tgbotapi.Update) error {
	stats, err := b.playerService.GetProkurorStats()
	if err != nil {
		log.Printf("Error getting prokuror stats: %v", err)
		return b.sendMessage(c.Message.Chat.ID, "Error fetching statistics")
	}

	response := "🚨 *ВЕРХОВНЫЙ ПРОКУРОР* 🚨\n\n"
	response += fmt.Sprintf("👮‍♂️ *%s* 👮‍♂️\n", escapeMarkdown(stats.Nickname))
	response += fmt.Sprintf("🚔 *Статистика:* %s 🚓\n", escapeMarkdown(stats.Stats))
	response += "\n🏛️ *Закон и порядок* ⚖️\n"
	response += "🚨 *Справедливость восторжествует* 🚨"

	return b.sendMessage(c.Message.Chat.ID, response)
}

func (b *Bot) handleHappyBirthday(c *tgbotapi.Update) error {
	emojis := "🎂🎁🎉🎊🥳🍾🥂🎇✨"
	nickname := "даня тапок"

	response := fmt.Sprintf("*С Днем Рождения, Даник* %s\n\n", emojis)
	response += "Братан, от души желаю:\n"
	response += "💪 Силы и мощи как у быка\n"
	response += "💰 Бабла, чтобы некуда было девать\n"
	response += "🏆 Побед на всех фронтах\n"
	response += "🔥 Огня в глазах и страсти в сердце\n"
	response += "🎮 И конечно же побед в играх\n\n"

	response += fmt.Sprintf("*Статистика игрока %s:*\n", escapeMarkdown(nickname))
	response += "🏅 *Win Rate:* 100% \\(16/16 побед\\)\n"
	response += "👑 *MVP:* в каждой игре\n"
	response += "🚀 *Статус:* Абсолютная легенда\n\n"
	response += "С днюхой, красавчик 🍻"

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
	u.Timeout = 10

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
		case "prokuror":
			err = b.handleProkuror(&update)
		case "happy_birthday":
			err = b.handleHappyBirthday(&update)
		}

		if err != nil {
			log.Printf("Error handling command: %v", err)
		}
	}

	return nil
}
