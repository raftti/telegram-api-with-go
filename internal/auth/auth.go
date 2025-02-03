package authentication

import (
	"context"
	"fmt"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

// Auth реализует интерфейс аутентификации пользователя
type Auth struct{}

// Phone запрашивает номер телефона
func (a Auth) Phone(ctx context.Context) (string, error) {
	var phone string
	fmt.Print("Введите номер телефона: ")
	_, err := fmt.Scan(&phone)
	return phone, err
}

// Code запрашивает код подтверждения
func (a Auth) Code(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
	var code string
	fmt.Print("Введите код из Telegram: ")
	_, err := fmt.Scan(&code)
	return code, err
}

// Password запрашивает пароль
func (a Auth) Password(ctx context.Context) (string, error) {
	var password string
	fmt.Print("Введите пароль (если требуется): ")
	_, err := fmt.Scan(&password)
	return password, err
}

// SignUp запрашивает данные для регистрации
func (a Auth) SignUp(ctx context.Context) (auth.UserInfo, error) {
	var firstName, lastName string
	fmt.Print("Введите имя: ")
	_, _ = fmt.Scan(&firstName)
	fmt.Print("Введите фамилию: ")
	_, _ = fmt.Scan(&lastName)
	return auth.UserInfo{FirstName: firstName, LastName: lastName}, nil
}

// AcceptTermsOfService запрашивает подтверждение условий использования
func (a Auth) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	fmt.Println("Примите условия:", tos.Text)
	fmt.Print("Принять? (yes/no): ")
	var response string
	_, _ = fmt.Scan(&response)
	if response == "yes" {
		return nil
	}
	return fmt.Errorf("пользователь не принял условия")
} 