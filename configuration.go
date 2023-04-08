package configuration

type IConfigurationCenter interface {
	SetRealPort(port int)
	GetPort() uint64
	RegisterMine() error
	DeRegisterMine() error
	AddConfigListener(id, group string) error
	Get(key string) interface{}
	GetByID(group, id, key string) interface{}
}
