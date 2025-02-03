package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"telegram-api-with-go/internal/config"
)

var (
	// Log глобальный логгер
	Log *slog.Logger
	// UseColors определяет, использовать ли цветной вывод
	UseColors = true
)

// Цвета для разных уровней логирования
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
)

// ColorHandler обработчик для цветного вывода логов
type ColorHandler struct {
	opts    slog.HandlerOptions
	writer  io.Writer
	attrs   []slog.Attr
	groups  []string
	mu      sync.Mutex
}

// NewColorHandler создает новый обработчик для цветного вывода
func NewColorHandler(w io.Writer, opts *slog.HandlerOptions) *ColorHandler {
	return &ColorHandler{
		opts:   *opts,
		writer: w,
	}
}

// Enabled проверяет, включен ли уровень логирования
func (h *ColorHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := h.opts.Level.Level()
	return level >= minLevel
}

// Handle обрабатывает запись лога
func (h *ColorHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Форматируем время
	timeStr := r.Time.Format("15:04:05.000")

	// Получаем цвет для уровня логирования
	levelColor := h.getColorOrEmpty(getLevelColor(r.Level))
	levelText := strings.ToUpper(r.Level.String())

	// Форматируем сообщение
	message := fmt.Sprintf("%s%s%s [%s%s%s] %s",
		h.getColorOrEmpty(colorGray), timeStr, h.getColorOrEmpty(colorReset),
		levelColor, levelText, h.getColorOrEmpty(colorReset),
		r.Message,
	)

	// Добавляем атрибуты с учетом групп
	if attrs := h.formatAttrs(r, h.attrs, h.groups); attrs != "" {
		message += fmt.Sprintf("\n%s%s%s",
			h.getColorOrEmpty(colorCyan),
			attrs,
			h.getColorOrEmpty(colorReset),
		)
	}

	// Добавляем информацию об источнике
	if h.opts.AddSource && r.PC != 0 {
		frames := runtime.CallersFrames([]uintptr{r.PC})
		frame, _ := frames.Next()
		source := fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
		message += fmt.Sprintf(" %s(%s)%s",
			h.getColorOrEmpty(colorPurple),
			source,
			h.getColorOrEmpty(colorReset),
		)
	}

	// Выводим сообщение
	_, err := fmt.Fprintln(h.writer, message)
	return err
}

// formatAttrs форматирует атрибуты с учетом групп
func (h *ColorHandler) formatAttrs(r slog.Record, attrs []slog.Attr, groups []string) string {
	if r.NumAttrs() == 0 && len(attrs) == 0 {
		return ""
	}

	allAttrs := make(map[string]interface{})
	
	// Добавляем атрибуты из записи
	r.Attrs(func(a slog.Attr) bool {
		allAttrs[a.Key] = a.Value.Any()
		return true
	})

	// Добавляем предустановленные атрибуты
	for _, attr := range attrs {
		allAttrs[attr.Key] = attr.Value.Any()
	}

	// Если есть группы, создаем вложенную структуру
	if len(groups) > 0 {
		nested := allAttrs
		for i := len(groups) - 1; i >= 0; i-- {
			nested = map[string]interface{}{
				groups[i]: nested,
			}
		}
		allAttrs = nested
	}

	// Форматируем в JSON
	data, err := json.MarshalIndent(allAttrs, "", "  ")
	if err != nil {
		return fmt.Sprintf("error formatting attributes: %v", err)
	}

	return string(data)
}

// getColorOrEmpty возвращает цветовой код или пустую строку, если цвета отключены
func (h *ColorHandler) getColorOrEmpty(color string) string {
	if UseColors {
		return color
	}
	return ""
}

// WithAttrs добавляет атрибуты к обработчику
func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := &ColorHandler{
		opts:    h.opts,
		writer:  h.writer,
		attrs:   append([]slog.Attr{}, h.attrs...),
		groups:  append([]string{}, h.groups...),
	}
	h2.attrs = append(h2.attrs, attrs...)
	return h2
}

// WithGroup добавляет группу к обработчику
func (h *ColorHandler) WithGroup(name string) slog.Handler {
	h2 := &ColorHandler{
		opts:    h.opts,
		writer:  h.writer,
		attrs:   append([]slog.Attr{}, h.attrs...),
		groups:  append([]string{}, h.groups...),
	}
	h2.groups = append(h2.groups, name)
	return h2
}

// getLevelColor возвращает цвет для уровня логирования
func getLevelColor(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return colorGray
	case slog.LevelInfo:
		return colorGreen
	case slog.LevelWarn:
		return colorYellow
	case slog.LevelError:
		return colorRed
	default:
		return colorReset
	}
}

// InitLogger инициализирует глобальный логгер
func InitLogger() {
	level := parseLogLevel(config.LogLevel)
	UseColors = os.Getenv("NO_COLOR") == "" // Отключаем цвета, если установлена NO_COLOR

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	handler := NewColorHandler(os.Stdout, opts)
	Log = slog.New(handler)
	slog.SetDefault(Log)
}

// parseLogLevel преобразует строковый уровень логирования в slog.Level
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
} 