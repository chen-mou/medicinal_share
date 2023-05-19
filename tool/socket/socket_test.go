package socket

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"io"
	"net"
	"testing"
	"time"
)

type Header map[string]string

func (h Header) WriteTo(w io.Writer) (int64, error) {
	count := int64(0)
	for k, v := range h {
		byt := []byte(k + ":" + v + "\n")
		n, err := w.Write(byt)
		if err != nil {
			return 0, err
		}
		count += int64(n)
	}
	return 0, nil
}

func TestSendAndGet(t *testing.T) {
	cm := NewConnManager("localhost:15889")
	cm.HeaderHandler("Token", func(conn *Conn, s string) error {
		conn.SetAuth(s)
		return nil
	})
	cm.Run()
	go func() {
		d := ws.Dialer{
			Header: Header{
				"Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEYXRhIjp7ImlkIjoiMTQ0In0sImV4cCI6MTY3MjY2OTM3OCwiaXNzIjoiR0FURVdBWV9TRVJWRVIifQ.BZIwaUeOXDbFY7fv6mxZRDQFmJEapESuP3v6BPvg6Z0",
			},
		}
		conn, _, _, err := d.Dial(context.TODO(), "ws://localhost:15889")
		if err != nil {
			panic(err)
		}
		for {
			wsutil.WriteClientMessage(conn, ws.OpText, []byte("text1"))
			msg, _, _ := wsutil.ReadServerData(conn)
			fmt.Println(string(msg))
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(2 * time.Second)

	go func() {
		d := ws.Dialer{
			Header: Header{
				"Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEYXRhIjp7ImlkIjoiMTYwIn0sImV4cCI6MTY3MjY3MDI2NiwiaXNzIjoiR0FURVdBWV9TRVJWRVIifQ.jySS8rj6w0yjw6m5C2c3HgJZ-etxEsuzKHpqY0uwwnc",
			},
		}
		//ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
		conn, _, _, err := d.Dial(context.TODO(), "ws://localhost:15889")
		if err != nil {
			panic(err)
		}
		for {
			wsutil.WriteClientMessage(conn, ws.OpText, []byte("text1"))
			msg, _, _ := wsutil.ReadServerData(conn)
			fmt.Println(string(msg))
			time.Sleep(1 * time.Second)
		}
	}()
	select {}
}

func TestTcp(t *testing.T) {
	_, err := net.Dial("tcp", "127.0.0.1:15777")
	if err != nil {
		panic(err)
	}

}
