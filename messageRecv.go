package workwx

type MessageRUrl struct {
	Signature string `url:"msg_signature"`
	Timestamp string `url:"timestamp"`
	Nonce     string `url:"nonce"`
}

type MessageRVerify struct {
	MessageRUrl
	EchoStr string `url:"echostr"`
}

type MessageRBody struct {
	ToUserName string `xml:"ToUserName"`
	AgentId    string `xml:"AgentId"`
	Encrypt    string `xml:"Encrypt"`
}

type MessageR struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	AgentId      string `xml:"AgentId"`
	MsgType      string `xml:"MsgType"`
	MsgId        string `xml:"MsgId"`
	CreateTime   string `xml:"CreateTime"`
	//text
	Content string `xml:"Content,omitempty"`

	//image
	PicUrl  string `xml:"PicUrl,omitempty"`
	MediaId string `xml:"MediaId,omitempty"`

	//voice, also use MediaId
	Format string `xml:"Format,omitempty"`

	//video, also use MediaId
	ThumbMediaId string `xml:"ThumbMediaId,omitempty"`

	//location
	X       float64 `xml:"Location_X,omitempty"`
	Y       float64 `xml:"Location_Y,omitempty"`
	Scale   int     `xml:"Scale,omitempty"`
	Label   string  `xml:"Label,omitempty"`
	AppType string  `xml:"AppType,omitempty"`

	//link, also use PicUrl
	Url         string `xml:"Url,omitempty"`
	Title       string `xml:"Title,omitempty"`
	Description string `xml:"Description,omitempty"`

	//event
	Event    string `xml:"Event,omitempty"`
	EventKey string `xml:"EventKey,omitempty"`
}

type MessageRResponse struct {
	ToUserName   string `xml:"ToUserName,cdata,omitempty"`
	FromUserName string `xml:"FromUserName,cdata,omitempty"`
	CreateTime   string `xml:"CreateTime"`
	MsgType      string `xml:"MsgType,cdata"`
}

type MessageRResponseText struct {
	MessageRResponse
	Content string `xml:"Content,cdata"`
}

type MessageRResponseImage struct {
	MessageRResponse
	Image MessageRResponseMediaData `xml:"Image"`
}

type MessageRResponseVoice struct {
	MessageRResponse
	Voice MessageRResponseMediaData `xml:"Voice"`
}

type MessageRResponseMediaData struct {
	MediaId string `xml:"MediaId,cdata"`
}

type MessageRResponseVideo struct {
	MessageRResponse
	Video MessageRRepsonseVideoData `xml:"Video"`
}

type MessageRRepsonseVideoData struct {
	MessageRResponseMediaData
	Title       string `xml:"Title,cdata,omitempty"`
	Description string `xml:"Description,cdata,omitempty"`
}

type MessageRResponseNews struct {
	MessageRResponse
	ArticleCount int                        `xml:"ArticleCount"`
	Articles     []MessageRResponseNewsData `xml:"Articles,omitempty"`
}

type MessageRResponseNewsData struct {
	Item MessageRResponseNewsDataItem `xml:"item"`
}

type MessageRResponseNewsDataItem struct {
	Title       string `xml:"Title,cdata,omitempty"`
	Description string `xml:"Description,cdata,omitempty"`
	PicUrl      string `xml:"PicUrl,cdata,omitempty"`
	Url         string `xml:"Url,cdata,omitempty"`
}

const (
	MessageTypeText  = "text"
	MessageTypeImage = "image"
	MessageTypeVoice = "voice"
	MessageTypeVideo = "video"
	MessageTypeNews  = "news"
)
