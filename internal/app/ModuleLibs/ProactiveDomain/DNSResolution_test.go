package ProactiveDomain

import (
	"fmt"
	"testing"
)

func TestDNSResolution(t *testing.T) {
	// domain := "djasfao.gxust.edu.cn"
	domain := "www.gxust.edu.cn"
	// res1, err1 := ResolutionIPv4(domain)
	// res2, err2 := ResolutionIPv6(domain)
	// res3, err3 := ResolutionCNAME(domain)
	// res4, err4 := ResolutionMX(domain)
	// res5, err5 := ResolutionNS(domain)
	// if err1 != nil {
	// 	fmt.Println(err1)
	// }
	// fmt.Println(res1)

	// if err2 != nil {
	// 	fmt.Println(err2)
	// }
	// fmt.Println(res2)

	// if err3 != nil {
	// 	fmt.Println(err3)
	// }
	// fmt.Println(res3)

	// if err4 != nil {
	// 	fmt.Println(err4)
	// }
	// fmt.Println(res4)

	// if err5 != nil {
	// 	fmt.Println(err5)
	// }
	// fmt.Println(res5)

	result := DomainResolution(domain)
	fmt.Println(result)

}
