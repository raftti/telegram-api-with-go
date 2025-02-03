package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"telegram-api-with-go/internal/logger"

	"github.com/gotd/td/tg"
)

// SpyService представляет сервис для слежения за пользователем
type SpyService struct {
	client *Client
	userID int64
	log    *slog.Logger
	lastOnline uint8
}

type StoredUserEvent struct {
	UserID     int64 `json:"user_id"`
	LastOnline int64 `json:"last_online"`
}

// NewSpyService создает новый сервис слежения
func NewSpyService(client *Client, userID int64) *SpyService {
	return &SpyService{
		client: client,
		userID: userID,
		log:    logger.Log,
	}
}

// StartSpying начинает слежение за пользователем
func (s *SpyService) StartSpying(ctx context.Context) {
	s.log.Info("Начало слежения за пользователем", "user_id", s.userID)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.log.Info("Завершение слежения за пользователем", "user_id", s.userID)
			return
		case <-ticker.C:
			s.checkUserStatus(ctx)
		}
	}
}

// checkUserStatus проверяет статус пользователя
func (s *SpyService) checkUserStatus(ctx context.Context) {
	api := s.client.client.API()
	
	// Получаем информацию о пользователе
	user, err := api.UsersGetUsers(ctx, []tg.InputUserClass{
		&tg.InputUser{
			UserID: s.userID,
		},
	})

	if err != nil {
		s.log.Error("Ошибка получения информации о пользователе", 
			"user_id", s.userID,
			"error", err,
		)
		return
	}

	userStatuses, _ := user[0].AsNotEmpty()
	isOffline, ok := userStatuses.Status.(*tg.UserStatusOffline)

	if !ok {
		s.log.Info("Пользователь онлайн", "user_id", s.userID)
	} else {
		s.lastOnline = uint8(isOffline.GetWasOnline())
		fmt.Println("Пользователь offline", s.lastOnline)
		go saveStatusToFile(s.userID, int64(isOffline.GetWasOnline()))
	}
}

// GetUserID возвращает ID отслеживаемого пользователя
func (s *SpyService) GetUserID() int64 {
	return s.userID
} 

func getLastStatus(fileName string) (*StoredUserEvent, error) {
	file, err := os.Open(fileName)
	if err != nil {
		// Если файл не существует, возвращаем nil без ошибки.
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var event StoredUserEvent
	if err := json.NewDecoder(file).Decode(&event); err != nil {
		return nil, err
	}
	return &event, nil
}

// updateLastStatus обновляет файл с последним статусом.
func updateLastStatus(fileName string, event StoredUserEvent) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(event)
}

// saveStatusToFile добавляет новую запись в историю статусов,
// только если новый lastOnline отличается от предыдущего.
func saveStatusToFile(userID int64, lastOnline int64) {
	const (
		historyFile   = "user_status.json"
		lastStatusFile = "last_status.json"
	)

	// Проверяем последний статус из отдельного файла.
	lastStatus, err := getLastStatus(lastStatusFile)
	if err != nil {
		panic(err)
	}

	// Если последний статус уже соответствует новому, выходим.
	if lastStatus != nil && lastStatus.LastOnline == lastOnline {
		return
	}

	// Формируем новый статус.
	newEvent := StoredUserEvent{
		UserID:     userID,
		LastOnline: lastOnline,
	}

	// Открываем файл истории для дозаписи.
	history, err := os.OpenFile(historyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer history.Close()

	// Записываем новый статус в файл истории.
	// Запишем как одну JSON-строку с переводом строки в конце.
	eventData, err := json.Marshal(newEvent)
	if err != nil {
		panic(err)
	}
	if _, err := history.Write(eventData); err != nil {
		panic(err)
	}
	if _, err := history.Write([]byte("\n")); err != nil {
		panic(err)
	}

	// Обновляем файл с последним статусом.
	if err := updateLastStatus(lastStatusFile, newEvent); err != nil {
		panic(err)
	}
}

