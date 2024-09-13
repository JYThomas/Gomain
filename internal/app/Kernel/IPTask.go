package Task

// IP扫描任务
type IPTask struct {
	UserId     int64
	TaskId     int64
	TaskName   string
	TaskType   int64
	Targes     string
	Strategies IPStrategies
}

// IP扫描任务配置
type IPStrategies struct {
	AliveDetectMode       bool
	PortScanMode          bool
	PortServiceDetectMode bool
	IP2DOMAINMode         bool
}

// 主机信息
type HostInfo struct {
	Domain          string
	IP              string
	PortServiceInfo PortServiceInfo
	IP2Domains      []string
	OSInfo          string
}

// 端口服务信息
type PortServiceInfo struct {
	IP            string
	Port          int64
	ServiceName   string
	ServiceFinger string
}

// Web服务信息收集
func (it IPTask) WebServiceCollection(urls []string) (web_info []WebInfo) {
	// 通过收集结果子域 收集web服务资产信息
	return
}

// IP资产收集
func (it IPTask) HostInfoCollection(ips []string) (host_info []HostInfo) {
	// 通过收集结果子域 收集主机服务资产信息
	return
}
