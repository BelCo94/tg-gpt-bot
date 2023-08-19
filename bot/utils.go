package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func chatTitle(chat *tgbotapi.Chat) string {
	if chat.Title != "" {
		return chat.Title
	}
	return chat.UserName
}

func formatUser(user *tgbotapi.User) string {
	return fmt.Sprintf("%s (id=%d)", user.UserName, user.ID)
}

func formatChat(chat *tgbotapi.Chat) string {
	return fmt.Sprintf("%s (id=%d)", chatTitle(chat), chat.ID)
}
