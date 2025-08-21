package http

import "context"

type TelegramClient interface {
	SendMessage(ctx context.Context, chatID, text string) error
}
