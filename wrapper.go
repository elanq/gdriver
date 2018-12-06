package gdriver

import (
	"golang.org/x/oauth2"
	"io/ioutil"
)

const(
	DefaultCredential = "credentials.json"
)

type Wrapper interface {
	Token() (c *oauth2.Config, error)
	Client() (*http.Client, error)
}

type DefaultWrapper struct {
	c *oauth2.Config
}

func (d *DefaultWrapper) Token() (c *oauth2.Config, error) {
	f, err := ioutil.ReadFile(DefaultCredential)
	if err != nil {
		return nil, err
	}

	c, err := google.ConfigFromJSON(f, drive.DriveMetadataReadonlyScope)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (d *DefaultWrapper) Client() (*http.Client, error) {
	return nil, nil
}
