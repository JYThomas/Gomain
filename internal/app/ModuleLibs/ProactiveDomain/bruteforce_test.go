package ProactiveDomain

import (
	"fmt"
	"testing"
)

func TestProactiveDomain(t *testing.T) {
	// 初始化测试模块
	bf := MODULE_BRUTEFORCE{ModuleName: "bruteforce"}
	domain := "gxu.edu.cn"
	dict_name := "subnames_ofa.txt"

	// 调用要测试的函数
	subdomains, err := bf.GetDomainNames(domain, dict_name)
	if err != nil {
		fmt.Println(err)
	}
	// 断言: 检查返回的子域名列表是否为空
	// if len(records.DNSRecord) == 0 {
	// 	t.Errorf("Expected subdomains, got empty result")
	// }

	t.Logf("\nrecords: %v", subdomains)
	t.Logf("\nrecords: %v", len(subdomains))

}
