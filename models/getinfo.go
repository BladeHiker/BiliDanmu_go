package models

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
)

func GetRealRoomID(short int) (realID uint32, err error) {
	url := fmt.Sprintf("%s?id=%d", RealID, short)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get token err: ", err)
		return 0, err
	}

	rawdata, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		fmt.Println("ioutil.ReadAll(resp.Body) err: ", err)
		return 0, err
	}
	realID = json.Get(rawdata, "data", "room_id").ToUint32()

	return realID, nil
}

// GetToken return the necessary token for connecting to the server
func GetToken(roomid uint32) (key string) {
	url := fmt.Sprintf("%s?room_id=%d&platform=pc&player=web", keyUrl, roomid)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get token err: ", err)
		return
	}

	rawdata, err := ioutil.ReadAll(resp.Body)

	_ = resp.Body.Close()
	if err != nil {
		fmt.Println("ioutil.ReadAll(resp.Body) err: ", err)
		return
	}
	key = json.Get(rawdata, "data").Get("token").ToString()
	return
}

func GetRoomInfo(roomid uint32) (roomInfo RoomInfo) {
	// roomInfo = &models.RoomInfo{}
	url := fmt.Sprintf("%s?room_id=%d", roomInfoUrl, roomid)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get roomInfo err: ", err)
		return
	}

	rawdata, err := ioutil.ReadAll(resp.Body)

	_ = resp.Body.Close()
	if err != nil {
		fmt.Println("ioutil.ReadAll(resp.Body) err: ", err)
		return
	}

	roomInfo.RoomId = roomid
	roomInfo.UpUid = json.Get(rawdata, "data").Get("room_info").Get("uid").ToUint32()
	roomInfo.Title = json.Get(rawdata, "data").Get("room_info").Get("title").ToString()
	roomInfo.Tags = json.Get(rawdata, "data").Get("room_info").Get("tags").ToString()
	roomInfo.LiveStatus = json.Get(rawdata, "data").Get("room_info").Get("live_status").ToBool()
	roomInfo.LockStatus = json.Get(rawdata, "data").Get("room_info").Get("lock_status").ToBool()

	return
}

func ZlibInflate(compress []byte) ([]byte, error) {
	var out bytes.Buffer
	c := bytes.NewReader(compress)
	r, err := zlib.NewReader(c)
	if err != zlib.ErrChecksum && err != zlib.ErrDictionary && err != zlib.ErrHeader && r != nil {
		_, _ = io.Copy(&out, r)
		if err := r.Close(); err != nil {
			fmt.Println("r.close err:", err)
			return nil, err
		}
		return out.Bytes(), nil
	}
	return nil, err
}

func (d *DanMuMsg) GetDanmuMsg(source []byte) {
	d.UID = json.Get(source, "info", 2, 0).ToUint32()
	d.Uname = json.Get(source, "info", 2, 1).ToString()
	d.Ulevel = json.Get(source, "info", 4, 0).ToUint32()
	d.Text = json.Get(source, "info", 1).ToString()
	d.MedalName = json.Get(source, "info", 3, 1).ToString()
	if d.MedalName == "" {
		d.MedalName = "无勋章"
	}
	d.MedalLevel = json.Get(source, "info", 3, 0).ToUint32()
	return
}

func (g *Gift) GetGiftMsg(source []byte) {
	g.UUname = json.Get(source, "data", "uname").ToString()
	g.Action = json.Get(source, "data", "action").ToString()
	nums := json.Get(source, "data", "num").ToUint32()
	g.Price = json.Get(source, "data", "price").ToUint32() * nums
	g.GiftName = json.Get(source, "data", "giftName").ToString()
}

// 返回字节数组表示数的十进制形式
func ByteArrToDecimal(src []byte) (sum int) {
	if src == nil {
		return 0
	}
	b := []byte(hex.EncodeToString(src))
	l := len(b)
	for i := l - 1; i >= 0; i-- {
		base := int(math.Pow(16, float64(l-i-1)))
		var mul int
		if int(b[i]) >= 97 {
			mul = int(b[i]) - 87
		} else {
			mul = int(b[i]) - 48
		}

		sum += base * mul
	}
	return
}
