package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	// telegram token
	TOKEN = ""
	// URL telegram
	URL = "https://api.telegram.org/bot"
)

func update(w http.ResponseWriter, r *http.Request) {

	message := &tgbotapi.Update{}

	var chatID int64
	chatID = 0
	msgText := ""

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		log.Println(err)
	}

	// if private or group
	if message.Message.Chat.ID != 0 {
		log.Println(message.Message.Chat.ID, message.Message.Text)
		chatID = message.Message.Chat.ID
		msgText = message.Message.Text
	} else {
		// if channel
		log.Println(message.ChannelPost.Chat.ID, message.ChannelPost.Text)
		chatID = message.ChannelPost.Chat.ID
		msgText = message.ChannelPost.Text
	}

	respMsg := fmt.Sprintf("%s%s/sendMessage?chat_id=%d&text=Received: %s", URL, TOKEN, chatID, msgText)

	// send echo resp
	_, err = http.Get(respMsg)
	if err != nil {
		log.Println(err)
	}
}

func main() {

	telegramToken, exists := os.LookupEnv("TELEGRAM_TOKEN")
	if exists {
		TOKEN = telegramToken
	} else {
		panic("Telegram token not defined")
	}

	telegramApiUrl, exists := os.LookupEnv("TELEGRAM_API_URL")
	if exists {
		URL = telegramApiUrl
	}

	listenPort := flag.String("port", "3000", "Listenning port")
	flag.Parse()

	http.HandleFunc("/", update)

	log.Println("Listenning on port", *listenPort, ".")
	if err := http.ListenAndServe(":"+*listenPort, logRequest(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
