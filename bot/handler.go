package bot

import (
	"fmt"

	"github.com/BelCo94/tg-gpt-bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *Bot) messageHandler(msg *tgbotapi.Message) error {
	if msg.Text == "" {
		fmt.Printf("Skip service (empty) message in chat %s from %s\n", formatChat(msg.Chat), formatUser(msg.From))
		return nil
	}
	if bot.isAdminCommand(msg) {
		fmt.Printf("Handle admin command in chat %s\n", formatChat(msg.Chat))
		return bot.handleAdminCommand(msg)
	}
	chat := bot.storage.GetChat(msg.Chat.ID)
	user := bot.storage.GetUser(int64(msg.From.ID))
	if (chat != nil && chat.IsDisabled) || (user != nil && user.IsDisabled) {
		fmt.Printf("Skip message in chat %s from %s\n", formatChat(msg.Chat), formatUser(msg.From))
		return nil
	}
	if chat == nil || user == nil {
		return bot.handleNewSource(msg)
	}
	fmt.Printf("Handle message in chat %s from %s\n", formatChat(msg.Chat), formatUser(msg.From))
	storedMessage := &storage.Message{
		ID:        int64(msg.MessageID),
		ChatID:    msg.Chat.ID,
		UserID:    user.ID,
		Text:      msg.Text,
		ReplyToID: -1,
	}
	switch {
	case msg.IsCommand() && msg.Command() == "ask":
		storedMessage.Text = msg.CommandArguments()
	case msg.ReplyToMessage != nil:
		storedMessage.ReplyToID = int64(msg.ReplyToMessage.MessageID)
	case msg.Chat.IsPrivate():
		previousMessage := bot.storage.GetLastMessageInChat(chat)
		if previousMessage != nil {
			storedMessage.ReplyToID = previousMessage.ID
		}
	}
	if err := bot.storage.CreateMessage(storedMessage); err != nil {
		return err
	}
	messages := bot.storage.GetMessageChain(storedMessage, bot.config.Base.ChatHistorySize)
	return bot.askOpenAI(messages)
}
