package storage

type User struct {
	ID         int64 `gorm:"primaryKey"`
	Name       string
	IsDisabled bool
}

type Chat struct {
	ID         int64 `gorm:"primaryKey"`
	Name       string
	IsDisabled bool
}

type Message struct {
	ID        int64 `gorm:"primaryKey"`
	ChatID    int64 `gorm:"primaryKey"`
	UserID    int64
	Text      string
	ReplyToID int64
}
