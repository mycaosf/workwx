package workwx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mycaosf/utils/net/httpc"
	"io/ioutil"
	"net/http"
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
	toUser  []string
	toParty []string
	toTag   []string
	agentId int
}

type sendMessageCommonReal struct {
	ToUser  string `json:"touser,omitempty"`
	ToParty string `json:"toparty,omitempty"`
	ToTag   string `json:"totag,omitempty"`
	MsgType string `json:"msgtype"`
	AgentId int    `json:"agentid"`
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

func (p *Message) SetReceiver(user, party, tag []string) {
	p.toUser = user
	p.toParty = party
	p.toTag = tag
}

func (p *Message) SetAgentId(agentId int) {
	p.agentId = agentId
}

//send text message
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
	to.MsgType = msgType
	to.AgentId = p.agentId
}

func (p *Message) send(data interface{}) error {
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
		header := make(http.Header)
		header.Add(httpc.HTTPHeaderContentType, contentJson)

		if resp, err = httpPost(url, header, buffer); err != nil {
			buffer.Seek(0, 0)
			if token, err = p.Get(true); err != nil {
				return err
			}

			url = urlSendMessage + token
			if resp, err = httpPost(url, header, buffer); err != nil {
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
		}
	}

	return
}

//send markdown message
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
	urlSendMessage = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="
	toJoinStr      = "|"
	contentJson    = "application/json"
)
