package PassiveDomain

import (
	"errors"
	"io"
	"strings"

	"github.com/JYThomas/Gomain/internal/pkg"
	"github.com/PuerkitoBio/goquery"
)

// 定义模块结构体
type MODULE_CHAZIYU struct {
	ModuleName string
}

func (m_chaziyu MODULE_CHAZIYU) GetDomainNames(domain string, retrycounts int) ([]string, error) {
	if retrycounts <= 0 {
		return []string{}, errors.New("Module CHAZIYU: Max retry attempts reached")
	}

	// 目标url
	url := "https://chaziyu.com/" + domain

	// 创建请求客户端
	client := pkg.MakeHttpClient()

	// 设置请求配置
	request, err := pkg.MakeReq(url)
	if err != nil {
		return []string{}, errors.New("Module CHAZIYU: Make Requests Fail")
	}

	// 发送 HTTP 请求
	resp, err := client.Do(request)
	if err != nil {
		return []string{}, errors.New("Module CHAZIYU: Send Requests Error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return m_chaziyu.GetDomainNames(domain, retrycounts-1)
	}

	// 处理响应内容，提取域名目标
	subdomains, err := ResolveHTML_CHAZIYU(resp.Body)
	if err != nil {
		return []string{}, errors.New("Module CHAZIYU: Fail to extract subdomains")
	}
	// 这里就先不对切片结果数据去重了 因为不同的数据源之间有可能是有重复的
	// 所以等到所有数据源结果跑完 然后再进行统一去重
	// subdomains = pkg.RemoveDuplicates(subdomains)
	return subdomains, nil
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
