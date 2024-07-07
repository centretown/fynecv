package comm

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"

	"nhooyr.io/websocket"
)

const (
	CONFIG = 20 + iota
	STATES
	SUBSCRIBE
)

const (
	dial = "ws://melon:8123/api/websocket"
)

type WebSockClient struct {
	conn      *websocket.Conn
	ctx       context.Context
	MessageID int
	Err       error
	Quit      chan int
	Buffer    chan []byte
}

func NewWebSockClient() (*WebSockClient, error) {
	ctx := context.Background()
	conn, resp, err := websocket.Dial(ctx, dial, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Dial", resp.Status)
	hs := &WebSockClient{
		ctx:       ctx,
		conn:      conn,
		Quit:      make(chan int),
		Buffer:    make(chan []byte),
		MessageID: STATES,
	}
	return hs, err
}

const BUFFER_SIZE = 1024 * 32

func (hs *WebSockClient) Read() ([]byte, error) {
	var readBuffer []byte = make([]byte, BUFFER_SIZE)

	typ, rdrConn, err := hs.conn.Reader(hs.ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Type", typ)
	rdr := bufio.NewReaderSize(rdrConn, BUFFER_SIZE)

	for {

		count, err := rdr.Read(readBuffer)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if count > 0 {
			log.Println("read count", count)
			return readBuffer[:count], nil
		}
	}

}

func (hs *WebSockClient) Write(cmd string) error {
	log.Println(cmd)
	err := hs.conn.Write(hs.ctx, websocket.MessageText, []byte(cmd))
	if err != nil {
		log.Println(err)
	}
	return err
}

func (hs *WebSockClient) WriteID(cmd string) {
	log.Println(cmd)
	idCmd := fmt.Sprintf(cmd, hs.MessageID)
	hs.MessageID++
	err := hs.conn.Write(hs.ctx, websocket.MessageText, []byte(idCmd))
	if err != nil {
		log.Println(err)
	}
}
