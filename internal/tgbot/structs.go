package tgbot

type UserInfo struct {
	UserId       int
	ChatId       int64
	UserName     string
	FirstName    string
	LastName     string
	LanguageCode string
	IsBot        bool
	IsBanned     bool
}
