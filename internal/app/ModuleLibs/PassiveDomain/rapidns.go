package PassiveDomain

import (
	"errors"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/JYThomas/Gomain/internal/pkg"
	"github.com/PuerkitoBio/goquery"
)

type MODULE_RAPIDDNS struct {
	ModuleName string
}

// 获取域名数据
func (m_rapiddns MODULE_RAPIDDNS) GetDomainNames(domain string, retrycounts int) (DomainNames []string, err error) {
	// 加载配置文件数据源url
	Config, err := pkg.LoadConfig()
	if err != nil {
		return []string{}, errors.New("Module RAPIDDNS: Fail to load config file")
	}
	// 获取数据源爬虫目标链接
	BASICURL := Config.Section("PassiveDomain").Key("URL_RAPIDDNS").String()

	// 设置循环分页获取数据源数据
	for i := 1; i <= 10; i++ {
		// 构造请求链接
		// https://www.rapiddns.io/s/bilibili.com?page=1
		TargetURL := BASICURL + domain + "?page=" + strconv.Itoa(i)

		// 请求数据
		html, err := GetResponse_RAPPIDNS(TargetURL, retrycounts)
		if err != nil {
			// 如果在三次请求都没获取到数据的情况下 要么网络问题 要么没有数据 直接丢弃
			return []string{}, errors.New("Module RAPIDDNS: Fail to Get Response")
		}

		// 解析响应 提取域名资产
		subdomains, err := ResolveHTML_RAPIDDNS(html)
		if err != nil {
			return []string{}, errors.New("Module RAPIDDNS: Fail to Extract subdomains")
		}
		// 数据合并 将切片元素追加到返回结果中
		DomainNames = append(DomainNames, subdomains...)
	}

	return DomainNames, nil
}

// 发起目标请求 获取响应内容
func GetResponse_RAPPIDNS(TargetURL string, retrycounts int) (html io.Reader, err error) {
	if retrycounts <= 0 {
		return nil, errors.New("Module RAPIDDNS: Max retry attempts reached")
	}

	// 创建请求客户端
	client := pkg.MakeHttpClient()

	// 设置请求配置
	request, err := pkg.MakeReq(TargetURL)
	if err != nil {
		return nil, errors.New("Module RAPIDDNS: Make Requests Fail")
	}

	// 发送 HTTP 请求
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("Module RAPIDDNS: Send Requests Error")
	}
	defer resp.Body.Close()

	// 响应结果判断 重试三次
	if resp.StatusCode != 200 {
		return GetResponse_RAPPIDNS(TargetURL, retrycounts-1)
	}

	// 读取 HTML 内容到字符串
	htmlContent, err := ioutil.ReadAll(resp.Body)
	// 返回一个新的 io.Reader
	return strings.NewReader(string(htmlContent)), nil
}

// 响应结果解析
func ResolveHTML_RAPIDDNS(html io.Reader) (DomainNames []string, err error) {
	// 定义结果存储切片
	domains_slice := make([]string, 0)
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return []string{}, errors.New("Module RAPIDDNS: Fail to read HTML")
	}

	// 查找html中的域名
	doc.Find("Table tbody tr").Each(func(i int, row *goquery.Selection) {
		// 查找每行中的第二个 <td>，即域名部分
		subdomain := row.Find("td").Eq(0).Text()
		domains_slice = append(domains_slice, strings.TrimSpace(subdomain))
	})

	return domains_slice, nil
}
