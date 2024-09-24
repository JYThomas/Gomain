package ProactiveDomain

import (
	"errors"
	"fmt"
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
func NewDNSResolver() *DNSResolver {
	return &DNSResolver{
		Resolver: &dns.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// QueryDNS 执行 DNS 查询
func (dr *DNSResolver) QueryDNS(domain string, recordType uint16) ([]dns.RR, error) {
	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(domain), recordType)
	m.RecursionDesired = true

	// 加载DNS服务器配置
	// 加载配置文件数据源url
	Config, err := pkg.LoadConfig()
	if err != nil {
		return nil, err
	}
	// 获取数据源爬虫目标链接
	DNSRESOLVER := Config.Section("ProactiveDomain").Key("DNSRESOLVER").String()

	response, _, err := dr.Resolver.Exchange(&m, DNSRESOLVER) // 使用 Google DNS 服务器
	if err != nil {
		return nil, err
	}
	if response.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("failed to get DNS records: %s", dns.RcodeToString[response.Rcode])
	}
	return response.Answer, nil
}

// 域名解析函数
func DomainResolution(domain string) (Records ResolutionResults) {

	var wg sync.WaitGroup
	ResultChan := make(chan []DNSRecord, 5)

	// 启动并发任务
	wg.Add(5)

	go func() {
		defer wg.Done()
		if IPv46Record, err := ResolutionIPv4(domain); err == nil {
			ResultChan <- IPv46Record
		} else {
			ResultChan <- nil
		}
	}()

	go func() {
		defer wg.Done()
		if IPv46Record, err := ResolutionIPv6(domain); err == nil {
			ResultChan <- IPv46Record
		} else {
			ResultChan <- nil
		}
	}()

	go func() {
		defer wg.Done()
		if IPv46Record, err := ResolutionCNAME(domain); err == nil {
			ResultChan <- IPv46Record
		} else {
			ResultChan <- nil
		}
	}()

	go func() {
		defer wg.Done()
		if MXRecord, err := ResolutionMX(domain); err == nil {
			ResultChan <- MXRecord
		} else {
			ResultChan <- nil
		}
	}()

	go func() {
		defer wg.Done()
		if NSRecord, err := ResolutionNS(domain); err == nil {
			ResultChan <- NSRecord
		} else {
			ResultChan <- nil
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
func ResolutionIPv4(domain string) (IPv4Record []DNSRecord, err error) {
	resolver := NewDNSResolver()

	// 使用 LookupIPContext 执行解析
	records, err := resolver.QueryDNS(domain, dns.TypeA)

	fmt.Println(records)
	fmt.Println(err)

	if err != nil {
		return nil, errors.New("Module ProactiveDomain: Fail to Resolve IPv4/IPv6")
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
	resolver := NewDNSResolver()

	// 使用 LookupIPContext 执行解析
	records, err := resolver.QueryDNS(domain, dns.TypeAAAA)

	fmt.Println(records)
	fmt.Println(err)

	if err != nil {
		return nil, errors.New("Module ProactiveDomain: Fail to Resolve IPv4/IPv6")
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
	resolver := NewDNSResolver()
	records, err := resolver.QueryDNS(domain, dns.TypeCNAME)

	fmt.Println(records)
	fmt.Println(err)

	if err != nil {
		return nil, errors.New("Module ProactiveDomain: Fail to Resolve IPv4/IPv6")
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
	resolver := NewDNSResolver()
	records, err := resolver.QueryDNS(domain, dns.TypeMX)

	fmt.Println(records)
	fmt.Println(err)

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
	resolver := NewDNSResolver()
	records, err := resolver.QueryDNS(domain, dns.TypeNS)

	fmt.Println(records)
	fmt.Println(err)

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
