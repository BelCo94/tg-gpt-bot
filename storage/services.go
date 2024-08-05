package storage

import (
	"fmt"
)

func (storage *Storage) GetMessageChain(message *Message, chainLength uint16) []*Message {
	messages := make([]*Message, 0, chainLength)
	messages = append(messages, message)
	replyToID := message.ReplyToID
	for len(messages) < int(chainLength) && replyToID != -1 {
		var msg Message
		if err := storage.DB.Where("ID = ?", replyToID).First(&msg).Error; err != nil {
			fmt.Printf("Message with ID = %d not found in the DB\n", replyToID)
			break
		}
		if msg.IsImageGeneration {
			break // stop chain on image generation message
		}
		messages = append(messages, &msg)
		replyToID = msg.ReplyToID
	}
	return messages
}

func (storage *Storage) GetChat(chatId int64) *Chat {
	var obj Chat
	if err := storage.DB.Where("ID = ?", chatId).First(&obj).Error; err != nil {
		return nil
	}
	return &obj
}

func (storage *Storage) GetUser(userId int64) *User {
	var obj User
	if err := storage.DB.Where("ID = ?", userId).First(&obj).Error; err != nil {
		return nil
	}
	return &obj
}

func (storage *Storage) CreateOrUpdateUser(user *User) error {
	return storage.DB.Save(&user).Error
}

func (storage *Storage) CreateOrUpdateChat(chat *Chat) error {
	return storage.DB.Save(&chat).Error
}

func (storage *Storage) CreateMessage(message *Message) error {
	return storage.DB.Create(message).Error
}

func (storage *Storage) GetLastMessageInChat(chat *Chat) *Message {
	var obj Message
	if err := storage.DB.Where("chat_id = ?", chat.ID).Order("ID desc").Limit(1).First(&obj).Error; err != nil {
		return nil
	}
	return &obj
}
