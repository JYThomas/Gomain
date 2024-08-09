package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"Gomain/internal/utils/MakeRequests"
	"Gomain/internal/utils/HandleFunc"
)

// 整合处理 请求 解析 返回
func query(domain string)([]string, error){
	// 目标url
	url := "https://chaziyu.com/" + domain

	// 创建请求客户端
	client := MakeRequests.MakeHttpClient()

	// 设置请求配置
	request, err := MakeRequests.MakeReq(url)
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
		result := HandleFunc.RemoveDuplicates(subdomains)
		return result, nil
	}

	return nil, err
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


// 主函数测试调用
func main(){
	domain := "bilibili.com"
	subdomains, err := query(domain)
	if err != nil{
		panic(err)
	}
	fmt.Println(subdomains)
}
