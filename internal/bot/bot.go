package bot

import (
	"context"
	"log/slog"

	"telegram-api-with-go/internal/config"
	"telegram-api-with-go/internal/logger"
	"telegram-api-with-go/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Bot представляет Telegram бота
type Bot struct {
	api    *tgbotapi.BotAPI
	client *telegram.Client
	spy    *telegram.SpyService
	log    *slog.Logger
}

// New создает нового бота
func New(client *telegram.Client) (*Bot, error) {
	log := logger.Log

	api, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Error("Ошибка создания Telegram API", "error", err)
		return nil, err
	}

	spy := telegram.NewSpyService(client, config.DefaultSpyUserID)

	return &Bot{
		api:    api,
		client: client,
		spy:    spy,
		log:    log,
	}, nil
}

// Start запускает бота
func (b *Bot) Start(ctx context.Context) error {
	b.log.Info("Запуск бота", "username", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.api.GetUpdatesChan(u)
	if err != nil {
		b.log.Error("Ошибка получения канала обновлений", "error", err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			b.log.Info("Завершение работы бота")
			return ctx.Err()
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			b.handleUpdate(ctx, update)
		}
	}
} 