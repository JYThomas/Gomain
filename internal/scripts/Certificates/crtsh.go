package main

import (
	"fmt"
	"io"
	"encoding/json"
	"Gomain/internal/utils/MakeRequests"
	"Gomain/internal/utils/HandleFunc"

)

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

func query(domain string)([]string, error){
	url := "https://crt.sh/?q=" + domain
	
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

	if(resp.StatusCode != 200){
		return []string{}, nil
	}

	// 处理响应内容，提取域名目标
	subdomains, err := resolve_resp(resp.Body)
	if err != nil{
		panic(err)
	}
	result := HandleFunc.RemoveDuplicates(subdomains)
	return result, err
}

// 解析响应的html页面
func resolve_resp(html io.Reader)(result []string, err error){
	// 读取响应体 响应内容为json字符串
	content, err := io.ReadAll(html)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var certificates []Certificate
	// 将Json字符串内容解析为Certificate结构体切片 Json字符串内容与结构体结构相同
	err = json.Unmarshal([]byte(content), &certificates)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	for _, cert := range certificates {
		result = append(result, cert.CommonName)
	}

	return result, nil
}


func main(){
	// domain := "bilibili.com"
	domain := "gxust.edu.cn"
	subdomains, err := query(domain)
	if err != nil{
		panic(err)
	}
	fmt.Println(subdomains)
}
