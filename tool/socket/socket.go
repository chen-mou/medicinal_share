package socket

import (
	"context"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"io/ioutil"
	"medicinal_share/main/entity"
	"medicinal_share/main/model/user"
	"medicinal_share/tool"
	"medicinal_share/tool/encrypt/jwtutil"
	"medicinal_share/tool/pool"
	"net"
	"strconv"
	"time"
)

var p = pool.NewPool(500)

type ConnManager struct {
	cmap          map[string]net.Conn
	headerHandler map[string]func(*Conn, string) error
	handler       map[string]func(*Conn, string)
	uri           string
}

func NewConnManager(uri string) *ConnManager {
	return &ConnManager{
		uri:           uri,
		headerHandler: map[string]func(*Conn, string) error{},
		handler:       map[string]func(*Conn, string){},
	}
}

type Conn struct {
	conn net.Conn
	auth string
	info *entity.User
	Id   string
}

func NewConn(conn net.Conn) *Conn {
	return &Conn{conn: conn}
}

func (c Conn) Send(message string) {}

func (c *Conn) SetAuth(auth string) {
	c.auth = auth
}

func (c *Conn) Auth() error {
	if c.auth == "" {
		return errors.New("未登录")
	}
	data, err := jwtutil.Parse(c.auth)
	if err != nil {
		return err
	}
	val, ok := data["id"]
	if !ok {
		return errors.New("token有误")
	}
	id, err := strconv.ParseInt(val, 10, 64)
	usr := user.GetById(id)
	for _, v := range usr.Role {
		if v.Name == "doctor" {
			usr.DockerInfo = user.GetDoctorInfoById(id)
		}
	}
	c.info = usr
	return nil
}

func (cm ConnManager) Run() {
	listen, err := net.Listen("tcp", cm.uri)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				panic(err)
			}
			co := NewConn(conn)
			up := ws.Upgrader{
				OnHeader: func(key, value []byte) error {
					f, ok := cm.headerHandler[string(key)]
					if !ok {
						return nil
					}
					return f(co, tool.BytesToString(value))
				},
				OnBeforeUpgrade: func() (header ws.HandshakeHeader, err error) {
					return nil, co.Auth()
				},
			}
			_, err = up.Upgrade(conn)
			if err != nil {
				SendError(conn, err)
				return
			}
			ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
			err = p.Submit(ctx, func() {
				reader := wsutil.NewReader(conn, ws.StateServerSide)
				for {
					_, err := reader.NextFrame()
					if err != nil {
						panic(err)
					}
					data, err := ioutil.ReadAll(reader)
					fmt.Println(string(data))
					wsutil.WriteServerText(conn, []byte("receive:"))
				}
			})
			if err != nil {
				SendError(conn, err)
				return
			}
		}
	}()
}

func (cm *ConnManager) Message(path string, f func(conn *Conn, payload string)) *ConnManager {
	cm.handler[path] = f
	return cm
}

func (cm *ConnManager) HeaderHandler(key string, f func(*Conn, string) error) *ConnManager {
	cm.headerHandler[key] = f
	return cm
}

func SendError(conn net.Conn, err error) {
}
