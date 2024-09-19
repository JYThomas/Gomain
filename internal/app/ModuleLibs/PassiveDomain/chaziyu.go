package PassiveDomain

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/JYThomas/Gomain/internal/pkg"
	"github.com/PuerkitoBio/goquery"
)

// 定义模块结构体
type MODULE_CHAZIYU struct {
	ModuleName string
}

// 获取域名数据
func (m_chaziyu MODULE_CHAZIYU) GetDomainNames(domain string, retrycounts int) (DomainNames []string, err error) {
	// 加载配置文件数据源url
	Config, err := pkg.LoadConfig()
	if err != nil {
		return []string{}, errors.New("Module CHAZIYU: Fail to load config file")
	}
	// 获取数据源爬虫目标链接
	BASICURL := Config.Section("PassiveDomain").Key("URL_CHAZIYU").String()
	// 目标url
	TargetURL := BASICURL + domain

	// 请求数据
	html, err := GetResponse_CHAZIYU(TargetURL, retrycounts)
	if err != nil {
		// 如果在三次请求都没获取到数据的情况下 要么网络问题 要么没有数据 直接丢弃
		return []string{}, errors.New("Module CHAZIYU: Fail to Get Response")
	}

	// 处理响应内容，提取域名目标
	DomainNames, err = ResolveHTML_CHAZIYU(html)
	if err != nil {
		return []string{}, errors.New("Module CHAZIYU: Fail to extract subdomains")
	}
	// 这里就先不对切片结果数据去重了 因为不同的数据源之间有可能是有重复的
	// 所以等到所有数据源结果跑完 然后再进行统一去重
	// subdomains = pkg.RemoveDuplicates(subdomains)
	return DomainNames, nil
}

// 发起目标请求 获取响应内容
func GetResponse_CHAZIYU(TargetURL string, retrycounts int) (html io.Reader, err error) {
	if retrycounts <= 0 {
		return nil, errors.New("Module CHAZIYU: Max retry attempts reached")
	}

	// 创建请求客户端
	client := pkg.MakeHttpClient()

	// 设置请求配置
	request, err := pkg.MakeReq(TargetURL)
	if err != nil {
		return nil, errors.New("Module CHAZIYU: Make Requests Fail")
	}

	// 发送 HTTP 请求
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("Module CHAZIYU: Send Requests Error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return GetResponse_CHAZIYU(TargetURL, retrycounts-1)
	}

	// 读取 HTML 内容到字符串
	htmlContent, err := ioutil.ReadAll(resp.Body)
	// 返回一个新的 io.Reader
	return strings.NewReader(string(htmlContent)), nil
}

// 解析响应html页面
func ResolveHTML_CHAZIYU(html io.Reader) ([]string, error) {
	// 定义结果存储切片
	domains_slice := make([]string, 0)
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return []string{}, errors.New("Module CHAZIYU: Fail to read HTML")
	}
	// 查找所有的表格
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		// 查找表格中的每一行
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			// 查找行中的第二个 td
			row.Find("td").Eq(1).Each(func(k int, cell *goquery.Selection) {
				text := cell.Text()
				domains_slice = append(domains_slice, strings.TrimSpace(string(text)))
			})
		})
	})
	return domains_slice, nil
}
