package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	authentication "telegram-api-with-go/internal/auth"
	"telegram-api-with-go/internal/bot"
	"telegram-api-with-go/internal/config"
	"telegram-api-with-go/internal/logger"
	"telegram-api-with-go/internal/session"
	"telegram-api-with-go/internal/telegram"
)

func main() {
	// Загружаем конфигурацию
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализируем логгер
	logger.InitLogger()
	log := logger.Log

	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Запуск приложения")

	// Инициализируем хранилище сессии
	sessionStorage := &session.MemorySession{}
	log.Debug("Инициализировано хранилище сессии")

	// Создаем Telegram клиент
	client := telegram.NewClient(sessionStorage)
	log.Debug("Создан Telegram клиент")

	// Создаем аутентификатор
	authenticator := &authentication.Auth{}

	// Создаем и запускаем Telegram клиент в отдельной горутине
	go func() {
		log.Info("Запуск Telegram клиента")
		if err := client.Run(ctx, authenticator); err != nil {
			log.Error("Ошибка Telegram клиента", "error", err)
			os.Exit(1)
		}
	}()

	// Создаем бота
	bot, err := bot.New(client)
	if err != nil {
		log.Error("Ошибка создания бота", "error", err)
		os.Exit(1)
	}
	log.Info("Бот создан успешно")

	// Запускаем бота в отдельной горутине
	go func() {
		log.Info("Запуск бота")
		if err := bot.Start(ctx); err != nil {
			log.Error("Ошибка запуска бота", "error", err)
			os.Exit(1)
		}
	}()

	// Ожидаем сигнала завершения
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Info("Получен сигнал завершения")
	log.Info("Завершение работы...")
} 