package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/lesta-battleship/server-core/pkg/matchmaking/packets"
)

var (
	pathFlag  = flag.String("path", "/join/custom", "")
	queryFlag = flag.String("query", "", "")
)

type GuildChatMessage struct {
	Msg string `json:"content"`
}

// Type: History, Error

type GuildChatHistory struct {
	Type string                    `json:"type"`
	Data []GuildChatHistoryMessage `json:"data"`
	// Data []map[string]any `json:"data"`
}

type GuildChatHistoryMessage struct {
	Id      string `json:"_id"`
	GuildId int    `json:"guild_id"`
	UserId  int    `json:"user_id"`
	Content string `json:"content"`
	// Timestamp any    `json:"timestamp"`
}

// func (m *GuildChatHistoryMessage) UnmarshalJSON(bytes []byte) error {
// 	data := struct {
// 		Id        string `json:"_id"`
// 		GuildId   int    `json:"guild_id"`
// 		UserId    int    `json:"user_id"`
// 		Content   string `json:"content"`
// 		// Timestamp any    `json:"timestamp"`
// 	}{}
//
// 	if err := json.Unmarshal(bytes, &data); err != nil {
// 		return err
// 	}
//
// 	return nil
// }

type WebsocketClient struct {
	id string

	dialer *websocket.Dialer
	conn   *websocket.Conn
}

func NewWebsocketClient(id string) *WebsocketClient {
	return &WebsocketClient{
		id: id,

		dialer: websocket.DefaultDialer,
		conn:   nil,
	}
}

func (c *WebsocketClient) Connect(url string) error {
	// header := http.Header{}
	// header.Add("X-XSRF-TOKEN", c.id)

	// conn, _, err := c.dialer.Dial(url, header)
	conn, _, err := c.dialer.Dial(url, nil)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *WebsocketClient) GetPacket() any {
	var packet GuildChatHistory
	if err := c.conn.ReadJSON(&packet); err != nil {
		log.Println("read err: ", err)

		return packets.Packet{}
	}

	return packet
}

func (c *WebsocketClient) SendPacket(packet any) {
	err := c.conn.WriteJSON(packet)
	if err != nil {
		log.Println("write err: ", err)

		return
	}
}

func (c *WebsocketClient) Stop() {
	c.conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	if err := c.conn.Close(); err != nil {
	}
}

func main() {
	flag.Parse()

	u := "ws://37.9.53.187:8000/ws/guild/1/1"
	// u := url.URL{Scheme: "ws", Host: "/ws/guild/1/1", Path: *pathFlag, RawQuery: *queryFlag}
	log.Printf("connecting to %s", u)

	id := rand.Text()
	client := NewWebsocketClient(id)
	if err := client.Connect(u); err != nil {
		return
	}

	done, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		for {
			packet := client.GetPacket()

			fmt.Printf("%+v", packet)
		}
	}()

	go func() {
		for {
			select {
			default:
				var text string
				fmt.Scanln(&text)

				chatMsg := GuildChatMessage{Msg: text}

				client.SendPacket(chatMsg)
			case <-done.Done():
				client.Stop()
			}
		}
	}()

	// go func() {
	// 	for {
	// 		select {
	// 		default:
	// 			var text string
	// 			fmt.Scanln(&text)
	//
	// 			switch text {
	// 			case "quit":
	// 				client.SendPacket(packets.NewDisconnect(client.id))
	//
	// 				cancel()
	//
	// 				return
	// 			case "create":
	// 				client.SendPacket(packets.NewCreateRoom(client.id))
	// 			case "join":
	// 				var roomId string
	// 				fmt.Scanln(&roomId)
	//
	// 				client.SendPacket(packets.NewJoinRoom(client.id, roomId))
	// 			default:
	// 				client.SendPacket(packets.NewPlayerMessage(client.id, text))
	// 			}
	// 		case <-done.Done():
	// 			client.Stop()
	// 		}
	// 	}
	// }()

	<-done.Done()
}

func SendPacket(conn *websocket.Conn, packet packets.Packet) {
	err := conn.WriteJSON(packet)
	if err != nil {
		log.Println("write err: ", err)
		return
	}
}
