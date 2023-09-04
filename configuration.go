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
	RegisterMine(meta map[string]string) error
	DeRegisterMine() error
}

type IDiscoverCenter interface {
	DiscoverInstanceOne(group, service string, clusters ...string) (string, error)
}

type IConfigCenter interface {
	AddConfigListener(id, group string, callback func(string, string, string, string)) error
	Get(key string) interface{}
	GetByID(group, id, key string) interface{}
}
