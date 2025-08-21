package telegram

import (
	"context"
	"strconv"

	"github.com/go-telegram/bot"
)

type Client struct {
	bot *bot.Bot
}

// NewClient создает новый клиент для Telegram Bot API
func NewClient(bot *bot.Bot) *Client {
	return &Client{
		bot: bot,
	}
}

// SendMessage отправляет сообщение в указанный чат
func (c *Client) SendMessage(ctx context.Context, chatID, text string) error {
	// Конвертируем строковый chatID в int64
	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return err
	}

	_, err = c.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatIDInt,
		Text:      text,
		ParseMode: "HTML",
	})

	return err
}
