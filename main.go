package main

import (
	"context"
	"fmt"
	"os"

	"example.com/login/log"
	"example.com/login/login"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	sessionKey      = os.Getenv("SESSION_KEY")
	sessionPassword = os.Getenv("SESSION_PASSWORD")
)

func main() {
	var l *logrus.Entry
	var f *os.File
	var err error
	if os.Getenv("LOG_FILE") == "1" {
		f, err = os.OpenFile(fmt.Sprintf("%s.json", uuid.NewString()), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		// don't forget to close it
		defer f.Close()
		l = log.NewWithFile(f)
	} else {
		l = log.New()
	}

	ctx := context.Background()

	login, err := login.New(ctx, l, sessionKey, sessionPassword)
	if err != nil {
		l.Error("login.New() Failed")
		os.Exit(1)
	}
	err = login.Login()
	if err != nil {
		l.Error("Login Failed")
		os.Exit(1)
	}
	l.Info("Login Successful")
}
