package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

// 消息表
type Message struct {
	gorm.Model
	FormId   int64  //发送者
	TargetId int64  //接受者
	Type     int    //消息通知类型 群聊  私聊  广播
	Media    int    //消息类型  1文字 2图片 3音频
	Content  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int //其他数字统计
}

// 发送信息使用的
type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSet  set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

// 聊天 得到发送者id 接收者id  消息类型  消息内容 消息发送的类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	//校验token等合法性
	//token := query.Get("token")
	query := request.URL.Query()
	id := query.Get("userId")
	userId, _ := strconv.ParseInt(id, 10, 64)
	/*
		targetId := query.Get("targetId")
		context := query.Get("context")
		msgType := query.Get("type")
	*/
	isvalida := true
	conn, err := (&websocket.Upgrader{
		//token校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//获取conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSet:  set.New(set.ThreadSafe),
	}
	//用户关系

	//userid 和node绑定
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	//完成发送的逻辑
	go sendProc(node)
	//完成接受的逻辑
	go recvProc(node)
	sendMsg(userId, []byte("欢迎进入聊天室"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		udpsendChan <- data
	}
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

// 完成udp协程发送数据
func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 3000,
	})
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan:
			_, err := conn.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成udp协程接受数据
func udpRecvProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	for {
		var buf [512]byte
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度处理逻辑
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		sendMsg(msg.TargetId, data)
		/*case 2: //群发
			sendGroupMsg()
		case 3: //广播
			sendAllMsg()
			//case 4:
		*/
	}
}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
