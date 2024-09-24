package ProactiveDomain

// import (
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/miekg/dns"
// )

// // DNSResolver 结构体用于 DNS 解析
// type DNSResolver struct {
// 	Resolver *dns.Client
// }

// // NewDNSResolver 创建一个新的 DNSResolver
// func NewDNSResolver() *DNSResolver {
// 	return &DNSResolver{
// 		Resolver: &dns.Client{
// 			Timeout: 5 * time.Second,
// 		},
// 	}
// }

// // QueryDNS 执行 DNS 查询
// func (dr *DNSResolver) QueryDNS(domain string, recordType uint16) ([]dns.RR, error) {
// 	m := dns.Msg{}
// 	m.SetQuestion(dns.Fqdn(domain), recordType)
// 	m.RecursionDesired = true

// 	response, _, err := dr.Resolver.Exchange(&m, "8.8.8.8:53") // 使用 Google DNS 服务器
// 	if err != nil {
// 		return nil, err
// 	}
// 	if response.Rcode != dns.RcodeSuccess {
// 		return nil, fmt.Errorf("failed to get DNS records: %s", dns.RcodeToString[response.Rcode])
// 	}
// 	return response.Answer, nil
// }

// func dnsres(domain string) {
// 	resolver := NewDNSResolver()

// 	// 查询不同类型的 DNS 记录
// 	recordTypes := map[string]uint16{
// 		"A":     dns.TypeA,
// 		"AAAA":  dns.TypeAAAA,
// 		"CNAME": dns.TypeCNAME,
// 		"MX":    dns.TypeMX,
// 		"NS":    dns.TypeNS,
// 	}

// 	for name, recordType := range recordTypes {
// 		records, err := resolver.QueryDNS(domain, recordType) // 替换为要查询的域名
// 		if err != nil {
// 			log.Fatalf("Error querying %s records for example.com: %v", name, err)
// 		}
// 		fmt.Printf("%s records for example.com:\n", name)
// 		for _, record := range records {
// 			fmt.Println(record)
// 		}
// 	}
// }
