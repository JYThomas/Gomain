package Modules

import (
	"errors"
	"fmt"
	"strings"

	"github.com/JYThomas/Gomain/internal/app/ModuleLibs/PassiveDomain"
	"github.com/JYThomas/Gomain/internal/pkg"
)

// 定义一个接口用于后续结构体实例化的方法实现
type DomainCollector interface {
	GetDomainNames(domain string, param int) ([]string, error)
}

// 根据名称返回模块结构体的指针
func GetModuleByName(ModuleName string) DomainCollector {
	switch ModuleName {
	case "CHAZIYU":
		return &PassiveDomain.MODULE_CHAZIYU{ModuleName: "chaziyu"}
	case "CRTSH":
		return &PassiveDomain.MODULE_CRTSH{ModuleName: "crtsh"}
	case "RAPIDDNS":
		return &PassiveDomain.MODULE_RAPIDDNS{ModuleName: "rapiddns"}
	default:
		return nil
	}
}

// 从配置文件加载被动收集模块名称
func RunPassiveDomain(domain string) (PassiveDomainNames []string, err error) {
	// 加载配置文件
	Config, err := pkg.LoadConfig()
	if err != nil {
		return []string{}, errors.New("Module - PassiveDomain Main Process: Fail to load config file")
	}

	// 根据模块名称加载对应模块的被动收集函数脚本
	PassiveModulesString := Config.Section("PassiveDomain").Key("MODULES_NAME").String()
	PassiveModulesString = strings.Trim(PassiveModulesString, "[]")

	// 字符串数组
	PassiveModules := strings.Split(PassiveModulesString, ",")
	// 结构体指针数组存储结构体实例的指针
	ObjPassiveModules := []DomainCollector{}

	// 模块结构体实例化
	for _, ModuleName := range PassiveModules {
		ObjModules := GetModuleByName(ModuleName)
		if ObjModules != nil {
			ObjPassiveModules = append(ObjPassiveModules, ObjModules)
		}
	}

	fmt.Println(ObjPassiveModules)

	// 分配协程执行被动域名资产收集
	PassiveDomainNames = []string{}
	for _, objModule := range ObjPassiveModules {
		fmt.Println(objModule)
		Domains, err := objModule.GetDomainNames(domain, 3)
		if err != nil {
			continue
		}
		PassiveDomainNames = append(PassiveDomainNames, Domains...)
	}

	// 域名资产收集结果返回 统一好模块结果格式
	return PassiveDomainNames, nil
}
