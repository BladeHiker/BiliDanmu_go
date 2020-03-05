package models

import (
	"fmt"
)

// Pool`s fields map CMD value
type Pool struct {
	UserMsg         chan string
	UserGift        chan string
	UserEnter       chan string
	UserGuard       chan string
	MsgUncompressed chan string
	UserEntry       chan string
}

func NewPool() *Pool {
	return &Pool{
		UserMsg:         make(chan string, 10),
		UserGift:        make(chan string, 10),
		UserEnter:       make(chan string, 10),
		MsgUncompressed: make(chan string, 10),
	}
}

func (pool *Pool) Handle() {
	for {
		select {
		case uc := <-pool.MsgUncompressed:
			// 目前只处理未压缩数据的关注数变化信息
			if cmd := json.Get([]byte(uc), "cmd").ToString(); CMD(cmd) == CMDRoomRealTimeMessageUpdate {
				fans := json.Get([]byte(uc), "data", "fans").ToInt()
				fmt.Println("当前房间关注数变动：", fans)
			}
		case src := <-pool.UserMsg:
			m := NewDanmu()
			m.GetDanmuMsg([]byte(src))
			fmt.Printf("%d-%s | %d-%s: %s\n", m.MedalLevel, m.MedalName, m.Ulevel, m.Uname, m.Text)
		case src := <-pool.UserGift:
			g := NewGift()
			g.GetGiftMsg([]byte(src))
			fmt.Printf("%s %s 价值 %d 的 %s\n", g.UUname, g.Action, g.Price, g.GiftName)
		case src := <-pool.UserEnter:
			name := json.Get([]byte(src), "data", "uname").ToString()
			fmt.Printf("欢迎VIP %s 进入直播间", name)
		case src := <-pool.UserGuard:
			name := json.Get([]byte(src), "data", "username").ToString()
			fmt.Printf("欢迎房管 %s 进入直播间", name)
		case src := <-pool.UserEntry:
			cw := json.Get([]byte(src), "data", "copy_writing").ToString()
			fmt.Printf("%s", cw)
		}
	}
}
