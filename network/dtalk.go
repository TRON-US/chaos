package devtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MsgInfo struct {
	MsgType string  `json:"msgtype"`
	Text    MsgText `json:"text"`
}

type MsgText struct {
	Content string `json:"content"`
}

type PayLoad struct {
	Msg       string   `json:"msg"`
	Priority  Priority `json:"priority"`
	ServiceId string   `json:"service_id"`
}

type Priority int

const (
	PRIORITY_0 Priority = iota
	PRIORITY_1
	PRIORITY_2
	PRIORITY_3
	PRIORITY_4
	PRIORITY_5
	PRIORITY_6
	PRIORITY_7
	PRIORITY_8
	PRIORITY_9
)

var (
	url string
)

func init() {
	token := "9214b993e5692d6ba59a03d38b246febf2612c6563d1ec93ace1b0da35cfaae2"
	url = fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", token)
}

/*
	Trigger warning in Dtalk.
	{
		msg: $msg,
		priority: $lvl,
		server_id: $serviceId
	}
*/
func Warning(msg string, lvl Priority, serviceId string) error {
	pl := &PayLoad{Msg: msg, Priority: lvl, ServiceId: serviceId}
	bs, e := json.Marshal(pl)
	if e != nil {
		return e
	}
	return doSend(bs)
}

// Send message(plain text) to Dtalk Robot.
func SendText(msg string, serviceId string) error {
	m := fmt.Sprintf("%s. serverId: %s", msg, serviceId)
	return doSend([]byte(m))
}

func doSend(bs []byte) error {
	mi := &MsgInfo{MsgType: "text", Text: MsgText{Content: string(bs)}}
	bs, e := json.Marshal(mi)
	if e != nil {
		return e
	}
	body := bytes.NewBuffer(bs)
	res, err := http.Post(url, "application/json;charset=utf-8", body)
	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	return err
}
