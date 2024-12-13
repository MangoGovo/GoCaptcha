package request

import (
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"time"
)

type Client struct {
	*resty.Client
	UserAgent string
}

func New(userAgent string) Client {
	client := resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetRetryCount(5).
		SetTimeout(3 * time.Second)
	s := Client{
		Client:    client,
		UserAgent: userAgent,
	}
	// 利用中间件实现请求日志
	//s.OnAfterResponse(midware.LogMiddleware)
	return s
}

func (s Client) Request() *resty.Request {
	return s.R().
		EnableTrace()
}
