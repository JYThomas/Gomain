package main

import (
	"fmt"
	"net/http"
	"time"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 切片去重函数
func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range elements {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}

// 设置请求客户端(注意声明返回值类型)
func make_http_client() (*http.Client){
	// 创建一个自定义的 HTTP 客户端，设置超时时间
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置超时时间为 10 秒
	}
	return client
}

// 设置请求配置
func make_req(url string)(*http.Request, error){
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

// 解析响应html页面
func resolve_html(html io.Reader)([]string, error){
	// 定义结果存储切片
	domains_slice := make([]string, 0)
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil{
		panic(err)
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

// 整合处理 请求 解析 返回
func query(domain string)([]string, error){
	// 目标url
	url := "https://chaziyu.com/" + domain

	// 创建请求客户端
	client := make_http_client()

	// 设置请求配置
	request, err := make_req(url)
	if err != nil{
		panic(err)
	}

	// 发送 HTTP 请求
	resp, err := client.Do(request)
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()

	// 处理响应内容，提取域名目标
	if(resp.StatusCode == 200){
		subdomains, err := resolve_html(resp.Body)
		if err != nil{
			panic(err)
		}
		result := removeDuplicates(subdomains)
		return result, nil
	}

	return nil, err
}

// 主函数测试调用
func main(){
	domain := "bilibili.com"
	subdomains, err := query(domain)
	if err != nil{
		panic(err)
	}
	fmt.Println(subdomains)
}
