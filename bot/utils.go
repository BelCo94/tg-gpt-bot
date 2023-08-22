package bot

import (
	"fmt"
	"strings"

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

var escapableChars = map[byte]bool{
	'_': true,
	'*': true,
	'[': true,
	']': true,
	'(': true,
	')': true,
	'~': true,
	'`': true,
	'>': true,
	'#': true,
	'+': true,
	'-': true,
	'=': true,
	'|': true,
	'{': true,
	'}': true,
	'.': true,
	'!': true,
}

func prepareTextForMarkdownV2(text string) string {
	var result strings.Builder
	var inSingleTick, inTripleTick bool
	for i := 0; i < len(text); i++ {
		switch {
		case !inSingleTick && i < len(text)-2 && text[i:i+3] == "```" && (i == len(text)-3 || text[i+3] != '`'):
			inTripleTick = !inTripleTick
			result.WriteString(text[i : i+3])
			i += 2
		case !inTripleTick && text[i] == '`' && (i == len(text)-1 || text[i+1] != '`'):
			inSingleTick = !inSingleTick
			result.WriteByte(text[i])
		case inTripleTick || inSingleTick:
			result.WriteByte(text[i])
		case escapableChars[text[i]]:
			result.WriteString("\\")
			result.WriteByte(text[i])
		default:
			result.WriteByte(text[i])
		}
	}
	return result.String()
}
