package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/PavlentiyGo/notification-service/services/api-gateway/internal/domain"
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
			if err != nil || tgUser.ID == 0 {
				http.Error(w, "Invalid token or signature", http.StatusUnauthorized)
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

func validateTelegramInitData(initData string, botToken string) (domain.TelegramUser, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return domain.TelegramUser{}, err
	}

	telegramHash := values.Get("hash")
	if telegramHash == "" {
		return domain.TelegramUser{}, fmt.Errorf("hash not found")
	}
	values.Del("hash")

	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, values.Get(k)))
	}
	dataCheckString := strings.Join(pairs, "\n")

	mac := hmac.New(sha256.New, []byte("WebAppData"))
	mac.Write([]byte(botToken))
	secretKey := mac.Sum(nil)

	signMac := hmac.New(sha256.New, secretKey)
	signMac.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(signMac.Sum(nil))
	equal := hmac.Equal([]byte(calculatedHash), []byte(telegramHash))
	if !equal {
		return domain.TelegramUser{}, fmt.Errorf("invalid token: %w", err)
	}

	userJson := values.Get("user")

	var tgUser domain.TelegramUser
	if err = json.Unmarshal([]byte(userJson), &tgUser); err != nil {
		return domain.TelegramUser{}, fmt.Errorf("failed to unmarshal user json: %w", err)
	}

	return tgUser, nil
}
