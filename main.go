package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"example.com/login/login"
	"github.com/google/uuid"
)

var (
	sessionKey      = os.Getenv("SESSION_KEY")
	sessionPassword = os.Getenv("SESSION_PASSWORD")
)

func main() {
	var f *os.File
	var err error
	files, err := filepath.Glob("*.json")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
	f, err = os.OpenFile(fmt.Sprintf("%s.json", uuid.NewString()), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	// don't forget to close it
	defer f.Close()

	ctx := context.Background()

	login, err := login.New(ctx, f, sessionKey, sessionPassword)
	if err != nil {
		fmt.Println("login.New() Failed")
		os.Exit(1)
	}
	err = login.Login()
	if err != nil {
		fmt.Println("Login Failed")
		os.Exit(1)
	}
	fmt.Println("Login Successful")
}
