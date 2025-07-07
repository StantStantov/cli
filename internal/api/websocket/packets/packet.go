package packets

import (
	"fmt"
	"lesta-battleship/cli/internal/api/websocket/packets/guild"
	"reflect"
)

type Packet interface {
	Content() any
	IsPacket()
}

type PacketWrapper struct {
	content any `json:"-"`
}

func (p PacketWrapper) Content() any {
	return p.content
}
func (PacketWrapper) IsPacket() {}

// Заворачивает guild.Packet в packets.Packet.
func WrapGuild(packet guild.Packet) Packet {
	return PacketWrapper{content: packet}
}

// Разворачивает packets.Packet в guild.Packet.
// Результат разворота сохраняется в value.
//
// Возвращает ошибку при:
// 1. Передачи в параметр value не указателя на значение.
// 2. Передачи в параметр packet пакета, содержимое которого не реализует интерфейс guild.Packet.
func UnwrapAsGuild(packet Packet, value any) error {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("UnwrapAsGuild: Parameter value isn't a pointer")
	}
	rv = rv.Elem()

	content, ok := packet.Content().(guild.Packet)
	if !ok {
		return fmt.Errorf("UnwrapAsGuild: Can't type assert packet contents as guild.Packet")
	}

	rv.Set(reflect.ValueOf(content))
	return nil
}
