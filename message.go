package workwx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type ErrorMessage struct {
	ErrCode      int
	ErrMsg       string
	InvalidUser  []string
	InvalidParty []string
	InvalidTag   []string
}

func (p *ErrorMessage) Error() string {
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
	baseResponse
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

type Message struct {
	token
	toUser    []string
	toParty   []string
	toTag     []string
	chatId    string
	agentId   int
	className string
}

type sendMessageCommonReal struct {
	ToUser  string `json:"touser,omitempty"`
	ToParty string `json:"toparty,omitempty"`
	ToTag   string `json:"totag,omitempty"`
	ChatId  string `json:"chatid,omitempty"`
	MsgType string `json:"msgtype"`
	AgentId int    `json:"agentid,omitempty"`
}

type MessageText struct {
	Content string //max 2048 bytes
	Safe    bool
	IdTrans bool
}

type sendMessageTextReal struct {
	sendMessageCommonReal
	Text    ContentText `json:"text"`
	Safe    int         `json:"safe, omitempty"`
	IdTrans int         `json:"enable_id_trans,omitempty"`
}

type ContentText struct {
	Content string `json:"content"`
}

type MessageMarkdown struct {
	Content string //max 2048 bytes
}

type sendMessageMarkdownReal struct {
	sendMessageCommonReal
	Markdown ContentText `json:"markdown"`
}

// receivers: user, party, tag, chatid
func (p *Message) SetReceiver(receivers ...[]string) {
	for len(receivers) < 3 {
		receivers = append(receivers, nil)
	}

	if len(receivers) > 3 && len(receivers[3]) == 1 {
		p.SetReceiverChatId(receivers[3][0])
	} else {
		p.toUser = receivers[0]
		p.toParty = receivers[1]
		p.toTag = receivers[2]
		p.className = messageClass
		p.chatId = ""
	}
}

func (p *Message) SetReceiverChatId(chatId string) {
	p.toUser = nil
	p.toParty = nil
	p.toTag = nil
	p.chatId = chatId
	p.className = chatClass
}

func (p *Message) SetAgentId(agentId int) {
	p.agentId = agentId
}

// send text message
func (p *Message) Text(text *MessageText) error {
	var data sendMessageTextReal
	p.toRealCommon(&data.sendMessageCommonReal, "text")

	data.Text.Content = text.Content
	if text.Safe {
		data.Safe = 1
	}
	if text.IdTrans {
		data.IdTrans = 1
	}

	return p.send(&data)
}

func (p *Message) toRealCommon(to *sendMessageCommonReal, msgType string) {
	if p.toUser != nil {
		to.ToUser = strings.Join(p.toUser, toJoinStr)
	}
	if p.toParty != nil {
		to.ToParty = strings.Join(p.toParty, toJoinStr)
	}
	if p.toTag != nil {
		to.ToTag = strings.Join(p.toTag, toJoinStr)
	}
	if p.chatId != "" {
		to.ChatId = p.chatId
	}

	to.MsgType = msgType
	to.AgentId = p.agentId
}

func (p *Message) send(data interface{}) error {
	if buf, err := json.Marshal(data); err != nil {
		return err
	} else {
		buffer := bytes.NewReader(buf)
		var r sendMessageResponseReal
		if err = p.postJson(p.className, messageApiSend, buffer, &r); err != nil {
			return err
		}

		return sendMessageRet(&r)
	}
}

func sendMessageRet(r *sendMessageResponseReal) (err error) {
	if r.ErrCode != 0 {
		ret := &ErrorMessage{ErrCode: r.ErrCode, ErrMsg: r.ErrMsg}
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

	return
}

// send markdown message
func (p *Message) Markdown(markdown *MessageMarkdown) error {
	var data sendMessageMarkdownReal
	p.toRealCommon(&data.sendMessageCommonReal, "markdown")
	data.Markdown.Content = markdown.Content

	return p.send(&data)
}

var (
	receiverTypes = []string{"user", "party", "tag"}
)

const (
	messageClass   = "message"
	messageApiSend = "send"
	toJoinStr      = "|"
)
