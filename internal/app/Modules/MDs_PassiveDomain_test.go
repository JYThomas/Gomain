package Modules

import (
	"fmt"
	"testing"
)

func TestMDs_PassiveDomain(t *testing.T) {
	// 定义测试域名
	domain := "iqiyi.com"

	// 重试次数
	// retrycounts := 3

	// 调用要测试的函数
	subdomains, err := RunPassiveDomain(domain)

	fmt.Println(subdomains)

	// 断言: 检查错误是否为 nil
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 断言: 检查返回的子域名列表是否为空
	if len(subdomains) == 0 {
		t.Errorf("Expected subdomains, got empty result")
	}

	t.Logf("Subdomains: %v", subdomains)
}
