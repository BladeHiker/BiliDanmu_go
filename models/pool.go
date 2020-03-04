package models

import "fmt"

type Pool struct {
	UserMsg         chan string
	UserGift        chan string
	UserEnter       chan string
	MsgUncompressed chan string
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
		}
	}
}
