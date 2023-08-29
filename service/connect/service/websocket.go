package service

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 房间存储
var rooms = make(map[string]*room)

// websocket房间模型
// 关于房间的生命周期：对于个人来说，一方打开另一方的页面时，建立websocket连接。当双方页面都处于关闭状态的时候，websocket连接关闭。
// 客户端维持一个计数器，与服务器的计数器进行比对，确保消息同步。(PS:溢出
type room struct {
	//房间内websocket客户端
	clients map[string]*websocket.Conn
	//房间id
	id string //对id拼接，小id在前大id在后
	//房间消息计数器
	counter int //服务器接收到一条消息后，先于counter对比，如果小于counter则先返回一条要求同步的消息。在对比完成后，计数器加1。
	//并发锁
	mu sync.Mutex
	//房间当前连接数
	connectCount int
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, //允许跨域
}

func connect(c *gin.Context) {
	uid, isExist := c.Get("user_id")
	if !isExist {
		respError(c, 401, errors.New("token无效"))
		return
	}

	friendID := c.Param("id") //客户端直接拼起来就是了
	//获取好友信息
	friend, err := getFriendInfo(friendID)
	if err != nil {
		respError(c, 500, err)
		return
	}

	//与房间进行连接
	connectRoom(uid.(string), friend, c)
}

func connectRoom(uid string, friend Friend, c *gin.Context) {
	//检查房间是否已经存在，如果没有则创建房间
	r, ok := rooms[friend.Id]
	if !ok {
		r = createRoom(friend.Id)
		r.counter = friend.Counter
	}

	//建立websocket连接
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	//加入房间
	r.addClient(uid, conn)

	//开始接收消息
	go r.read(conn)
}

func createRoom(roomID string) *room {
	//从数据库获取counter

	var r = new(room)
	r.id = roomID
	r.clients = make(map[string]*websocket.Conn)
	r.counter = 0

	rooms[roomID] = r

	return r
}

// 添加客户端
func (r *room) addClient(id string, conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[id] = conn
	r.connectCount++
}
