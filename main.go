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
	Interval int64 `json:"check_interval"`
}

type Response struct {
	Id string
	Answers map[string] string
}
type Answer struct {
	Responses []Response
}

func inviteAll(config Config){
	typeFormUrl := fmt.Sprintf("https://api.typeform.com/v0/form/%s?key=%s&completed=true&since=%d", config.TUID, config.TKey, time.Now().Unix() - config.Interval)
	fmt.Println(typeFormUrl)
	res, err := http.Get(typeFormUrl)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(string(body))
	var answer Answer
	json.Unmarshal(body, &answer)
	fmt.Println(answer.Responses[0].Answers)
	fmt.Println("FirstName is " + answer.Responses[0].Answers[config.NameField])
	time.Sleep(1 * time.Second)
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

			inviteAll(config)

		}
	}()
	startHttpServer()
}

func startHttpServer(){
	http.Handle("/", http.FileServer(http.Dir("public")))

	http.ListenAndServe(":3000", nil)

}