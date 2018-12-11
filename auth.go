package gdriver

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

//NewClient creates http client connection to google drive API
func NewClient(w Wrapper) (*http.Client, error) {
	token, err := FileToken(DefaultToken)
	if err != nil {
		token, err = WebToken(w)
		if err != nil {
			return nil, err
		}
		if err = writeCredential(DefaultToken, token); err != nil {
			return nil, err
		}
	}

	conf, err := w.OauthConfig()
	if err != nil {
		return nil, err
	}

	if err = w.SetToken(token); err != nil {
		return nil, err
	}

	client := conf.Client(context.Background(), token)
	return client, nil
}

//AuthURL get authentication URL
func AuthURL(w Wrapper) (string, error) {
	config, err := w.OauthConfig()
	if err != nil {
		return "", err
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL, nil
}

//WebToken accepts authcode and convert it as oauth token
func WebToken(w Wrapper) (*oauth2.Token, error) {
	config, err := w.OauthConfig()
	if err != nil {
		return nil, err
	}

	if w.AuthCode() == "" {
		return nil, errors.New("Empty authcode, please set it first")
	}

	return config.Exchange(context.TODO(), w.AuthCode())
}

func FileToken(filename string) (*oauth2.Token, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func writeCredential(path string, token *oauth2.Token) error {
	log.Println("saving credential to ", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	return err
}
