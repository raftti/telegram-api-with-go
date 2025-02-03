package telegram

import (
	"context"
	"fmt"
	"log/slog"

	"telegram-api-with-go/internal/config"
	"telegram-api-with-go/internal/logger"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

// Client представляет клиент Telegram
type Client struct {
	client *telegram.Client
	log    *slog.Logger
}

// NewClient создает новый экземпляр клиента Telegram
func NewClient(sessionStorage telegram.SessionStorage) *Client {
	client := telegram.NewClient(config.APIID, config.APIHash, telegram.Options{
		SessionStorage: sessionStorage,
	})
	return &Client{
		client: client,
		log:    logger.Log,
	}
}

// Run запускает клиент Telegram
func (c *Client) Run(ctx context.Context, clientAuth auth.UserAuthenticator) error {
	return c.client.Run(ctx, func(ctx context.Context) error {
		status, err := c.client.Auth().Status(ctx)
		if err != nil {
			c.log.Error("Ошибка получения статуса авторизации", "error", err)
			return err
		}

		if !status.Authorized {
			c.log.Info("Требуется авторизация")
			flow := auth.NewFlow(clientAuth, auth.SendCodeOptions{})
			if err := c.client.Auth().IfNecessary(ctx, flow); err != nil {
				c.log.Error("Ошибка авторизации", "error", err)
				return err
			}
		}

		c.log.Info("Telegram клиент авторизован")
		<-ctx.Done()
		return ctx.Err()
	})
}

// GetChats получает список чатов
func (c *Client) GetChats(ctx context.Context) ([]string, error) {
	c.log.Info("Запрос списка чатов")
	var chats []string

	api := c.client.API()
	res, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetDate: 0,
		OffsetPeer: &tg.InputPeerEmpty{},
		Limit:      50,
		Hash:       0,
	})
	if err != nil {
		c.log.Error("Ошибка получения диалогов", "error", err)
		return nil, fmt.Errorf("ошибка MessagesGetDialogs: %w", err)
	}

	var dialogs *tg.MessagesDialogsSlice
	switch v := res.(type) {
	case *tg.MessagesDialogs:
		dialogs = &tg.MessagesDialogsSlice{
			Chats: v.Chats,
		}
	case *tg.MessagesDialogsSlice:
		dialogs = v
	default:
		c.log.Error("Неожиданный тип результата", "type", fmt.Sprintf("%T", res))
		return nil, fmt.Errorf("неожиданный тип результата %T", res)
	}

	for _, chatObj := range dialogs.Chats {
		var title string
		switch chat := chatObj.(type) {
		case *tg.Chat:
			title = fmt.Sprintf("Группа: %s", chat.Title)
		case *tg.ChatForbidden:
			title = fmt.Sprintf("Запрещенная группа: %s", chat.Title)
		}
		if title != "" {
			chats = append(chats, title)
		}
	}

	c.log.Info("Получен список чатов", "count", len(chats))
	return chats, nil
} 