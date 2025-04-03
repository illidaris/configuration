package configuration

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/illidaris/aphrodite/pkg/backup"
	"github.com/illidaris/aphrodite/pkg/netex"
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

	Meta         *SimpleConfig
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

func (n *NacosCenter) GetRpcServiceName() string {
	return fmt.Sprintf("%s_%s", n.ServiceName, "rpc")
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

func (n *NacosCenter) GetRpcPort() uint64 {
	if n.RPort > 0 {
		return n.RPort
	}
	portStr := backup.ReadFrmDisk(KEY_TMP_RPORT_PATH)
	port := cast.ToInt(portStr)
	if port > 0 {
		return uint64(port)
	}
	p, _ := netex.GetFreePort()
	if p > 0 {
		backup.WriteToDisk(KEY_TMP_RPORT_PATH, cast.ToString(p))
		return uint64(p)
	}
	return uint64(p)
}

func (n *NacosCenter) GetIP() string {
	return n.IP
}

func (n *NacosCenter) GetMeta() SimpleConfig {
	return *n.Meta
}

func NewNacos(param vo.NacosClientParam, meta *SimpleConfig) (IConfigurationCenter, error) {
	serv := meta.Service
	if serv.IP == "" {
		serv.IP = GetIPX()
	}
	n := &NacosCenter{
		Meta:        meta,
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
	return NewNacos(p, cfg)
}
