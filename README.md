# Telegram Bot with User Tracking

Telegram бот с возможностью отслеживания активности пользователей и просмотра списка чатов.

## Функциональность

- Отслеживание онлайн-статуса пользователя
- Просмотр списка чатов
- Простой и понятный интерфейс
- Структурированное логирование

## Команды

- `/spy` - начать отслеживание пользователя
- `/chats` - получить список чатов

## Установка

1. Клонируйте репозиторий:

```bash
git clone https://github.com/yourusername/telegram-api-with-go.git
cd telegram-api-with-go
```

2. Установите зависимости:

```bash
go mod download
```

3. Настройте окружение:

   - Скопируйте файл `.env.example` в `.env`:
     ```bash
     cp .env.example .env
     ```
   - Отредактируйте `.env` файл, указав свои значения:
     - `TELEGRAM_API_ID` и `TELEGRAM_API_HASH` можно получить на https://my.telegram.org
     - `TELEGRAM_BOT_TOKEN` можно получить у @BotFather в Telegram
     - `DEFAULT_SPY_USER_ID` - ID пользователя для отслеживания
     - `LOG_LEVEL` - уровень логирования (debug, info, warn, error)

4. Запустите бота:

```bash
go run cmd/bot/main.go
```

## Структура проекта

```
telegram-api-with-go/
├── cmd/
│   └── bot/           # Точка входа в приложение
├── internal/
│   ├── auth/          # Аутентификация
│   ├── session/       # Управление сессией
│   ├── bot/           # Логика бота
│   ├── telegram/      # Работа с Telegram API
│   ├── logger/        # Логирование
│   └── config/        # Конфигурация
└── pkg/
    └── utils/         # Вспомогательные функции
```

## Конфигурация

Проект использует переменные окружения для конфигурации. Все настройки хранятся в файле `.env`:

```env
# Telegram API credentials
TELEGRAM_API_ID=your_api_id
TELEGRAM_API_HASH=your_api_hash
TELEGRAM_BOT_TOKEN=your_bot_token

# User settings
DEFAULT_SPY_USER_ID=target_user_id

# File paths
SESSION_FILE=session.data

# Logging
LOG_LEVEL=debug  # debug, info, warn, error
```

## Логирование

Бот использует структурированное логирование с помощью `slog`. Логи выводятся в формате JSON и содержат:

- Уровень логирования
- Временную метку
- Источник (файл и строку)
- Контекстную информацию

## Требования

- Go 1.21 или выше
- Telegram API ID и Hash
- Telegram Bot Token

## Лицензия

MIT
