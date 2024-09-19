package PassiveDomain

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/JYThomas/Gomain/internal/pkg"
)

// 定义模块结构体
type MODULE_CRTSH struct {
	ModuleName string
}

// 定义用于解析 JSON 的结构体
type Certificate struct {
	IssuerCAID     int    `json:"issuer_ca_id"`
	IssuerName     string `json:"issuer_name"`
	CommonName     string `json:"common_name"`
	NameValue      string `json:"name_value"`
	ID             int    `json:"id"`
	EntryTimestamp string `json:"entry_timestamp"`
	NotBefore      string `json:"not_before"`
	NotAfter       string `json:"not_after"`
	SerialNumber   string `json:"serial_number"`
	ResultCount    int    `json:"result_count"`
}

// 发起请求获取域名
func (m_crtsh MODULE_CRTSH) GetDomainNames(domain string, retrycounts int) ([]string, error) {
	if retrycounts <= 0 {
		return []string{}, errors.New("Module CRTSH: Max retry attempts reached")
	}

	url := "https://crt.sh/?q=" + domain

	// 创建请求客户端
	client := pkg.MakeHttpClient()

	// 设置请求配置
	request, err := pkg.MakeReq(url)
	if err != nil {
		// panic 用于处理程序中的严重错误或不可恢复的异常
		// panic(err)
		return []string{}, errors.New("Module CRTSH: Make Requests Fail")
	}

	// 发送 HTTP 请求
	resp, err := client.Do(request)
	if err != nil {
		return []string{}, errors.New("Module CRTSH: Send Requests Error")
	}
	defer resp.Body.Close()

	// 响应错误 重试
	if resp.StatusCode != 200 {
		return m_crtsh.GetDomainNames(domain, retrycounts-1)
	}

	// 处理响应内容，提取域名目标
	subdomains, err := ResolveHTML_CRTSH(resp.Body)
	if err != nil {
		return []string{}, errors.New("Module CRTSH: Extract domains Error")
	}
	result := pkg.RemoveDuplicates(subdomains)
	return result, err
}

// 解析响应的html页面
func ResolveHTML_CRTSH(html io.Reader) (result []string, err error) {
	// 读取响应体 响应内容为json字符串
	content, err := io.ReadAll(html)
	if err != nil {
		return []string{}, errors.New("Module CRTSH: failed to read response body")
	}

	var certificates []Certificate
	// 将Json字符串内容解析为Certificate结构体切片 Json字符串内容与结构体结构相同
	err = json.Unmarshal([]byte(content), &certificates)

	if err != nil {
		return []string{}, errors.New("Module CRTSH: failed to unmarshal JSON")
	}

	for _, cert := range certificates {
		result = append(result, cert.CommonName)
	}

	return result, nil
}
