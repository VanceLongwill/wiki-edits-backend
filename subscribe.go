package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type WikiChange struct {
	Action     string `json:"action"`
	ChangeSize int    `json:"change_size"`
}

type WikiSubscriber struct {
	log  *log.Logger
	url  string
	done <-chan struct{}
}

func (w *WikiSubscriber) Subscribe() (<-chan *WikiChange, <-chan error) {
	results := make(chan *WikiChange)
	errors := make(chan error)

	go func() {
		log.Printf("connecting to %s", w.url)
		c, _, err := websocket.DefaultDialer.Dial(w.url, nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		defer c.Close()

		defer close(results)
		defer close(errors)

		done := make(chan struct{})

		go func() {
			defer close(done)
			for {
				change := new(WikiChange)
				if err := c.ReadJSON(change); err != nil {
					log.Println("read:", err)
					return
				}
				select {
				case results <- change:
				case <-w.done:
					return
				}
			}
		}()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
				if err != nil {
					return
				}
			case <-w.done:
				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}
	}()
	return results, errors
}
