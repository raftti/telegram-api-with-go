package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// handleUpdate обрабатывает обновления от Telegram
func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	b.log.Info("Получено сообщение",
		"user", update.Message.From.UserName,
		"text", update.Message.Text,
		"chat_id", update.Message.Chat.ID,
	)

	switch update.Message.Text {
	case "/spy":
		b.handleSpyCommand(ctx, update)
	case "/chats":
		b.handleChatsCommand(ctx, update)
	default:
		b.handleUnknownCommand(update)
	}
}

// handleSpyCommand обрабатывает команду /spy
func (b *Bot) handleSpyCommand(ctx context.Context, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if b.spy != nil {
		userID := b.spy.GetUserID()
		msg.Text = fmt.Sprintf("Теперь вы следите за пользователем %d.", userID)
		b.log.Info("Запуск слежения за пользователем",
			"user_id", userID,
			"initiator", update.Message.From.UserName,
		)
		go b.spy.StartSpying(ctx)
	} else {
		msg.Text = "Ошибка: сервис слежения не инициализирован"
		b.log.Error("Сервис слежения не инициализирован",
			"chat_id", update.Message.Chat.ID,
		)
	}

	if _, err := b.api.Send(msg); err != nil {
		b.log.Error("Ошибка отправки сообщения",
			"chat_id", update.Message.Chat.ID,
			"error", err,
		)
	}
}

// handleChatsCommand обрабатывает команду /chats
func (b *Bot) handleChatsCommand(ctx context.Context, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Запрашиваю список чатов...")
	b.api.Send(msg)

	b.log.Info("Запрос списка чатов",
		"user", update.Message.From.UserName,
		"chat_id", update.Message.Chat.ID,
	)

	chats, err := b.client.GetChats(ctx)
	if err != nil {
		msg.Text = "Ошибка: " + err.Error()
		b.log.Error("Ошибка получения списка чатов",
			"error", err,
			"chat_id", update.Message.Chat.ID,
		)
		b.api.Send(msg)
		return
	}

	if len(chats) == 0 {
		msg.Text = "Чаты не найдены."
		b.log.Info("Чаты не найдены",
			"chat_id", update.Message.Chat.ID,
		)
	} else {
		text := "Список чатов:\n"
		for i, chat := range chats {
			text += fmt.Sprintf("%d. %s\n", i+1, chat)
		}
		msg.Text = text
		b.log.Info("Отправка списка чатов",
			"count", len(chats),
			"chat_id", update.Message.Chat.ID,
		)
	}

	if _, err := b.api.Send(msg); err != nil {
		b.log.Error("Ошибка отправки списка чатов",
			"error", err,
			"chat_id", update.Message.Chat.ID,
		)
	}
}

// handleUnknownCommand обрабатывает неизвестные команды
func (b *Bot) handleUnknownCommand(update tgbotapi.Update) {
	b.log.Warn("Получена неизвестная команда",
		"command", update.Message.Text,
		"user", update.Message.From.UserName,
		"chat_id", update.Message.Chat.ID,
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда. Доступные команды: /spy, /chats")
	if _, err := b.api.Send(msg); err != nil {
		b.log.Error("Ошибка отправки сообщения о неизвестной команде",
			"error", err,
			"chat_id", update.Message.Chat.ID,
		)
	}
} 