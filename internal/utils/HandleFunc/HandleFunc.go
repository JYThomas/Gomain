package HandleFunc
// package main
import (
	"fmt"
	"strings"
)


// 切片去重函数
func RemoveDuplicates(elements []string) []string {
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

// 去除泛解析格式的域名 *.baidu.com
func FilterWildcardSign(domain string)(domain_extracted string){
	if strings.Contains(domain, "*."){
		parts := strings.Split(domain, "*.")
		domain_extracted := strings.Join(parts[1:], ".")
		return domain_extracted
	}

	return domain
}

// func main(){
// 	res := FilterWildcardSign("baidu.com")
// 	fmt.Println(res)
// }