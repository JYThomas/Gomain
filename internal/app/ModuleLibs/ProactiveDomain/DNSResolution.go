package ProactiveDomain

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/JYThomas/Gomain/internal/pkg"
	"github.com/miekg/dns"
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

// DNSResolver 结构体用于 DNS 解析
type DNSResolver struct {
	Resolver *dns.Client
}

// NewDNSResolver 创建一个新的 DNSResolver
func CreateDNSResolver() *DNSResolver {
	return &DNSResolver{
		Resolver: &dns.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// 加载DNS服务器配置
func LoadDNSServer() (DNSRESOLVER string) {
	// 加载DNS服务器配置
	// 默认DNS服务器, 加载配置文件出错的情况下返回默认的DNS服务器配置
	DNSServer := "8.8.8.8:53"
	Config, err := pkg.LoadConfig()
	if err != nil {
		return DNSServer
	}
	// 获取数据源爬虫目标链接
	DNSRESOLVER = Config.Section("ProactiveDomain").Key("DNSRESOLVER").String()
	return DNSRESOLVER
}

// QueryDNS 执行 DNS 查询
func (dr *DNSResolver) QueryDNS(domain string, recordType uint16) ([]dns.RR, error) {
	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(domain), recordType)
	m.RecursionDesired = true

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 加载DNS服务器
	DNSRESOLVER := LoadDNSServer()

	// DNS请求
	response, _, err := dr.Resolver.ExchangeContext(ctx, &m, DNSRESOLVER)
	// fmt.Println(err)

	if err != nil {
		return nil, errors.New("Fail to resolve")
	}
	if response.Rcode != dns.RcodeSuccess {
		return nil, errors.New("failed to get DNS records")
	}

	return response.Answer, nil
}

// 域名解析函数
func DomainResolution(domain string) (Records ResolutionResults) {

	var wg sync.WaitGroup
	ResultChan := make(chan []DNSRecord, 5)

	// 启动并发任务
	wg.Add(5)

	// IPv4
	go func(d string) {
		defer wg.Done()
		if IPv4Record, err := ResolutionIPv4(d); err == nil {
			ResultChan <- IPv4Record
		} else {
			ResultChan <- []DNSRecord{}
		}
	}(domain)

	// IPv6
	go func(d string) {
		defer wg.Done()
		if IPv6Record, err := ResolutionIPv6(d); err == nil {
			ResultChan <- IPv6Record
		} else {
			ResultChan <- []DNSRecord{}
		}
	}(domain)

	// CNAME
	go func(d string) {
		defer wg.Done()
		if CNAMERecord, err := ResolutionCNAME(d); err == nil {
			ResultChan <- CNAMERecord
		} else {
			ResultChan <- []DNSRecord{}
		}
	}(domain)

	// MX
	go func(d string) {
		defer wg.Done()
		if MXRecord, err := ResolutionMX(d); err == nil {
			ResultChan <- MXRecord
		} else {
			ResultChan <- []DNSRecord{}
		}
	}(domain)

	// NS
	go func(d string) {
		defer wg.Done()
		if NSRecord, err := ResolutionNS(d); err == nil {
			ResultChan <- NSRecord
		} else {
			ResultChan <- []DNSRecord{}
		}
	}(domain)

	// 等待所有协程完成并关闭通道
	go func() {
		wg.Wait()
		close(ResultChan)
	}()

	// 收集所有结果
	var allRecords []DNSRecord
	for records := range ResultChan {
		if records != nil {
			allRecords = append(allRecords, records...)
		}
	}

	// 如果没有结果 说明生成的域名是无效域名
	return ResolutionResults{
		DomainName: domain,
		DNSRecord:  allRecords,
	}
}

// IPv4
func ResolutionIPv4(domain string) (IPv4Record []DNSRecord, err error) {
	resolver := CreateDNSResolver()
	// 使用 LookupIPContext 执行解析
	records, err := resolver.QueryDNS(domain, dns.TypeA)

	if err != nil {
		return nil, errors.New("Module ProactiveDomain: Fail to Resolve IPv4")
	}
	for _, record := range records {
		if aRecord, ok := record.(*dns.A); ok {
			record := DNSRecord{
				DomainName:  domain,
				ResolveType: "A",
				Host:        aRecord.A.String(),
			}
			IPv4Record = append(IPv4Record, record)
		}
	}
	return IPv4Record, nil
}

// IPv6
func ResolutionIPv6(domain string) (IPv6Record []DNSRecord, err error) {
	resolver := CreateDNSResolver()
	// 使用 LookupIPContext 执行解析
	records, err := resolver.QueryDNS(domain, dns.TypeAAAA)

	if err != nil {
		return nil, errors.New("Module ProactiveDomain: Fail to Resolve IPv6")
	}
	for _, record := range records {
		if aaaaRecord, ok := record.(*dns.AAAA); ok {
			record := DNSRecord{
				DomainName:  domain,
				ResolveType: "AAAA",
				Host:        aaaaRecord.AAAA.String(),
			}
			IPv6Record = append(IPv6Record, record)
		}
	}

	return IPv6Record, nil
}

// CNAME
func ResolutionCNAME(domain string) (CNAMERecord []DNSRecord, err error) {
	resolver := CreateDNSResolver()
	records, err := resolver.QueryDNS(domain, dns.TypeCNAME)

	if err != nil {
		return nil, errors.New("Module ProactiveDomain: Fail to Resolve CNAME")
	}
	for _, record := range records {
		if cnameRecord, ok := record.(*dns.CNAME); ok {
			record := DNSRecord{
				DomainName:  domain,
				ResolveType: "CNAME",
				Host:        cnameRecord.Target,
			}
			CNAMERecord = append(CNAMERecord, record)
		}
	}

	return CNAMERecord, nil
}

// MX
func ResolutionMX(domain string) (MXRecord []DNSRecord, err error) {
	resolver := CreateDNSResolver()
	records, err := resolver.QueryDNS(domain, dns.TypeMX)

	if err != nil {
		return nil, errors.New("Module ProactiveDomain: Fail to Resolve MX")
	}
	for _, record := range records {
		if mxRecord, ok := record.(*dns.MX); ok {
			record := DNSRecord{
				DomainName:  domain,
				ResolveType: "MX",
				Host:        mxRecord.Mx,
			}
			MXRecord = append(MXRecord, record)
		}
	}

	return MXRecord, nil
}

// NS
func ResolutionNS(domain string) (NSRecord []DNSRecord, err error) {
	resolver := CreateDNSResolver()
	records, err := resolver.QueryDNS(domain, dns.TypeNS)

	if err != nil {
		return nil, errors.New("Module ProactiveDomain: Fail to Resolve NS")
	}
	for _, record := range records {
		if nsRecord, ok := record.(*dns.NS); ok {
			record := DNSRecord{
				DomainName:  domain,
				ResolveType: "NS",
				Host:        nsRecord.Ns,
			}
			NSRecord = append(NSRecord, record)
		}
	}

	return NSRecord, nil
}
