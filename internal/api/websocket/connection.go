package websocket

import (
	"fmt"
	"lesta-battleship/cli/internal/api/websocket/packets"
	"net/http"

	"github.com/gorilla/websocket"
)

// Интерфейс, реализующий принцип записи и чтения WebsocketClient.
type Strategy interface {
	ReadPump(readChan chan<- packets.Packet, conn *websocket.Conn) error
	WritePump(writeChan <-chan packets.Packet, conn *websocket.Conn) error
}

const maxChanBuffer = 100

// Абстракция над websocket соединением к серверу.
//
// Считывает пакеты от сервера в readChan.
// Сохраняет пакеты для записи на сервер в writeChan.
// Сохраняет ошибки в errorChan.
//
// Работа зависит от переданного Strategy.
type WebsocketClient struct {
	readChan  chan packets.Packet
	writeChan chan packets.Packet
	errorChan chan error

	strategy Strategy

	dialer *websocket.Dialer
	conn   *websocket.Conn
}

func NewWebsocketClient(strategy Strategy) *WebsocketClient {
	return &WebsocketClient{
		readChan:  make(chan packets.Packet, maxChanBuffer),
		writeChan: make(chan packets.Packet, maxChanBuffer),
		errorChan: make(chan error),

		strategy: strategy,

		dialer: websocket.DefaultDialer,
	}
}

// Метод для устнановки websocket соединения с сервером.
//
// Параметр url является обязательным, header - опциональным.
//
// Возвращает ошибку при отсутствие возможности подключится к серверу.
func (c *WebsocketClient) Connect(url string, header http.Header) error {
	conn, _, err := c.dialer.Dial(url, header)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

// Метод для чтения пакета из канала readChan.
func (c *WebsocketClient) GetPacket() packets.Packet {
	return <-c.readChan
}

// Метод для записи пакета в канал writeChan.
func (c *WebsocketClient) SendPacket(packet packets.Packet) {
	c.writeChan <- packet
}

// Метод для возвращения канала с ошибками. Только для чтения.
func (c *WebsocketClient) ErrorChan() <-chan error {
	return c.errorChan
}

// Запускает чтение пакетов от сервера.
// Сохраняет пакеты в канал readChan для дальшейшего чтения клиентом.
//
// При ошибке заканчивает чтение и разрывает websocket соединение с сервером.
//
// Рекомендуется запускать в горутине.
func (c *WebsocketClient) ReadPump() {
	defer func() {
		if _, ok := <-c.readChan; !ok {
			close(c.readChan)
		}

		c.Stop()
	}()

	err := c.strategy.ReadPump(c.readChan, c.conn)
	c.errorChan <- fmt.Errorf("WebsocketClient: %w", err)
}

// Запускает запись пакетов от клиента.
// Считывает пакеты из канала writeChan для дальнейшей передачи на сервер.
//
// При ошибке заканчивает запись и разрывает websocket соединение.
//
// Рекомендуется запускать в горутине.
func (c *WebsocketClient) WritePump() {
	defer func() {
		if _, ok := <-c.writeChan; !ok {
			close(c.writeChan)
		}

		c.Stop()
	}()

	err := c.strategy.WritePump(c.writeChan, c.conn)
	c.errorChan <- fmt.Errorf("WebsocketClient: %w", err)
}

// Метод для разрыва websocket соединения с сервером.
//
// При разрыве соединения ReadPump и WritePump также заканчивают свою работу.
func (c *WebsocketClient) Stop() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
		}
		c.conn = nil
	}
}
