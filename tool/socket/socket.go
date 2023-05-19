package socket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"io/ioutil"
	user2 "medicinal_share/main/entity"
	"medicinal_share/main/model/user"
	"medicinal_share/tool"
	redis1 "medicinal_share/tool/db/redis"
	"medicinal_share/tool/encrypt/jwtutil"
	"medicinal_share/tool/encrypt/md5"
	"medicinal_share/tool/pool"
	"net"
	"net/url"
	"strconv"
	"time"
)

var p = pool.NewPool(500)

const (
	idPrefix = "Socket:ID:"
	channel  = "Socket_Message"
)

var notFound string

func init() {
	byt, _ := json.Marshal(map[string]any{
		"code": 404,
		"msg":  "没有找到目标方法",
	})
	notFound = string(byt)
}

type ConnManager struct {
	cmap          map[string]*Conn
	headerHandler map[string]func(*Conn, string) error
	handler       map[string]func(*Conn, string)
	uri           string
}

func NewConnManager(uri string) *ConnManager {
	return &ConnManager{
		uri: uri,
		headerHandler: map[string]func(*Conn, string) error{
			"Method": func(conn *Conn, s string) error {
				conn.method = s
				return nil
			},
		},
		handler: map[string]func(*Conn, string){},
		cmap:    map[string]*Conn{},
	}
}

type Conn struct {
	conn   net.Conn
	auth   string
	method string
	info   *user2.User
	//redisConn *redis.ClusterClient
	Id string
}

type Message struct {
	Method string `json:"method"`
	Data   string `json:"data"`
}

type Result struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

func NewConn(conn net.Conn) *Conn {
	return &Conn{conn: conn}
}

func (c Conn) send(message string) {
	byt, _ := json.Marshal(Result{
		Code: 0,
		Data: message,
	})
	wsutil.WriteServerText(c.conn, byt)
}

func SendTo(message string, sendTo int64) {
	id := strconv.FormatInt(sendTo, 10)
	cmd := redis1.DB.Get(context.TODO(), idPrefix+id)
	res, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return
		}
		panic(err)
	}
	m := RedisMessage{
		Msg:    message,
		SendTo: res,
	}
	i, err := redis1.DB.Publish(context.TODO(), channel, m).Result()
	for i == 0 || err != nil {
		i, err = redis1.DB.Publish(context.TODO(), channel, m).Result()
	}
}

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
			user.UpdateDoctorStatus(usr.Id, user.Online)
		}
	}
	c.info = usr
	return nil
}

func (c *Conn) GetCurrentUser() *user2.User {
	return c.info
}

type RedisMessage struct {
	SendTo string
	Msg    string
}

func (cm ConnManager) redisListener(client *redis.ClusterClient) {
	sub := client.Subscribe(context.TODO(), "Socket_Message")
	for {
		msg, _ := sub.ReceiveMessage(context.TODO())
		m := msg.Payload

		v := RedisMessage{}
		json.Unmarshal(tool.StringToBytes(m), v)

		conn, ok := cm.cmap[v.SendTo]
		if !ok {
			continue
		}
		go conn.send(v.Msg)
	}
}

func (cm ConnManager) Run() {
	listen, err := net.Listen("tcp", cm.uri)
	if err != nil {
		panic(err)
	}
	go cm.run(listen)
	go func() {
		cm.redisListener(redis1.DB)
	}()
}

// run 监听uri
func (cm ConnManager) run(listen net.Listener) {
	defer func() {
		err := recover()
		if err != nil {
			panic(err)
		}
		fmt.Println("gobwas 挂掉了")
	}()
	for {
		conn, err := listen.Accept()
		if err != nil {
			SendError(conn, err)
			continue
		}
		co := NewConn(conn)
		up := ws.Upgrader{
			OnRequest: func(uri []byte) error {
				url, err := url.Parse(string(uri))
				if err != nil {
					return err
				}
				token := url.Query().Get("token")
				co.SetAuth(token)
				return nil
			},
			OnHeader: func(key, value []byte) error {
				return cm.onHeader(key, value, co)
			},
			OnBeforeUpgrade: func() (header ws.HandshakeHeader, err error) {
				return cm.onBeforeUpgrade(co)
			},
		}
		_, err = up.Upgrade(conn)
		if err != nil {
			SendError(conn, err)
			continue
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
				//fmt.Println(strconv.FormatInt(co.info.Id, 10) + ":" + string(data))
				//wsutil.WriteServerText(conn, tool.StringToBytes("receive"))
				fmt.Printf("收到消息:%s\n", data)
				res := Message{}
				json.Unmarshal(data, &res)
				f, ok := cm.handler[res.Method]
				if !ok {
					co.send(notFound)
				}
				go func() {
					defer func() {
						err := recover()
						if err != nil {
							switch err.(type) {
							case error:
								SendError(co.conn, err.(error))
							case string:
								co.send(err.(string))
							}
						}
					}()
					f(co, res.Data)
				}()
			}
		})
		if err != nil {
			SendError(conn, err)
			return
		}
	}
}

// onHeader 解析Header的时候执行
func (cm ConnManager) onHeader(key, value []byte, co *Conn) error {
	f, ok := cm.headerHandler[string(key)]
	if !ok {
		return nil
	}
	return f(co, tool.BytesToString(value))
}

// onBeforeUpgrade 协议升级前执行
func (cm ConnManager) onBeforeUpgrade(co *Conn) (ws.HandshakeHeader, error) {
	err := co.Auth()
	if err != nil {
		return nil, err
	}
	id := strconv.FormatInt(co.info.Id, 10)
	co.Id = md5.Hash(getSocketId() + id)
	cm.cmap[co.Id] = co
	cmd := redis1.DB.Do(context.TODO(), "set", idPrefix+id, co.Id)
	_, err = cmd.Result()
	return nil, err
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
	byt, _ := json.Marshal(Result{
		Code: 1,
		Data: err.Error(),
	})
	wsutil.WriteServerText(conn, byt)
}
