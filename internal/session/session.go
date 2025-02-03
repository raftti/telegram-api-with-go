package session

import (
	"context"
	"os"
	"sync"

	"telegram-api-with-go/internal/config"

	"github.com/gotd/td/session"
)

// MemorySession - потокобезопасное хранение сессии
type MemorySession struct {
	mux  sync.RWMutex
	data []byte
}

// LoadSession загружает данные сессии из файла
func (s *MemorySession) LoadSession(ctx context.Context) ([]byte, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	data, err := os.ReadFile(config.SessionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, session.ErrNotFound
		}
		return nil, err
	}

	s.data = data
	return data, nil
}

// StoreSession сохраняет данные сессии в файл
func (s *MemorySession) StoreSession(ctx context.Context, data []byte) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	err := os.WriteFile(config.SessionFile, data, 0600)
	if err != nil {
		return err
	}

	s.data = data
	return nil
} 