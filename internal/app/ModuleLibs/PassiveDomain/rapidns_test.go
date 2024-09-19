package PassiveDomain

import (
	"testing"
)

func TestGetDomainNames_RAPIDDNS(t *testing.T) {
	// 初始化测试模块
	rapiddns := MODULE_RAPIDDNS{ModuleName: "rapiddns"}

	// 定义测试域名
	domain := "bilibili.com"
	retrycounts := 2

	// 调用要测试的函数
	subdomains, err := rapiddns.GetDomainNames(domain, retrycounts)

	// 断言: 检查错误是否为 nil
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 断言: 检查返回的子域名列表是否为空
	if len(subdomains) == 0 {
		t.Errorf("Expected subdomains, got empty result")
	}

	// 打印结果以供手动检查
	t.Logf("Subdomains: %v", subdomains)
}
