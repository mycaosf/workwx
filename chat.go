package workwx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Chat struct {
	token
	chatId string
}

// Create create chat. If you want to set chatid, you should call SetChatId before Create.
func (p *Chat) Create(name, owner string, userList []string) (chatId string, err error) {
	var r chatCreateResponse
	data := &ChatInfo{
		Name: name, Owner: owner, UserList: userList, ChatId: p.chatId,
	}

	if err = p.send(chatApiCreate, data, &r); err == nil {
		if err = r.parse(); err == nil {
			if r.ChatId != "" {
				chatId = r.ChatId
				p.chatId = r.ChatId
			} else {
				err = fmt.Errorf("Create error: code: %d, msg: %s", r.ErrCode, r.ErrMsg)
			}
		}
	}

	return
}

func (p *Chat) send(api string, data, r interface{}) error {
	if buf, err := json.Marshal(data); err != nil {
		return err
	} else {
		buffer := bytes.NewReader(buf)

		return p.postJson(chatClass, api, buffer, r)
	}
}

func (p *Chat) Rename(name string) error {
	data := &chatModifyParam{Name: name}

	return p.modify(data)
}

func (p *Chat) ResetOwner(owner string) error {
	data := &chatModifyParam{Owner: owner}

	return p.modify(data)
}

func (p *Chat) AddUsers(users []string) error {
	data := &chatModifyParam{AddUserList: users}

	return p.modify(data)
}

func (p *Chat) DelUsers(users []string) error {
	data := &chatModifyParam{DelUserList: users}

	return p.modify(data)
}

func (p *Chat) Get() (ret ChatInfo, err error) {
	var r chatGetResponse
	if err = p.getJson(chatClass, chatApiGet, &r, chatIdStr, p.chatId); err == nil {
		if err = r.parse(); err == nil {
			ret = r.Info
		}
	}

	return
}

func (p *Chat) modify(data *chatModifyParam) error {
	var r Error
	data.ChatId = p.chatId
	err := p.send(chatApiModify, data, &r)
	if err == nil {
		err = r.parse()
	}

	return err
}

func (p *Chat) SetChatId(id string) {
	p.chatId = id
}

type ChatInfo struct {
	Name     string   `json:"name,omitempty"`
	Owner    string   `json:"owner,omitempty"`
	UserList []string `json:"userlist"`
	ChatId   string   `json:"chatid,omitempty"`
}

type chatCreateResponse struct {
	baseResponse
	ChatId string `json:"chatid"`
}

type chatModifyParam struct {
	ChatId      string   `json:"chatid"`
	Name        string   `json:"name,omitempty"`
	Owner       string   `json:"owner,omitempty"`
	AddUserList []string `json:"add_user_list,omitempty"`
	DelUserList []string `json:"del_user_list,omitempty"`
}

type chatGetResponse struct {
	baseResponse
	Info ChatInfo `json:"chat_info"`
}

const (
	chatClass     = "appchat"
	chatApiCreate = "create"
	chatApiModify = "update"
	chatApiGet    = "get"
	chatIdStr     = "&chatid="
)
