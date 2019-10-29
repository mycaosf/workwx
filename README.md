# workwx
The workwx repository holds packages for using work weixin.

The workwx project is experimental. Use this at your own risk.

Send Messsage exmaples:
```go
	msg := &Message{}
	msg.Set(corpId, secret)
	msg.SetReceiver([]string{userId}, nil, nil)
	msg.SetAgentId(agentId)

	data := MessageText{ Content: "test" }
	if err := msg.Text(&data); err != nil {
	}

```
