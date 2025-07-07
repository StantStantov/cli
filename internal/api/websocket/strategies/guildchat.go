package strategies

import (
	"lesta-battleship/cli/internal/api/websocket/packets"
	"lesta-battleship/cli/internal/api/websocket/packets/guild"

	"github.com/gorilla/websocket"
)

// Стратегия для WebsocketClient.
//
// Ожидает от сервера пакеты типа guild.Packet.
//
// При получении пакета guild.Disconnect принудительно заканчивает работу.
type GuildChatStrategy struct{}

func (c GuildChatStrategy) ReadPump(readChan chan<- packets.Packet, conn *websocket.Conn) error {
	isFirstMessage := true
	for {
		var packet guild.Packet = new(guild.ChatMessage)
		if isFirstMessage {
			packet = new(guild.ChatHistory)
			isFirstMessage = false
		}

		if err := conn.ReadJSON(&packet); err != nil {
			return err
		}

		readChan <- packets.WrapGuild(packet)
	}
}

func (c GuildChatStrategy) WritePump(writeChan <-chan packets.Packet, conn *websocket.Conn) error {
	for packet := range writeChan {
		var unwrap guild.Packet
		packets.UnwrapAsGuild(packet, &unwrap)

		switch unwrap.(type) {
		case *guild.Disconnect:
			return nil
		default:
			err := conn.WriteJSON(packet.Content())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
