package Task

type RunDomainTask interface {
	ProactiveCollection() []string
	PassiveCollection() []string
	DNSResolution(subdomains []string) []DNSInfo
	ICPScrapy(subdomains []string) []ICPInfo
	WebServiceCollection(urls []string) []WebInfo
	HostInfoCollection(ips []string) []HostInfo
}

type RunIPTask interface {
	HostInfoCollection(ips []string) []HostInfo
	WebServiceCollection(urls []string) []WebInfo
}
