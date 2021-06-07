package models

import (
	"fmt"
	"strings"
)

// Pool`s fields map CMD value
type Pool struct {
	UserMsg         chan []byte
	UserGift        chan []byte
	UserEnter       chan []byte
	UserGuard       chan []byte
	MsgUncompressed chan []byte
	UserEntry       chan []byte
}

func NewPool() *Pool {
	return &Pool{
		UserMsg:         make(chan []byte, 10),
		UserGift:        make(chan []byte, 10),
		UserEnter:       make(chan []byte, 10),
		MsgUncompressed: make(chan []byte, 10),
		UserEntry:       make(chan []byte, 10),
		UserGuard:       make(chan []byte, 10),
	}
}

func (pool *Pool) Handle() {
	for {
		select {
		case uc := <-pool.MsgUncompressed:
			// 目前只处理未压缩数据的关注数变化信息
			if cmd := json.Get(uc, "cmd").ToString(); CMD(cmd) == CMDRoomRealTimeMessageUpdate {
				fans := json.Get(uc, "data", "fans").ToInt()
				fmt.Println("当前房间关注数变动：", fans)
			}
		case src := <-pool.UserMsg:
			m := NewDanmu()
			m.GetDanmuMsg(src)
			fmt.Printf("%d-%s | %d-%s: %s\n", m.MedalLevel, m.MedalName, m.Ulevel, m.Uname, m.Text)
		case src := <-pool.UserGift:
			g := NewGift()
			g.GetGiftMsg(src)
			fmt.Printf("%s %s 价值 %d 的 %s\n", g.UUname, g.Action, g.Price, g.GiftName)
		case src := <-pool.UserEnter:
			name := json.Get(src, "data", "uname").ToString()
			fmt.Printf("欢迎VIP %s 进入直播间\n", name)
		case src := <-pool.UserGuard:
			name := json.Get(src, "data", "username").ToString()
			fmt.Printf("欢迎房管 %s 进入直播间\n", name)
		case src := <-pool.UserEntry:
			cw := json.Get(src, "data", "copy_writing").ToString()
			cw = strings.Replace(cw, "<%", "", 1)
			cw = strings.Replace(cw, "%>", "", 1)
			fmt.Printf("%s\n", cw)
		}
	}
}
