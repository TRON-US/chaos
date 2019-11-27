// Package dtalk implements the dtalk webhook calliing.
package dtalk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Priority int

const (
	// priority 0-9
	Priority0 Priority = iota
	Priority1
	Priority2
	Priority3
	Priority4
	Priority5
	Priority6
	Priority7
	Priority8
	Priority9
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

type ResponseBody struct {
	Code int    `json:"errcode"`
	Msg  string `json:"errmsg"`
}

// Trigger warning in Dtalk.
func Warning(url, msg string, priority Priority, serviceId string) error {
	pl := &PayLoad{Msg: msg, Priority: priority, ServiceId: serviceId}
	bs, e := json.Marshal(pl)
	if e != nil {
		return e
	}
	return doSend(url, bs)
}

// Send message(plain text) to Dtalk Robot.
func SendText(url, msg string, serviceId string) error {
	m := fmt.Sprintf("%s. serverId: %s", msg, serviceId)
	return doSend(url, []byte(m))
}

// Http post to dtalk.
func doSend(url string, bs []byte) error {
	// Init message information.
	mi := &MsgInfo{MsgType: "text", Text: MsgText{Content: string(bs)}}

	// Convert byte to json
	bs, err := json.Marshal(mi)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(bs))
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	resBodyByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	responseBody := &ResponseBody{}
	err = json.Unmarshal(resBodyByte, responseBody)
	if err != nil {
		return err
	}

	if responseBody.Code != 0 {
		return errors.New(responseBody.Msg)
	}
	return err
}
