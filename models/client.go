package models

import (
	"fmt"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"net/url"
)

var
(
	RealID      = "http://api.live.bilibili.com/room/v1/Room/room_init" 				// params: id=xxx
	DanMuServer = "ks-live-dmcmt-bj6-pm-02.chat.bilibili.com:443"
	keyUrl      = "https://api.live.bilibili.com/room/v1/Danmu/getConf"                 // params: room_id=xxx&platform=pc&player=web
	roomInfoUrl = "https://api.live.bilibili.com/xlive/web-room/v1/index/getInfoByRoom" // params: room_id=xxx
	json        = jsoniter.ConfigCompatibleWithStandardLibrary
)

// Client instance
type Client struct {
	Room      RoomInfo     `json:"room"`
	Request   *RequestInfo `json:"request"`
	conn      *websocket.Conn
	Connected bool `json:"connected"`
}

// Basic information of the live room
type RoomInfo struct {
	RoomId     uint32 `json:"room_id"`
	UpUid      uint32 `json:"up_uid"`
	Title      string `json:"title"`
	Online     uint32 `json:"online"`
	Tags       string `json:"tags"`
	LiveStatus bool   `json:"live_status"`
	LockStatus bool   `json:"lock_status"`
}

// data on handshake packets
type RequestInfo struct {
	Uid       uint8  `json:"uid"`
	Roomid    uint32 `json:"roomid"`
	Protover  uint8  `json:"protover"`
	Platform  string `json:"platform"`
	Clientver string `json:"clientver"`
	Type      uint8  `json:"type"`
	Key       string `json:"key"`
}

// NewRequestInfo return initialized structure
func NewRequestInfo(roomid uint32) *RequestInfo {
	t := GetToken(roomid)
	return &RequestInfo{
		Uid:       0,
		Roomid:    roomid,
		Protover:  2,
		Platform:  "web",
		Clientver: "1.10.2",
		Type:      2,
		Key:       t,
	}
}

// new websocket("wss)
func NewClient(roomid uint32) (c *Client, err error) {
	return &Client{
		Room:      GetRoomInfo(roomid),
		Request:   NewRequestInfo(roomid),
		conn:      nil,
		Connected: false,
	}, nil
}

func (c *Client) Start() (err error) {
	u := url.URL{Scheme: "wss", Host: DanMuServer, Path: "/sub",}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		fmt.Println("net.Dial err: ", err)
		return err
	}
	c.conn = conn

	fmt.Println("当前直播间状态：", c.Room.LiveStatus)

	fmt.Println("连接弹幕服务器 ", DanMuServer, " 成功，正在发送握手包...")
	r, err := json.Marshal(c.Request)

	if err != nil {
		fmt.Println("marshal err ,", err)
		return
	}
	if err = c.SendPackage(0, 16, 1, 7, 1, r); err != nil {
		fmt.Println("SendPackage err,", err)
		return
	}
	go c.ReceiveMsg()
	go c.HeartBeat()
	return
}
