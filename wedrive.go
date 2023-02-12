package workwx

func wedrivePost(p *Token, api string, req any, res any) {
	p.postJson(wedriveClass, api, req, res)
}

const (
	wedriveClass = "wedrive"
)
