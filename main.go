package main

import (
	"fmt"
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	TUID string 	`json:"typeform_uid"`
	TKey string `json:"typeform_key"`
	EmailField string `json:"typeform_email_field"`
	NameField string `json:"typeform_name_field"`
	LastNameField string `json:"typeform_lastname_field"`
	SlackChannel string `json:"slack_channel_name"`
	SlackToken string `json:"slack_token_name"`
}

func main()  {
	dat, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("failed to open file config.json")
	}
	var config Config
	err = json.Unmarshal(dat, &config)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(config.SlackToken)

	go func() {
		for {

			fmt.Println("yay")
			time.Sleep(1 * time.Second)
		}
	}()
	startHttpServer()
}

func startHttpServer(){
	http.Handle("/", http.FileServer(http.Dir("public")))

	http.ListenAndServe(":3000", nil)

}