package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/domain"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type tgUserKey struct{}

func Authorize(botToken string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			initData := r.Header.Get("Authorization")
			if initData == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}
			tgUser, err := validateTelegramInitData(initData, botToken)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid token or signature: %s", err), http.StatusUnauthorized)
				return
			}

			newCtx := context.WithValue(r.Context(), tgUserKey{}, tgUser)

			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}
func UserFromCtx(ctx context.Context) (domain.TelegramUser, error) {
	val, ok := ctx.Value(tgUserKey{}).(domain.TelegramUser)
	if !ok {
		return domain.TelegramUser{}, fmt.Errorf("wrong user id in ctx")
	}
	return val, nil
}

func validateTelegramInitData(initDataStr string, botToken string) (domain.TelegramUser, error) {
	// 1. Устанавливаем время жизни данных (например, 24 часа), чтобы предотвратить использование старых данных

	err := initdata.Validate(initDataStr, botToken, 0)
	if err != nil {
		return domain.TelegramUser{}, err // Здесь отсекаются неверные хэши и просроченные сессии
	}

	// 3. Если валидация успешна, парсим данные в готовую структуру
	data, err := initdata.Parse(initDataStr)
	if err != nil {
		return domain.TelegramUser{}, err
	}
	tgUser := domain.TelegramUser{
		ID:        int32(data.User.ID), //TODO int64
		Username:  data.User.Username,
		FirstName: data.User.FirstName,
		LastName:  data.User.LastName,
	}

	return tgUser, nil
}
