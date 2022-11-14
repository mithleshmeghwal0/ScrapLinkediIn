package login

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"

	"example.com/login/network"
	"example.com/login/network/http_client"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

const startURL = "https://www.linkedin.com"

const loginURL = "https://www.linkedin.com/uas/login-submit"

type linkedInLogin struct {
	ctx             context.Context
	log             *logrus.Entry
	client          network.Client
	sessionKey      string
	sessionPassword string
	csrfToken       string
}

func New(ctx context.Context, log *logrus.Entry, sessionkey, sessionPassword string) (*linkedInLogin, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cookiejar.New : %v", err)
	}

	return &linkedInLogin{
		ctx: ctx,
		client: http_client.New(&http.Client{
			Jar: jar,
		}, log),
		log:             log,
		sessionKey:      sessionkey,
		sessionPassword: sessionPassword,
	}, nil
}

func (l *linkedInLogin) setCSRFToken(token string) error {
	l.csrfToken = token
	return nil
}

func (l *linkedInLogin) Login() error {
	token, err := l.extractCSRFToken()
	if err != nil {
		return fmt.Errorf("extractCSRFToken: %v", err)
	}
	_ = l.setCSRFToken(token)
	return l.performLogin()
}

func (l *linkedInLogin) extractCSRFToken() (string, error) {

	res := l.client.GetBytes(startURL, nil)
	if res.Err != nil {
		return "", fmt.Errorf("GetBytes: %v", res.Err)
	}

	var input = []byte("input")
	var loginCsrfParam = []byte("loginCsrfParam")
	var name = []byte("name")
	var value = []byte("value")

	tokenizer := html.NewTokenizer(bytes.NewReader(res.Response))
	for {
		t := tokenizer.Next()
		if t == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				return "", fmt.Errorf(" tokenizer.Next(): io.EOF")
			}

			return "", fmt.Errorf("tokenizer.Next(): html.ErrorToken")
		}

		tag, hasAttr := tokenizer.TagName()
		if bytes.Equal(tag, input) {
			if hasAttr {
				var nameloginCsrfParamFound bool
				for {
					attrKey, attrValue, moreAttr := tokenizer.TagAttr()
					if bytes.Equal(attrKey, name) && bytes.Equal(attrValue, loginCsrfParam) {
						nameloginCsrfParamFound = true
						continue
					}
					if bytes.Equal(attrKey, value) && nameloginCsrfParamFound {
						return string(attrValue), nil
					}
					if !moreAttr {
						break
					}
				}
			}
		}
	}
}

func (l *linkedInLogin) performLogin() error {
	res := l.client.PostFormDataBytes(map[string]interface{}{
		"loginCsrfParam":   l.csrfToken,
		"session_key":      l.sessionKey,
		"session_password": l.sessionPassword,
	}, loginURL, nil)
	if res.Err != nil {
		return fmt.Errorf("PostFormDataBytes: %v", res.Err)
	}
	return nil
}
