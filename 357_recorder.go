package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"time"
)

type Message struct {
	Tag   string `json:"tag"`
	Links struct {
		ExternalURL string `json:"external_url"`
		URI         string `json:"uri"`
	} `json:"links"`
	Node string `json:"node"`
}

func SocketReceiver(out chan Message) {
	c, _, err1 := websocket.DefaultDialer.Dial("wss://socket.r357.eu/socket", nil)
	if err1 != nil {
		log.Fatal(err1)
	}
	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		var message Message
		if err = json.Unmarshal(data, &message); err != nil {
			log.Fatal(err)
		}
		out <- message
	}
}

func main() {
	msgChan := make(chan Message, 1)
	go SocketReceiver(msgChan)
	quit := time.After(time.Minute)
	for {
		select {
		case msg := <-msgChan:
			fmt.Println(msg.Tag)

		case <-quit:
			os.Exit(0)
		}
	}
}
