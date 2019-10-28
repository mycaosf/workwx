package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ErrorSendMessage struct {
	ErrCode      int
	ErrMsg       string
	InvalidUser  []string
	InvalidParty []string
	InvalidTag   []string
}

func (p *ErrorSendMessage) Error() string {
	ret := fmt.Sprintf("errcode: %d, errmsg: %s", p.ErrCode, p.ErrMsg)
	msgs := [][]string{p.InvalidUser, p.InvalidParty, p.InvalidTag}
	for i := 0; i < len(msgs); i++ {
		msg := msgs[i]
		if msg != nil {
			ret += "invlaid " + receiverTypes[i] + ": " + strings.Join(msg, toJoinStr)
		}
	}

	return ret
}

type sendMessageResponseReal struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

type SendMessage struct {
	token
	toUser  []string
	toParty []string
	toTag   []string
	agentId int
}

type sendMessageDataCommonReal struct {
	ToUser  string `json:"touser,omitempty"`
	ToParty string `json:"toparty,omitempty"`
	ToTag   string `json:"totag,omitempty"`
	MsgType string `json:"msgtype"`
	AgentId int    `json:"agentid"`
}

type SendMessageDataText struct {
	Content string //max 2048 bytes
	Safe    bool
	IdTrans bool
}

type sendMessageDataTextReal struct {
	sendMessageDataCommonReal
	Text    ContentText `json:"text"`
	Safe    int         `json:"safe, omitempty"`
	IdTrans int         `json:"enable_id_trans,omitempty"`
}

type ContentText struct {
	Content string `json:"content"`
}

type SendMessageDataMarkdown struct {
	Content string //max 2048 bytes
}

type sendMessageDataMarkdownReal struct {
	sendMessageDataCommonReal
	Markdown ContentText `json:"markdown"`
}

func (p *SendMessage) SetReceiver(user, party, tag []string) {
	p.toUser = user
	p.toParty = party
	p.toTag = tag
}

func (p *SendMessage) SetAgentId(agentId int) {
	p.agentId = agentId
}

//send text message
func (p *SendMessage) Text(text *SendMessageDataText) error {
	var data sendMessageDataTextReal
	p.toRealCommon(&data.sendMessageDataCommonReal, "text")

	data.Text.Content = text.Content
	if text.Safe {
		data.Safe = 1
	}
	if text.IdTrans {
		data.IdTrans = 1
	}

	return p.send(&data)
}

func (p *SendMessage) toRealCommon(to *sendMessageDataCommonReal, msgType string) {
	if p.toUser != nil {
		to.ToUser = strings.Join(p.toUser, toJoinStr)
	}
	if p.toParty != nil {
		to.ToParty = strings.Join(p.toParty, toJoinStr)
	}
	if p.toTag != nil {
		to.ToTag = strings.Join(p.toTag, toJoinStr)
	}
	to.MsgType = msgType
	to.AgentId = p.agentId
}

func (p *SendMessage) send(data interface{}) error {
	if buf, err := json.Marshal(data); err != nil {
		return err
	} else {
		var token string
		var resp *http.Response
		if token, err = p.Get(false); err != nil {
			return err
		}

		url := urlSendMessage + token
		buffer := bytes.NewReader(buf)

		if resp, err = http.Post(url, contentJson, buffer); err != nil {
			buffer.Seek(0, 0)
			if token, err = p.Get(true); err != nil {
				return err
			}

			url = urlSendMessage + token
			if resp, err = http.Post(url, contentJson, buffer); err != nil {
				return err
			}

		}

		return sendMessageRet(resp)
	}
}

func sendMessageRet(resp *http.Response) (err error) {
	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err == nil {
		var r sendMessageResponseReal
		if err = json.Unmarshal(body, &r); err == nil {
			if r.ErrCode != 0 {
				ret := &ErrorSendMessage{ErrCode: r.ErrCode, ErrMsg: r.ErrMsg}
				if r.InvalidUser != "" {
					ret.InvalidUser = strings.Split(r.InvalidUser, toJoinStr)
				}
				if r.InvalidParty != "" {
					ret.InvalidParty = strings.Split(r.InvalidParty, toJoinStr)
				}
				if r.InvalidTag != "" {
					ret.InvalidTag = strings.Split(r.InvalidTag, toJoinStr)
				}

				err = ret
			}
		}
	}

	return
}

//send markdown message
func (p *SendMessage) Markdown(markdown *SendMessageDataMarkdown) error {
	var data sendMessageDataMarkdownReal
	p.toRealCommon(&data.sendMessageDataCommonReal, "markdown")
	data.Markdown.Content = markdown.Content

	return p.send(&data)
}

var (
	receiverTypes = []string{"user", "party", "tag"}
)

const (
	urlSendMessage = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="
	toJoinStr      = "|"
	contentJson    = "application/json"
)
