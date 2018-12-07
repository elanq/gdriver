package gdriver

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v2"
)

const (
	DefaultCredential = "credentials.json"
	DefaultToken      = "token.json"
)

type Wrapper interface {
	OauthConfig() (*oauth2.Config, error)
	Client() (*http.Client, error)
	AuthCode() string
}

type DefaultWrapper struct {
	conf     *oauth2.Config
	client   *http.Client
	authCode string
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

func (d *DefaultWrapper) SetAuthCode() error {
	var authCode string

	fmt.Println("enter your oauth code from given link before:")
	if _, err := fmt.Scan(&authCode); err != nil {
		return err
	}
	d.authCode = authCode

	return nil
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
