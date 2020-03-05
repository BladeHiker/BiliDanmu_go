package models

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type DanMuMsg struct {
	UID        uint32 `json:"uid"`
	Uname      string `json:"uname"`
	Ulevel     uint32 `json:"ulevel"`
	Text       string `json:"text"`
	MedalLevel uint32 `json:"medal_level"`
	MedalName  string `json:"medal_name"`
}

func NewDanmu() *DanMuMsg {
	return &DanMuMsg{
		UID:        0,
		Uname:      "",
		Ulevel:     0,
		Text:       "",
		MedalLevel: 0,
		MedalName:  "无勋章",
	}
}

type Gift struct {
	UUname   string `json:"u_uname"`
	Action   string `json:"action"`
	Price    uint32 `json:"price"`
	GiftName string `json:"gift_name"`
}

func NewGift() *Gift {
	return &Gift{
		UUname:   "",
		Action:   "",
		Price:    0,
		GiftName: "",
	}
}

type WelCome struct {
}

type Notice struct {
}

type CMD string

var (
	CMDDanmuMsg                  CMD = "DANMU_MSG"                     // 普通弹幕信息
	CMDSendGift                  CMD = "SEND_GIFT"                     // 普通的礼物，不包含礼物连击
	CMDWELCOME                   CMD = "WELCOME"                       // 欢迎VIP
	CMDWelcomeGuard              CMD = "WELCOME_GUARD"                 // 欢迎房管
	CMDEntry                     CMD = "ENTRY_EFFECT"                  // 欢迎舰长等头衔
	CMDRoomRealTimeMessageUpdate CMD = "ROOM_REAL_TIME_MESSAGE_UPDATE" // 房间关注数变动
)

func (c *Client) SendPackage(packetlen uint32, magic uint16, ver uint16, typeID uint32, param uint32, data []byte) (err error) {
	packetHead := new(bytes.Buffer)

	if packetlen == 0 {
		packetlen = uint32(len(data) + 16)
	}
	var pdata = []interface{}{
		packetlen,
		magic,
		ver,
		typeID,
		param,
	}

	// 将包的头部信息以大端序方式写入字节数组
	for _, v := range pdata {
		if err = binary.Write(packetHead, binary.BigEndian, v); err != nil {
			fmt.Println("binary.Write err: ", err)
			return
		}
	}

	// 将包内数据部分追加到数据包内
	sendData := append(packetHead.Bytes(), data...)

	// fmt.Println("本次发包消息为：", sendData)

	if err = c.conn.WriteMessage(websocket.BinaryMessage, sendData); err != nil {
		fmt.Println("c.conn.Write err: ", err)
		return
	}

	return
}

func (c *Client) ReceiveMsg() {
	pool := NewPool()
	go pool.Handle()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("ReadMsg err :", err)
			continue
		}

		switch msg[11] {
		case 8:
			fmt.Println("握手包收发完毕，连接成功")
			c.Connected = true
		case 3:
			onlineNow := ByteArrToDecimal(msg[16:])
			if uint32(onlineNow) != c.Room.Online {
				c.Room.Online = uint32(onlineNow)
				fmt.Println("当前房间人气变动：", uint32(onlineNow))
			}
		case 5:
			if inflated, err := ZlibInflate(msg[16:]); err != nil {
				// 代表是未压缩数据
				pool.MsgUncompressed <- string(msg[16:])
			} else {
				for len(inflated) > 0 {
					l := ByteArrToDecimal(inflated[:4])
					c := json.Get(inflated[16:l], "cmd").ToString()
					switch CMD(c) {
					case CMDDanmuMsg:
						pool.UserMsg <- string(inflated[16:l])
					case CMDSendGift:
						pool.UserGift <- string(inflated[16:l])
					case CMDWELCOME:
						pool.UserGift <- string(inflated[16:l])
					case CMDWelcomeGuard:
						pool.UserGuard <- string(inflated[16:l])
					case CMDEntry:
						pool.UserEntry <- string(inflated[16:l])
					}
					inflated = inflated[l:]
				}
			}
		}
	}
}

func (c *Client) HeartBeat() {
	for {
		if c.Connected {
			obj := []byte("5b6f626a656374204f626a6563745d")
			if err := c.SendPackage(0, 16, 1, 2, 1, obj); err != nil {
				log.Println("heart beat err: ", err)
				continue
			}
			time.Sleep(30 * time.Second)
		}
	}
}
