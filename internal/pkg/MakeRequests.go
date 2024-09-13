package pkg

import (
	"net/http"
	"time"
)

// 设置请求客户端(注意声明返回值类型)
func MakeHttpClient() *http.Client {
	// 创建一个自定义的 HTTP 客户端，设置超时时间
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置超时时间为 10 秒
	}
	return client
}

// 设置请求配置
func MakeReq(url string) (*http.Request, error) {
	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")
	req.Header.Set("Accept", "application/json")

	return req, nil
}
