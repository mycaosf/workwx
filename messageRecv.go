package workwx

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"strconv"
	"time"
)

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

type CData struct {
	Data string `xml:",cdata"`
}

type MessageRResponse struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CData    `xml:"ToUserName,omitempty"`
	FromUserName CData    `xml:"FromUserName,omitempty"`
	CreateTime   string   `xml:"CreateTime"`
	MsgType      CData    `xml:"MsgType"`
}

type MessageRResponseText struct {
	MessageRResponse
	Content CData `xml:"Content"`
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
	MediaId CData `xml:"MediaId"`
}

type MessageRResponseVideo struct {
	MessageRResponse
	Video MessageRRepsonseVideoData `xml:"Video"`
}

type MessageRRepsonseVideoData struct {
	MessageRResponseMediaData
	Title       CData `xml:"Title,omitempty"`
	Description CData `xml:"Description,omitempty"`
}

type MessageRResponseNews struct {
	MessageRResponse
	ArticleCount int                        `xml:"ArticleCount"`
	Articles     []MessageRResponseNewsData `xml:"Articles,omitempty"`
}

type MessageRResponseNewsData struct {
	XMLName     xml.Name `xml:"item"`
	Title       CData    `xml:"Title,omitempty"`
	Description CData    `xml:"Description,omitempty"`
	PicUrl      CData    `xml:"PicUrl,omitempty"`
	Url         CData    `xml:"Url,omitempty"`
}

type MessageResponseBody struct {
	XMLName   xml.Name `xml:"xml"`
	Encrypt   CData    `xml:"Encrypt"`
	Signature CData    `xml:"MsgSignature"`
	Nonce     CData    `xml:"Nonce"`
	Timestamp string   `xml:"TimeStamp"`
}

func BuildMessageResponseBody(c *Crypto, v interface{}) (*MessageResponseBody, error) {
	if data, err := xml.Marshal(v); err != nil {
		return nil, err
	} else {
		var ret MessageResponseBody

		nonce := make([]byte, 16)
		rand.Read(nonce)
		nonceStr := base64.StdEncoding.EncodeToString(nonce)
		ret.Nonce.Data = nonceStr

		message := c.Encrypt(data)
		ret.Encrypt.Data = message
		ret.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
		ret.Signature.Data = c.Signature(ret.Timestamp, nonceStr, message)

		return &ret, nil
	}
}

const (
	MessageTypeText  = "text"
	MessageTypeImage = "image"
	MessageTypeVoice = "voice"
	MessageTypeVideo = "video"
	MessageTypeNews  = "news"
)
