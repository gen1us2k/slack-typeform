package main

import (
	"fmt"
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	"strings"
	"net/url"
	"log"
)

type Config struct {
	TUID string 	`json:"typeform_uid"`
	TKey string `json:"typeform_key"`
	EmailField string `json:"typeform_email_field"`
	NameField string `json:"typeform_name_field"`
	LastNameField string `json:"typeform_lastname_field"`
	SlackChannel string `json:"slack_channel_name"`
	SlackToken string `json:"slack_token_name"`
	Interval int64 `json:"check_interval"`
	ListenPort string `json:"http_port"`
	IpAddr string `json:"bind_addr"`
}

type Response struct {
	Id string
	Answers map[string] string
}
type Answer struct {
	Responses []Response
}


func inviteAll(config Config){
	log.Println("inviteAll")
	typeFormUrl := fmt.Sprintf("https://api.typeform.com/v0/form/%s?key=%s&completed=true&since=%d", config.TUID, config.TKey, time.Now().Unix() - config.Interval)
	res, err := http.Get(typeFormUrl)

	if err != nil {
		log.Panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if err != nil {
			log.Panic(err)
		}
	}

	var answer Answer
	json.Unmarshal(body, &answer)

	for _, user := range answer.Responses {
		log.Println("Sending invite to %s %s with email %s", user.Answers[config.NameField], user.Answers[config.LastNameField], user.Answers[config.EmailField])
		slackUrl := fmt.Sprintf("https://%s.slack.com/api/users.admin.invite", config.SlackChannel)

		invite := url.Values{}
		invite.Add("email", user.Answers[config.EmailField])
		invite.Add("first_name", user.Answers[config.NameField])
		invite.Add("last_name", user.Answers[config.LastNameField])
		invite.Add("token", config.SlackToken)
		invite.Add("set_active", "true")

		req, err := http.NewRequest("POST", slackUrl, strings.NewReader(invite.Encode()))
		if err != nil {
			log.Panic(err)
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			log.Panic(err)
		}

		defer resp.Body.Close()
	}
}

func main()  {
	dat, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Panic("Failed to open file config.json")
	}

	var config Config
	err = json.Unmarshal(dat, &config)

	if err != nil {
		log.Panic(err)
	}

	go func() {
		for _ = range time.Tick(time.Duration(config.Interval) * time.Second) {
			inviteAll(config)
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("public")))
	listenAddr := fmt.Sprintf("%s:%s", config.IpAddr, config.ListenPort)
	fmt.Println(listenAddr)
	http.ListenAndServe(listenAddr, nil)
}