package configuration

type IConfigurationCenter interface {
	IConfigCenter
	IRegisterCenter
	IDiscoverCenter
	SetILogger(log ILogger)
}

type IRegisterCenter interface {
	SetRealPort(port int) error
	GetPort() uint64
	GetServiceName() string
	RegisterMine(meta map[string]string) error
	RegisterServ(meta map[string]string, srvname string, port uint64) error
	DeRegisterMine() error
	DeRegisterServ(srvname string, port uint64) error
}

type IDiscoverCenter interface {
	DiscoverInstanceOne(group, service string, clusters ...string) (string, error)
}

type IConfigCenter interface {
	AddConfigListener(id, group string, callback func(string, string, string, string)) error
	Get(key string) interface{}
	GetByID(group, id, key string) interface{}
}
