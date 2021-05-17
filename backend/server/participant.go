package server

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Participant struct {
	Uuid string
	Conn *websocket.Conn
	Coop *Pool
	//other user information?
}

type Message struct {
	Message map[string]interface{}
	Client *Participant
}

func (c *Participant) Read(pool *Pool) {
	defer func() {
		pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		msg := Message{Client: c} //Attach the client to the Message struct? Sure...
		err := c.Conn.ReadJSON(&msg.Message) //Read JSON from Connect

		if msg.Message != nil {
			msg.Message["size"] = len(pool.Clients) //For each msg send the size of the room back to clients
		}

		fmt.Println("LOG MESSAGE:")
		fmt.Println(msg.Message["join"])

		//Set UUID to the participant on Join
		if msg.Message["join"] == true {
			fmt.Println("LOG UUID:")
			fmt.Println( msg.Message["uuid"].(string))
			c.Uuid = msg.Message["uuid"].(string)
		}

		if err != nil {
			log.Println(err)
			return
		}

		pool.Broadcast <- msg
		fmt.Printf("Message Received: %+v\n", msg)
	}

}