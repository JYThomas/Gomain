package main

import (
	"fmt"
	"io"
	"math"
	"errors"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/JYThomas/Gomain/internal/utils/MakeRequests"
	// "Gomain/internal/utils/HandleFunc"
)

// 爬虫请求函数
func query(domain string)(io.Reader, error){
	// 目标url
	url := "https://www.rapiddns.io/s/" + domain

	// 创建请求客户端
	client := MakeRequests.MakeHttpClient()

	// 设置请求配置
	request, err := MakeRequests.MakeReq(url)
	if err != nil{
		return []string{}, errors.New("Fail to create requests")
	}

	// 发送 HTTP 请求
	resp, err := client.Do(request)
	if err != nil{
		return []string{}, errors.New("Fail to send requests")
	}
	defer resp.Body.Close()

	// 处理响应内容，提取域名目标
	if(resp.StatusCode != 200){
		return []string{}, errors.New("Response Error: not 200")
	}

	return resp.Body, nil
}

// 获取查询结果分页信息
func GetRecordNumber(html io.Reader)(int, int, error){
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return []string{}, errors.New("Fail to read HTML")
	}
	// 获取总记录数(total标签)
	RecordNumber := 0
	fmt.Println(doc)
	doc.Find("div.d-flex").Each(func(i int, s *goquery.Selection){
		// 获取span元素中的文本、
		spanText := s.Find("span").Text()
		fmt.Println(spanText)
		RecordNumber, err = strconv.Atoi(spanText)
		if err != nil {
			return []string{}, errors.New("Fail to extract subdomains")
		}
	})
	PageNumber := int(math.Ceil(float64(RecordNumber) / 100.0))

	return RecordNumber, PageNumber, nil

}

// 页面数据解析函数
// func resolve_html(html io.Reader)([]string, error){
// 	// 定义结果存储切片
// 	// domain_slice := make([]string, 0)
// 	doc, err := goquery.NewDocumentFromReader(html)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// 解析HTML
	
// }


// 主函数调用测试
func main(){
	domain := "bilibili.com"
	// 首先获取响应内容第一页。获取总记录数
	FirstPageDoc, err := query(domain)
	if err != nil {
		panic(err)
	}
	// 提取记录数
	RecordNumber, PageNumber, err := GetRecordNumber(FirstPageDoc)
	if err != nil {
		panic(err)
	}
	fmt.Println(RecordNumber, PageNumber)
	// 
}





