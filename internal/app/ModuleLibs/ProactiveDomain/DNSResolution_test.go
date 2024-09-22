package ProactiveDomain

import (
	"testing"
)

func TestDNSResolution(t *testing.T) {
	domain := "www.bilibili.com"

	// 调用要测试的函数
	records := DomainResolution(domain)

	// 断言: 检查返回的子域名列表是否为空
	// if len(records.DNSRecord) == 0 {
	// 	t.Errorf("Expected subdomains, got empty result")
	// }

	t.Logf("\nrecords: %v", records)
}
