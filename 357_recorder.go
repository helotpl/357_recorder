package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
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

type Message2 struct {
	S string `json:"s"`
	T string `json:"t"`
	P struct {
		N string `json:"n"`
		A string `json:"a"`
		T string `json:"t"`
		D int    `json:"d"`
		P int    `json:"p"`
	} `json:"p"`
}

func getCurrent() string {
	resp, err := http.Get("https://stream.radio357.pl/now/playing.json")
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var message Message2
	if err = json.Unmarshal(body, &message); err != nil {
		log.Println(err)
	}
	return message.S
}

func SocketReceiver(out chan Message) {
	c, _, err1 := websocket.DefaultDialer.Dial("wss://socket.r357.eu/socket", nil)
	if err1 != nil {
		log.Fatal(err1)
	}
	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			continue
		}
		var message Message
		if err = json.Unmarshal(data, &message); err != nil {
			fmt.Println(err)
		}
		out <- message
	}
}

func durationGraph(duration time.Duration) string {
	m := duration.Truncate(time.Minute)
	minutes := int(math.Round(m.Minutes()))
	tenSeconds := int(math.Round((duration - m).Round(time.Second*10).Seconds() / 10))
	if tenSeconds == 6 {
		minutes += 1
		tenSeconds = 0
	}
	return strings.Repeat("#", minutes) + strings.Repeat(":", tenSeconds/2) + strings.Repeat(".", tenSeconds%2)
}

func main() {
	noTrack := "Radio 357 - Najlepszy radiowy adres na świecie!"

	duration := flag.Int("dur", 3600, "duration of logging")
	interval := flag.Int("i", 2, "interval between downloads")
	names := flag.Bool("names", false, "record only track names")
	flag.Parse()
	onlyTrackNames := *names
	if flag.NArg() != 1 {
		log.Fatal("Please supply non option parameter = output filename")
	}
	fileName := flag.Args()[0]
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()
	quit := time.After(time.Duration(*duration) * time.Second)

	prevt := time.Now()
	current := getCurrent()
	previous := current
	if current != noTrack {
		_, err = fmt.Fprintln(file, time.Now().Format(time.Stamp), current)
		if err != nil {
			fmt.Println(err)
		}
	} else if !onlyTrackNames {
		_, err = fmt.Fprint(file, time.Now().Format(time.Stamp))
		if err != nil {
			fmt.Println(err)
		}
	}
	for {
		select {
		case t := <-ticker.C:
			current = getCurrent()
			if current != previous {
				if previous == noTrack && !onlyTrackNames {
					d := t.Sub(prevt)
					_, err = fmt.Fprintln(file, "", d.Round(time.Second), durationGraph(d))
					if err != nil {
						fmt.Println(err)
					}
				}
				previous = current
				prevt = t
				if current != noTrack {
					_, err = fmt.Fprintln(file, t.Format(time.Stamp), current)
					if err != nil {
						fmt.Println(err)
					}
				} else if !onlyTrackNames {
					_, err = fmt.Fprint(file, t.Format(time.Stamp))
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		case <-quit:
			if previous == noTrack && !onlyTrackNames {
				d := time.Now().Sub(prevt)
				_, err = fmt.Fprintln(file, "", d.Round(time.Second), durationGraph(d))
				if err != nil {
					log.Fatal(err)
				}
			}
			os.Exit(0)
		}
	}
}
