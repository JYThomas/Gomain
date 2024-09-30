package ProactiveDomain

// import (
// 	"fmt"
// 	"sync"

// 	"github.com/miekg/dns"
// )

// type DNSRecord struct {
// 	DomainName string
// 	RecordType string
// 	Records    []string
// }

// func resolveDNS(domain string, recordType uint16, wg *sync.WaitGroup, results chan<- DNSRecord) {
// 	defer wg.Done()

// 	client := new(dns.Client)
// 	msg := new(dns.Msg)
// 	msg.SetQuestion(dns.Fqdn(domain), recordType)
// 	msg.RecursionDesired = true

// 	response, _, err := client.Exchange(msg, "8.8.8.8:53") // 使用 Google 的公共 DNS
// 	if err != nil {
// 		fmt.Printf("Failed to resolve %s: %v\n", domain, err)
// 		return
// 	}

// 	var records []string
// 	switch recordType {
// 	case dns.TypeA:
// 		for _, ans := range response.Answer {
// 			if aRecord, ok := ans.(*dns.A); ok {
// 				records = append(records, aRecord.A.String())
// 			}
// 		}
// 	case dns.TypeAAAA:
// 		for _, ans := range response.Answer {
// 			if aaaaRecord, ok := ans.(*dns.AAAA); ok {
// 				records = append(records, aaaaRecord.AAAA.String())
// 			}
// 		}
// 	case dns.TypeCNAME:
// 		for _, ans := range response.Answer {
// 			if cnameRecord, ok := ans.(*dns.CNAME); ok {
// 				records = append(records, cnameRecord.Target)
// 			}
// 		}
// 	case dns.TypeMX:
// 		for _, ans := range response.Answer {
// 			if mxRecord, ok := ans.(*dns.MX); ok {
// 				records = append(records, mxRecord.Mx)
// 			}
// 		}
// 	case dns.TypeNS:
// 		for _, ans := range response.Answer {
// 			if nsRecord, ok := ans.(*dns.NS); ok {
// 				records = append(records, nsRecord.Ns)
// 			}
// 		}
// 	default:
// 		fmt.Printf("Unsupported record type: %d\n", recordType)
// 		return
// 	}

// 	results <- DNSRecord{DomainName: domain, RecordType: dns.TypeToString[recordType], Records: records}
// }

// func BatchResolveDNS(domains []string, recordTypes []uint16, concurrencyLimit int) map[string]map[string][]string {
// 	var wg sync.WaitGroup
// 	results := make(chan DNSRecord)
// 	resultMap := make(map[string]map[string][]string)
// 	sem := make(chan struct{}, concurrencyLimit) // 控制并发数

// 	for _, domain := range domains {
// 		for _, recordType := range recordTypes {
// 			wg.Add(1)
// 			sem <- struct{}{} // 获取信号量
// 			go func(domain string, recordType uint16) {
// 				defer func() { <-sem }() // 释放信号量
// 				resolveDNS(domain, recordType, &wg, results)
// 			}(domain, recordType)
// 		}
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(results)
// 	}()

// 	for result := range results {
// 		if resultMap[result.DomainName] == nil {
// 			resultMap[result.DomainName] = make(map[string][]string)
// 		}
// 		resultMap[result.DomainName][result.RecordType] = result.Records
// 	}

// 	return resultMap
// }
