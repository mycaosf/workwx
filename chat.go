package workwx

type Chat struct {
	Token
	chatId string
}

// Create create chat. If you want to set chatid, you should call SetChatId before Create.
func (p *Chat) Create(name, owner string, userList []string) (ret ChatCreateResponse) {
	data := &ChatInfo{
		Name: name, Owner: owner, UserList: userList, ChatId: p.chatId,
	}

	p.postJson(chatClass, chatApiCreate, data, &ret)

	return
}

func (p *Chat) Rename(name string) Error {
	data := &chatModifyParam{Name: name}

	return p.modify(data)
}

func (p *Chat) ResetOwner(owner string) Error {
	data := &chatModifyParam{Owner: owner}

	return p.modify(data)
}

func (p *Chat) AddUsers(users []string) Error {
	data := &chatModifyParam{AddUserList: users}

	return p.modify(data)
}

func (p *Chat) DelUsers(users []string) Error {
	data := &chatModifyParam{DelUserList: users}

	return p.modify(data)
}

func (p *Chat) Get() (ret ChatGetResponse) {
	p.getJson(chatClass, chatApiGet, &ret, chatIdStr, p.chatId)

	return
}

func (p *Chat) modify(data *chatModifyParam) (ret Error) {
	data.ChatId = p.chatId
	p.postJson(chatClass, chatApiModify, data, &ret)

	return
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

type ChatCreateResponse struct {
	Error
	ChatId string `json:"chatid"`
}

type chatModifyParam struct {
	ChatId      string   `json:"chatid"`
	Name        string   `json:"name,omitempty"`
	Owner       string   `json:"owner,omitempty"`
	AddUserList []string `json:"add_user_list,omitempty"`
	DelUserList []string `json:"del_user_list,omitempty"`
}

type ChatGetResponse struct {
	Error
	Info ChatInfo `json:"chat_info"`
}

const (
	chatClass     = "appchat"
	chatApiCreate = "create"
	chatApiModify = "update"
	chatApiGet    = "get"
	chatIdStr     = "&chatid="
)
