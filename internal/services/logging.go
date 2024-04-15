package services

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type LogMsg struct {
	Event string
	Level string
	Time  time.Time
}

var Logging = logrus.New()

func InitLogger() error {
	Logging.SetFormatter(&logrus.JSONFormatter{})
	Logging.SetLevel(logrus.TraceLevel)
	fileLog, err := os.OpenFile("./app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logging.Info(fmt.Errorf("Logging:init:" + err.Error()))
	} else {
		Logging.SetOutput(fileLog)
	}
	Logging.Info("Bot started is normally")

	return nil
}
