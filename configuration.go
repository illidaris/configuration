package configuration

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
)

type IConfigurationCenter interface {
	IConfigCenter
	IRegisterCenter
	IDiscoverCenter
	SetILogger(log ILogger)
	GetServiceInfo() ServiceInfo
}

type IRegisterCenter interface {
	SetRealPort(port int) error
	GetIP() string
	GetPort() uint64
	GetRpcPort() uint64
	GetServiceName() string
	RegisterMine(meta map[string]string) error
	RegisterServ(meta map[string]string, srvname string, port uint64) error
	DeRegisterMine() error
	DeRegisterServ(srvname string, port uint64) error
	GetMeta() SimpleConfig
}

type IDiscoverCenter interface {
	GetNamingClient() naming_client.INamingClient
	DiscoverInstanceOne(group, service string, clusters ...string) (string, error)
	DiscoverInstances(group, service string, clusters ...string) ([]string, error)
}

type IConfigCenter interface {
	AddConfigListener(id, group string, callback func(string, string, string, string)) error
	GetConfigClient() config_client.IConfigClient
	Get(key string) interface{}
	GetByID(group, id, key string) interface{}
}
