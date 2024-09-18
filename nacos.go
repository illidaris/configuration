package configuration

import (
	"io/fs"
	"os"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var _ = IConfigurationCenter(&NacosCenter{})
var _ = IDiscoverCenter(&NacosCenter{})

type NacosCenter struct {
	ServiceInfo
	RealPort uint64

	ConfigClient config_client.IConfigClient
	NamingClient naming_client.INamingClient

	logger     ILogger
	settingMap map[string]map[string]*viper.Viper
}

func (n *NacosCenter) SetILogger(log ILogger) {
	n.logger = log
}

func (n *NacosCenter) GetConfigClient() config_client.IConfigClient {
	return n.ConfigClient
}

func (n *NacosCenter) GetNamingClient() naming_client.INamingClient {
	return n.NamingClient
}

func (n *NacosCenter) GetServiceName() string {
	return n.ServiceName
}

func (n *NacosCenter) GetServiceInfo() ServiceInfo {
	return n.ServiceInfo
}

func (n *NacosCenter) SetRealPort(port int) error {
	if n.RealPort != uint64(port) {
		n.RealPort = uint64(port)
		return os.WriteFile(KEY_TMP_PORT_PATH, []byte(cast.ToString(n.RealPort)), fs.ModePerm)
	}
	return nil
}

func (n *NacosCenter) GetPort() uint64 {
	if n.Port == 0 { // random port
		return n.RealPort
	}
	return n.Port
}

func (n *NacosCenter) GetIP() string {
	return n.IP
}

func NewNacos(param vo.NacosClientParam, serv ServiceInfo) (IConfigurationCenter, error) {
	if serv.IP == "" {
		serv.IP = GetIPX()
	}
	n := &NacosCenter{
		ServiceInfo: serv,
		settingMap:  map[string]map[string]*viper.Viper{},
	}
	if defaultLogger != nil {
		n.logger = defaultLogger
	} else {
		n.logger = &DefaultLogger{}
	}
	bs, err := os.ReadFile(KEY_TMP_PORT_PATH)
	if err == nil {
		n.RealPort = cast.ToUint64(string(bs))
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
