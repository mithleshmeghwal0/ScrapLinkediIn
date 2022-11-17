package login

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"

	"example.com/login/network"
	"example.com/login/network/http_client"
	"golang.org/x/net/html"
)

const startURL = "https://www.linkedin.com"

const loginURL = "https://www.linkedin.com/uas/login-submit"

type linkedInLogin struct {
	ctx             context.Context
	file            io.Writer
	client          network.Client
	sessionKey      string
	sessionPassword string
	csrfToken       string
}

func New(ctx context.Context, f io.Writer, sessionkey, sessionPassword string) (*linkedInLogin, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cookiejar.New : %v", err)
	}

	return &linkedInLogin{
		ctx: ctx,
		client: http_client.New(&http.Client{
			Jar: jar,
		}, f),
		file:            f,
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
	return l.performLogin(0)
}

func (l *linkedInLogin) extractCSRFToken() (string, error) {

	res := l.client.GetBytes(startURL,
		map[string]string{
			"Accept":     "text/html",
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
		})
	if res.Err != nil {
		return "", fmt.Errorf("GetBytes: %v", res.Err)
	}
	return searchloginCsrfParam(res.Response)
}

func searchloginCsrfParam(res []byte) (string, error) {

	var input = []byte("input")
	var loginCsrfParam = []byte("loginCsrfParam")
	var name = []byte("name")
	var value = []byte("value")

	tokenizer := html.NewTokenizer(bytes.NewReader(res))
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

func (l *linkedInLogin) performLogin(try int) error {
	if try == 2 {
		return errors.New("performLogin: max retries")
	}
	res := l.client.PostFormDataBytes(map[string]interface{}{
		"loginCsrfParam":   l.csrfToken,
		"session_key":      l.sessionKey,
		"session_password": l.sessionPassword,
	}, loginURL,
		map[string]string{
			"Accept":          "text/html",
			"Accept-Encoding": "gzip, deflate, br",
			"Connection":      "keep-alive",
			"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
		})
	if res.Err != nil {
		fmt.Println("here")
		token, err := searchloginCsrfParam(res.Response)
		if err != nil {
			fmt.Println("token not found")
			return fmt.Errorf("PostFormDataBytes: %v", res.Err)
		}
		l.setCSRFToken(token)
		try = try + 1
		l.performLogin(try)
	}
	fmt.Println("here 2")
	return nil
}
