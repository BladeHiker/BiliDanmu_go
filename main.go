package main

import (
	"biliDanMu/models"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var roomid uint32
	fmt.Print("请输入房间号，长短 ID 均可：")
	_, err := fmt.Scanf("%d", &roomid)
	if err != nil {
		log.Println("房间号输入错误，请退出重新输入！")
		os.Exit(0)
	}

	// 兼容房间号短 ID
	if roomid >= 100 && roomid < 1000 {
		roomid,err = models.GetRealRoomID(int(roomid))
		if err != nil {
			log.Println("房间号输入错误，请退出重新输入！")
			os.Exit(0)
		}
	}

	c, err := models.NewClient(roomid)
	if err != nil {
		fmt.Println("models.NewClient err: ", err)
		return
	}
	if err = c.Start(); err != nil {
		fmt.Println("c.Start err :", err)
		return
	}

	time.Sleep(time.Minute * 3)
}
