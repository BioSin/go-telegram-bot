package main

import (
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"time"
	"strconv"
)

const (
	ApiToken = "366368668:AAGXroejg7Nh1lGdwvTBmAc6BtdIgIFl_3E"
	ApiUrl   = "https://api.telegram.org/bot%s/%s"
)

var (
	httpClient http.Client
)

func main() {
	fmt.Println("Hello World!")

	httpClient = http.Client{}

	getMe()

	offset := 0

	ch := getUpdatesChan(offset)

	// Читаем из канала №1
	for update := range ch {
		process(update)
	}

	// Читаем из канала №2
	//msg := <-ch

	// Читаем из канала №3
	//for {
	//	msg := <-ch
	//}

	//v := Foo{
	//	Value: "Get method call",
	//}
	//v.Get();
}

func getMe() (error) {
	const methodName = "getMe"

	r, err := sendRequest(methodName, nil)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", r)

	return nil
}

func getUpdates(offset int) ([]Update, error) {
	const methodName = "getUpdates"
	updates := []Update{}

	data := url.Values{}

	if offset > 0 {
		data.Add("offset", strconv.Itoa(offset))
	}

	r, err := sendRequest(methodName, data)
	if err != nil {
		return updates, err
	}

	err = json.Unmarshal(r.Result, &updates)
	if err != nil {
		return updates, err
	}

	return updates, nil
}

func getUpdatesChan(offset int) chan Update {
	ch := make(chan Update)

	go func() {
		for {
			updates, err := getUpdates(offset)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				time.Sleep(time.Second * 5)
				continue
			}

			for _, update := range updates {
				offset = update.Id + 1
				ch <- update
			}

			time.Sleep(time.Second * 1)
		}
	}()

	return ch
}

func sendMessage(msg MessageData) (bool, error) {
	const methodName = "sendMessage"

	data := url.Values{}
	data.Add("chat_id", strconv.Itoa(msg.ChatId))
	data.Add("text", msg.Text)
	data.Add("parse_mode", msg.ParseMode)
	data.Add("reply_to_message_id", strconv.Itoa(msg.ReplyTo))

	_, err := sendRequest(methodName, data)
	if err != nil {
		return false, err
	}

	return true, nil
}
func process(update Update) {
	msg := MessageData{
		ChatId: update.Message.Chat.Id,
		Text: fmt.Sprintf("Hi, %s! I'm your dad, Luk!", update.Message.From.FirstName),
		ParseMode: "HTML",
		ReplyTo: update.Message.Id,
	}

	_, err := sendMessage(msg)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}
}

func sendRequest(methodName string, data url.Values) (ApiResponse, error) {
	apiUrl := fmt.Sprintf(ApiUrl, ApiToken, methodName)

	r := ApiResponse{}

	resp, err := httpClient.PostForm(apiUrl, data)
	if err != nil {
		return r, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(bytes, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}

type ApiResponse struct {
	Ok          bool
	Description string
	Result      json.RawMessage
}

type Update struct {
	Id      int `json:"update_id"`
	Message Message
}

type Message struct {
	Id   int `json:"message_id"`
	From User
	Chat Chat
	Text string
}

type User struct {
	Id        int
	FirstName string `json:"first_name"`
}

type Chat struct {
	Id int
}

type MessageData struct {
	ChatId    int `json:"char_id"`
	Text      string
	ParseMode string `json:"parse_mode"`
	ReplyTo   int `json:"reply_to_message_id"`
}

// Just for try
//type IFoo interface {
//	Get()
//}
//
//type Foo struct {
//	Value string
//}
//
//func (this *Foo) Get() {
//	fmt.Println(this.Value);
//}
//
//func foo(v IFoo) {
//	v.Get()
//}
