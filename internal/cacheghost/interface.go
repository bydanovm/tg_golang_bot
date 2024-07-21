package cacheghost

import "github.com/mbydanov/tg_golang_bot/internal/database"

type UserInfo interface {
	GetUserInfo(usrId int) database.Users
	GetTracking(usrId int) database.TrackingCrypto
	GetLimits(usrId int) database.Limits
	URLockU() bool
	URUnlock() bool
}

func GetUserName(ui UserInfo) (userName string) {
	return userName
}

func GetUserId(ui UserInfo) (userId int) {
	return userId
}

func GetChatId(ui UserInfo) (chatId int) {
	return chatId
}
