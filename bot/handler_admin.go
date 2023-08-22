package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BelCo94/tg-gpt-bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *Bot) isAdminCommand(msg *tgbotapi.Message) bool {
	return msg.From.ID == int(bot.config.Base.AdminID) && msg.IsCommand() && msg.Command() == "admin"
}

func (bot *Bot) handleAdminCommand(msg *tgbotapi.Message) error {
	text := bot.performAdminCommand(msg.CommandArguments())
	newMsg := tgbotapi.NewMessage(bot.config.Base.AdminID, text)
	_, err := bot.tgbot.Send(newMsg)
	return err
}

func (bot *Bot) performAdminCommand(cmd string) string {
	params := strings.Split(cmd, " ")
	if len(params) != 2 {
		return "Wrong number of pararmeters"
	}
	if params[0] != "enable" && params[0] != "disable" {
		return "Wrong action"
	}
	target, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return "Wrong id"
	}
	chat := &storage.Chat{
		ID:         target,
		Name:       "", // todo: get current name
		IsDisabled: params[0] == "disable",
	}
	if err := bot.storage.CreateOrUpdateChat(chat); err != nil {
		return "failed"
	}
	if target > 0 {
		user := &storage.User{
			ID:         target,
			Name:       "", //todo: get current name
			IsDisabled: params[0] == "disable",
		}
		if err := bot.storage.CreateOrUpdateUser(user); err != nil {
			return "failed"
		}
	}
	return "done"
}

func (bot *Bot) handleNewSource(msg *tgbotapi.Message) error {
	text := fmt.Sprintf("I was mentioned in the new chat:\nChat ID: `%d`\nChat Title: `%s`\nUser ID: `%d`\nUser name: `%s`", msg.Chat.ID, chatTitle(msg.Chat), msg.From.ID, msg.From.UserName)
	newMsg := tgbotapi.NewMessage(bot.config.Base.AdminID, text)
	newMsg.ParseMode = "MarkdownV2"
	_, err := bot.tgbot.Send(newMsg)
	return err
}
