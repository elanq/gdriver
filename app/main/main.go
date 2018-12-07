package main

import (
	"fmt"

	"github.com/elanq/gdriver"
)

func main() {
	//request flow
	//create wrapper with config
	//call auth URL
	//enter auth code from url
	//save credential
	//return client
	//use client to create drive instance
	wrapper, err := gdriver.NewDefaultWrapperWithConfig()
	if err != nil {
		panic(err)
	}

	if wrapper.AuthCode() == "" {
		url, err := gdriver.AuthURL(wrapper)
		if err != nil {
			panic(err)
		}

		fmt.Println("please open this url ", url)

		if err = wrapper.SetAuthCode(); err != nil {
			panic(err)
		}
	}

	api, err := wrapper.Drive()
	if err != nil {
		panic(err)
	}

	fmt.Println(api)

}
