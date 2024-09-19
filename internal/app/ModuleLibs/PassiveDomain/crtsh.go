package PassiveDomain

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"strings"

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
func (m_crtsh MODULE_CRTSH) GetDomainNames(domain string, retrycounts int) (DomainNames []string, err error) {
	// 加载配置文件数据源url
	Config, err := pkg.LoadConfig()
	if err != nil {
		return []string{}, errors.New("Module CRTSH: Fail to load config file")
	}
	// 获取数据源爬虫目标链接
	BASICURL := Config.Section("PassiveDomain").Key("URL_CRTSH").String()
	TargetURL := BASICURL + domain

	// 请求数据
	html, err := GetResponse_CRTSH(TargetURL, retrycounts)
	if err != nil {
		// 如果在三次请求都没获取到数据的情况下 要么网络问题 要么没有数据 直接丢弃
		return []string{}, errors.New("Module CRTSH: Fail to Get Response")
	}

	// 处理响应内容，提取域名目标
	DomainNames, err = ResolveHTML_CRTSH(html)
	if err != nil {
		return []string{}, errors.New("Module CRTSH: Extract domains Error")
	}

	// 汇总时统一去重
	// result := pkg.RemoveDuplicates(subdomains)
	return DomainNames, nil
}

// 发起目标请求 获取响应内容
func GetResponse_CRTSH(TargetURL string, retrycounts int) (html io.Reader, err error) {
	if retrycounts <= 0 {
		return nil, errors.New("Module CRTSH: Max retry attempts reached")
	}

	// 创建请求客户端
	client := pkg.MakeHttpClient()

	// 设置请求配置
	request, err := pkg.MakeReq(TargetURL)
	if err != nil {
		// panic 用于处理程序中的严重错误或不可恢复的异常
		// panic(err)
		return nil, errors.New("Module CRTSH: Make Requests Fail")
	}

	// 发送 HTTP 请求
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("Module CRTSH: Send Requests Error")
	}

	// 响应错误 重试
	if resp.StatusCode != 200 {
		return GetResponse_CRTSH(TargetURL, retrycounts-1)
	}

	// 读取 HTML 内容到字符串
	htmlContent, err := ioutil.ReadAll(resp.Body)
	// 返回一个新的 io.Reader
	return strings.NewReader(string(htmlContent)), nil
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
