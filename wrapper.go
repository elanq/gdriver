package gdriver

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v2"
)

const (
	DefaultCredential = "credentials.json"
	DefaultToken      = "token.json"
)

var (
	ErrTokenNil = errors.New("Token has nil value, please authenticate first")
)

type Wrapper interface {
	OauthConfig() (*oauth2.Config, error)
	Client() (*http.Client, error)
	AuthCode() string
	Drive() (*drive.Service, error)
	SetToken(*oauth2.Token) error
	Token() (*oauth2.Token, bool)
}

type DefaultWrapper struct {
	conf     *oauth2.Config
	client   *http.Client
	authCode string
	drive    *drive.Service
	token    *oauth2.Token
}

func NewDefaultWrapperWithConfig() (*DefaultWrapper, error) {
	wrapper := &DefaultWrapper{}
	conf, err := wrapper.newConfig()
	if err != nil {
		return nil, err
	}

	wrapper.conf = conf

	return wrapper, nil

}

func (d *DefaultWrapper) SetToken(token *oauth2.Token) error {
	if token == nil {
		return ErrTokenNil
	}

	d.token = token
	return nil
}

func (d *DefaultWrapper) SetAuthCode() error {
	var authCode string

	fmt.Println("enter your oauth code from given link before:")
	if _, err := fmt.Scan(&authCode); err != nil {
		return err
	}
	d.authCode = authCode

	return nil
}

func (d *DefaultWrapper) Token() (*oauth2.Token, bool) {
	if d.token != nil {
		log.Println("retrieving token by wrapper object")
		return d.token, true
	}

	if token, err := FileToken(DefaultToken); err == nil {
		log.Println("retrieving token by file")
		return token, true
	}

	if token, err := WebToken(d); err == nil {
		log.Println("retrieving token by web token")
		return token, true
	}

	return d.token, false
}

func (d *DefaultWrapper) AuthCode() string {
	return d.authCode
}

func (d *DefaultWrapper) OauthConfig() (*oauth2.Config, error) {
	if d.conf != nil {
		return d.conf, nil
	}

	return d.newConfig()
}

func (d *DefaultWrapper) Client() (*http.Client, error) {
	if d.client != nil {
		return d.client, nil
	}

	client, err := NewClient(d)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (d *DefaultWrapper) Drive() (*drive.Service, error) {
	client, err := d.Client()
	if err != nil {
		return nil, err
	}

	if d.drive != nil {
		return d.drive, nil
	}

	return drive.New(client)
}

func (d *DefaultWrapper) newConfig() (*oauth2.Config, error) {
	f, err := ioutil.ReadFile(DefaultCredential)
	if err != nil {
		return nil, err
	}

	c, err := google.ConfigFromJSON(f, drive.DriveMetadataReadonlyScope)
	if err != nil {
		return nil, err
	}

	d.conf = c

	return c, nil
}
