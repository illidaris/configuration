package sentinel

import (
	"errors"

	baseSentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/illidaris/configuration"
)

var (
	ErrDefaultCenterNil = errors.New("DefaultCenter is nil")
	ErrNacosConnFailed  = errors.New("NacosConnFailed")
)

func InitSentinel() error {
	if configuration.DefaultCenter == nil {
		return ErrDefaultCenterNil
	}
	var (
		baseConfig     = configuration.DefaultCenter.GetServiceInfo()
		sentinelConfig = baseConfig.Sentinel
		client         = configuration.DefaultCenter.GetConfigClient()
	)
	// 初始化Sentinal
	baseSentinel.InitWithConfigFile(sentinelConfig.ExtConfigFile)
	//注册流控规则Handler
	h := datasource.NewFlowRulesHandler(datasource.FlowRuleJsonArrayParser)
	//创建NacosDataSource数据源
	nds, err := NewNacosDataSource(client, sentinelConfig.GroupName, sentinelConfig.ServiceName, h)
	if err != nil {
		return err
	}
	//nacos数据源初始化
	err = nds.Initialize()
	if err != nil {
		return err
	}
	return nil
}
