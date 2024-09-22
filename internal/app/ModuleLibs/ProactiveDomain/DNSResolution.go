package ProactiveDomain

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

type ResolutionResults struct {
	DomainName string
	DNSRecord  []DNSRecord
}

type DNSRecord struct {
	DomainName  string
	ResolveType string
	Host        string
}

// 域名解析函数
func DomainResolution(domain string) (Records ResolutionResults) {
	var wg sync.WaitGroup
	ResultChan := make(chan []DNSRecord, 3)

	// 启动并发任务
	wg.Add(3)

	// IPv4、IPv6
	go func() {
		defer wg.Done()
		IPv46Record, err := ResolutionIPv46(domain)
		if err == nil {
			ResultChan <- IPv46Record
		}
	}()

	// MX
	go func() {
		defer wg.Done()
		MXRecord, err := ResolutionMX(domain)
		if err == nil {
			ResultChan <- MXRecord
		}
	}()

	// NS
	go func() {
		defer wg.Done()
		NSRecord, err := ResolutionNS(domain)
		if err == nil {
			ResultChan <- NSRecord
		}
	}()

	// 等待所有协程完成并关闭通道
	go func() {
		wg.Wait()
		close(ResultChan)
	}()

	// 收集所有结果
	var allRecords []DNSRecord
	for records := range ResultChan {
		allRecords = append(allRecords, records...)
	}

	// 如果没有结果 说明生成的域名是无效域名
	return ResolutionResults{
		DomainName: domain,
		DNSRecord:  allRecords,
	}
}

// IPv4、IPv6
func ResolutionIPv46(domain string) (IPv46Record []DNSRecord, err error) {
	// 使用 context 设置超时时间，例如5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用 LookupIPContext 执行解析
	ips, err := net.DefaultResolver.LookupIP(ctx, "ip", domain)

	if err != nil {
		return nil, errors.New("Module ProactiveDomain: Fail to Resolve IPv4/IPv6")
	}
	for _, ip := range ips {
		resolveType := "A"
		if ip.To4() == nil {
			resolveType = "AAAA" // 如果不是 IPv4，则为 IPv6
		}
		record := DNSRecord{
			DomainName:  domain,
			ResolveType: resolveType,
			Host:        ip.String(),
		}
		IPv46Record = append(IPv46Record, record)
	}

	return IPv46Record, nil
}

// MX
func ResolutionMX(domain string) (MXRecord []DNSRecord, err error) {
	// 使用 context 设置超时时间，例如5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用 LookupIPContext 执行解析
	mxs, err := net.DefaultResolver.LookupMX(ctx, domain)

	if err != nil {
		return nil, errors.New("Module ProactiveDomain: Fail to Resolve MX")
	}
	for _, mx := range mxs {
		record := DNSRecord{
			DomainName:  domain,
			ResolveType: "MX",
			Host:        mx.Host,
		}
		MXRecord = append(MXRecord, record)
	}

	return MXRecord, nil
}

// NS
func ResolutionNS(domain string) (NSRecord []DNSRecord, err error) {
	// 使用 context 设置超时时间，例如5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用 LookupIPContext 执行解析
	nss, err := net.DefaultResolver.LookupNS(ctx, domain)

	if err != nil {
		return nil, errors.New("Module ProactiveDomain: Fail to Resolve NS")
	}
	for _, ns := range nss {
		record := DNSRecord{
			DomainName:  domain,
			ResolveType: "NS",
			Host:        ns.Host,
		}
		NSRecord = append(NSRecord, record)
	}

	return NSRecord, nil
}
