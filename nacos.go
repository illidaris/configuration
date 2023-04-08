package configuration

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

var _ = IConfigurationCenter(&NacosCenter{})

type NacosCenter struct {
	ServiceInfo
	RealPort int

	ConfigClient config_client.IConfigClient
	NamingClient naming_client.INamingClient

	logger     ILogger
	settingMap map[string]map[string]*viper.Viper
}

func (n *NacosCenter) SetRealPort(port int) {
	n.RealPort = port
}

func (n *NacosCenter) GetPort() uint64 {
	if n.Port == 0 {
		return uint64(n.RealPort)
	}
	return n.Port
}

func NewNacos(param vo.NacosClientParam, serv ServiceInfo) (IConfigurationCenter, error) {
	if serv.IP == "" {
		serv.IP = GetIPX()
	}
	n := &NacosCenter{
		ServiceInfo: serv,
		logger:      &DefaultLogger{},
		settingMap:  map[string]map[string]*viper.Viper{},
	}
	// create config client
	client, err := clients.NewConfigClient(param)
	if err != nil {
		return nil, err
	}
	nClient, err := clients.NewNamingClient(param)
	if err != nil {
		return nil, err
	}
	n.ConfigClient = client
	n.NamingClient = nClient
	return n, nil
}

func NewSimpleNacos(cfg *SimpleConfig) (IConfigurationCenter, error) {
	p, err := cfg.ToParam()
	if err != nil {
		return nil, err
	}
	return NewNacos(p, cfg.Service)
}
