package ProactiveDomain

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/miekg/dns"
// )

// func TestDNSresolve(t *testing.T) {
// 	SubdomainNames, err := GetSuspectDomainNames("gxust.edu.cn", "subnames_ofa.txt")
// 	if err != nil {
// 		fmt.Println("ProactiveDomain Module: Load Suspect Domain Names Fail")
// 	}

// 	recordTypes := []uint16{dns.TypeA, dns.TypeAAAA, dns.TypeCNAME, dns.TypeMX, dns.TypeNS}

// 	results := BatchResolveDNS(SubdomainNames[:1000], recordTypes, 50)

// 	for domain, records := range results {
// 		fmt.Printf("Domain: %s\n", domain)
// 		for recordType, values := range records {
// 			fmt.Printf("  %s: %v\n", recordType, values)
// 		}
// 	}
// }
