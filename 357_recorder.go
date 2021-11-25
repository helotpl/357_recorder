package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Message struct {
	Tag   string `json:"tag"`
	Links struct {
		ExternalURL string `json:"external_url"`
		URI         string `json:"uri"`
	} `json:"links"`
	Node string `json:"node"`
}

func main() {
	c, _, err := websocket.DefaultDialer.Dial("wss://socket.r357.eu/socket", nil)
	if err != nil {
		log.Fatal(err)
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
		fmt.Println(message.Tag)
	}
}
