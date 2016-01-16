package main

import (
	"fmt"
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	"html/template"
	"strings"
	"net/url"
	"log"
)

// Config used for service configuration
type Config struct {
	TUID          string 	`json:"typeform_uid"`
	TKey          string 	`json:"typeform_key"`
	EmailField    string 	`json:"typeform_email_field"`
	NameField     string 	`json:"typeform_name_field"`
	LastNameField string 	`json:"typeform_lastname_field"`
	SlackChannel  string 	`json:"slack_channel_name"`
	SlackToken    string 	`json:"slack_token_name"`
	Interval      int64 	`json:"check_interval"`
	ListenPort    string 	`json:"http_port"`
	IPAddr        string 	`json:"bind_addr"`
}

// Response stores Answers for invite sending
type Response struct {
	ID string
	Answers map[string] string
}

// Answer is for parsing json response from typeform
type Answer struct {
	Responses []Response
}

func inviteAll(config Config){
	log.Println("inviteAll")
	typeFormURL := fmt.Sprintf("https://api.typeform.com/v0/form/%s?key=%s&completed=true&since=%d", config.TUID, config.TKey, time.Now().Unix() - config.Interval)
	res, err := http.Get(typeFormURL)

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
		log.Printf("Sending invite to %s %s with email %s\n", user.Answers[config.NameField], user.Answers[config.LastNameField], user.Answers[config.EmailField])
		slackURL := fmt.Sprintf("https://%s.slack.com/api/users.admin.invite", config.SlackChannel)

		invite := url.Values{}
		invite.Add("email", user.Answers[config.EmailField])
		invite.Add("first_name", user.Answers[config.NameField])
		invite.Add("last_name", user.Answers[config.LastNameField])
		invite.Add("token", config.SlackToken)
		invite.Add("set_active", "true")

		req, err := http.NewRequest("POST", slackURL, strings.NewReader(invite.Encode()))
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		mainPage(w, r, config.TUID)
	})
	listenAddr := fmt.Sprintf("%s:%s", config.IPAddr, config.ListenPort)
	fmt.Printf("Listening on: %s\n", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}

func mainPage(w http.ResponseWriter, r *http.Request, UID string) {
	data := make(map[string] string)
	data["UID"] = UID
	t, _ := template.ParseFiles("public/index.html")
	t.Execute(w, data)
}