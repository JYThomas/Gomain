package Task

// 域名扫描任务
type DomainTask struct {
	UserId     int64
	TaskId     int64
	TaskName   string
	TaskType   int64
	Targets    string
	Strategies DomainStrategies
}

// 域名扫描任务配置
type DomainStrategies struct {
	ProactiveMode            bool
	PassiveMode              bool
	AltDNSMode               bool
	WebServiceDetectMode     bool
	WebFingerRecognitionMode bool
	SiteScrapyMode           bool
	SiteScreenShotMode       bool
	ICPScrapyMode            bool
	HostScanMode             IPStrategies
}

// 定义DNS解析结果结构
type DNSInfo struct {
	Domain     string
	RecordType string
	IP         string
}

// ICP备案信息
type ICPInfo struct {
	Domain         string
	ICPNumber      string
	Department     string
	DepartmentType string
}

// Web资产信息
type WebInfo struct {
	Domain         string
	URL            string
	Title          string
	StatusCode     int64
	Headers        string
	WebFingers     string
	ScreenShot     string
	SiteScrapyInfo SiteScrapyInfo
}

// 站点爬虫信息
type SiteScrapyInfo struct {
	Domain     string
	URL        string
	Title      string
	StatusCode int64
	Headers    string
}

// 被动收集
func (dt DomainTask) PassiveCollection() (subdomains []string) {
	// TargetDomains := strings.Split(dt.Targets, ",")
	// 使用协程并发执行多个模块的资产情报收集
	return
}

// DNS解析结果
func (dt DomainTask) DNSResolution(subdomains []string) (dns_info []DNSInfo) {
	// 通过收集结果子域切片 获取域名DNS解析结果
	return
}

// ICP备案信息查询
func (dt DomainTask) ICPScrapy(subdomains []string) (icp_info []ICPInfo) {
	// 通过收集结果子域 获取ICP备案信息
	return
}

// Web服务信息收集
func (dt DomainTask) WebServiceCollection(urls []string) (web_info []WebInfo) {
	// 通过收集结果子域 收集web服务资产信息
	return
}

// IP资产收集
func (dt DomainTask) HostInfoCollection(ips []string) (host_info []HostInfo) {
	// 通过收集结果子域 收集主机服务资产信息
	return
}
