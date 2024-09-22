package ProactiveDomain

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/JYThomas/Gomain/internal/pkg"
)

// 定义模块结构体
type MODULE_BRUTEFORCE struct {
	ModuleName string
}

// 获取子域字典
func GetSuspectDomainNames(domain string, dict_name string) (SubdomainNames []string, err error) {
	// 获取当前工作目录
	startPath, err := filepath.Abs(".")
	if err != nil {
		return nil, errors.New("Load Config File Error: fail to find current File")
	}

	// 查找Gomain目录
	gomainDir, err := pkg.FindGomainDir(startPath)
	if err != nil {
		return nil, fmt.Errorf("Load Config File Error: Fail to find Gomain dir: %v", err)
	}

	// 字典相对路径
	configRelPath := "internal/app/ModuleLibs/ProactiveDomain/dicts/" + dict_name
	configPath := filepath.Join(gomainDir, configRelPath)

	// 打开文件
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("Error opening file %s: %v", configPath, err)
	}
	defer file.Close()

	// 存储读取出来的域名字典
	var SubLevelDomain []string

	// 使用 bufio.NewScanner 按行读取文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		SubLevelDomain = append(SubLevelDomain, line)
	}

	// 检查扫描过程中是否有错误
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file %s: %v", configPath, err)
	}

	// 生成域名目标
	for _, SubLevelName := range SubLevelDomain {
		SubdomainNames = append(SubdomainNames, SubLevelName+"."+domain)
	}

	return SubdomainNames, nil
}

// DNS解析域名
func (bf MODULE_BRUTEFORCE) GetDomainNames(domain string, dict_name string) (DomainNames []string, err error) {
	SubdomainNames, err := GetSuspectDomainNames(domain, dict_name)
	if err != nil {
		return []string{}, errors.New("ProactiveDomain Module: Load Suspect Domain Names Fail")
	}

	// 创建一个用于存储解析结果的通道
	resultChan := make(chan string, len(SubdomainNames))

	// 并发解析潜在域名
	var wg sync.WaitGroup

	// 并发解析潜在域名
	for _, subdomain := range SubdomainNames {
		wg.Add(1)
		go func(sub string) {
			defer wg.Done()
			// 执行域名解析
			results := DomainResolution(subdomain)

			if results.DNSRecord != nil {
				// 发送解析结果到通道
				resultChan <- subdomain
			}

		}(subdomain)
	}

	// 等待所有 goroutines 完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集所有解析结果 DNS解析成功的域名
	for resolved := range resultChan {
		DomainNames = append(DomainNames, resolved)
	}

	return DomainNames, nil
}
