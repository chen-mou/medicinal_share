package socket

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"io"
	"testing"
	"time"
)

type Header map[string]string

func (h Header) WriteTo(w io.Writer) (int64, error) {
	count := int64(0)
	for k, v := range h {
		byt := []byte(k + ":" + v)
		n, err := w.Write(append(byt))
		if err != nil {
			return 0, err
		}
		count += int64(n)
	}
	return 0, nil
}

func TestSendAndGet(t *testing.T) {
	cm := NewConnManager("localhost:15889")
	cm.Run()
	d := ws.Dialer{
		Header: Header{
			"test": "test",
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
		time.Sleep(2 * time.Second)
		wsutil.WriteClientMessage(conn, ws.OpText, []byte("text2"))
		msg, _, _ = wsutil.ReadServerData(conn)
		fmt.Println(string(msg))
	}
}
