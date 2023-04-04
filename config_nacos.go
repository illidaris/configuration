package configuration

import (
	"os"
	"path"

	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var (
	mode = "debug"
)

type Config vo.NacosClientParam

type ServiceInfo struct {
	GroupName   string
	ServiceName string
	ClusterName string
	IP          string
	Port        uint64
	Weight      float64
	Ephemeral   bool
	Enable      bool
	Healthy     bool
}

type SimpleConfig struct {
	Service       ServiceInfo
	ClientConfig  *SimpleClientConfig  // optional
	ServerConfigs []SimpleServerConfig // optional
}

func (c *SimpleConfig) ToParam() (vo.NacosClientParam, error) {
	param := vo.NacosClientParam{}
	sc := []constant.ServerConfig{}
	for _, cfg := range c.ServerConfigs {
		sc = append(sc, cfg.ToNacosServerConfig())
	}
	param.ServerConfigs = sc
	opts, err := c.ClientConfig.ToOption()
	if err != nil {
		return param, err
	}
	cc := constant.NewClientConfig(opts...)
	param.ClientConfig = cc
	return param, nil
}

type SimpleServerConfig struct {
	IpAddr   string // the nacos server address
	Port     uint64 // nacos server port
	GrpcPort uint64 // nacos server grpc port, default=server port + 1000, this is not required
}

func (c *SimpleServerConfig) ToNacosServerConfig() constant.ServerConfig {
	return constant.ServerConfig{
		IpAddr:   c.IpAddr,
		Port:     c.Port,
		GrpcPort: c.GrpcPort,
	}
}

type SimpleClientConfig struct {
	NamespaceId string // the namespaceId of Nacos.When namespace is public, fill in the blank string here.
	AppName     string // the appName
	Endpoint    string // the endpoint for get Nacos server addresses
	Username    string // the username for nacos auth
	Password    string // the password for nacos auth
}

func (c *SimpleClientConfig) ToOption() ([]constant.ClientOption, error) {
	opts := []constant.ClientOption{}
	pwd, err := os.Getwd()
	if err != nil {
		return opts, err
	}
	cacheDir := path.Join(pwd, "tmp", "nacos", "cache")
	logDir := path.Join(pwd, "tmp", "nacos", "logs")
	level := "info"
	if mode == "debug" {
		level = mode
	}
	opts = append(opts,
		constant.WithNamespaceId(c.NamespaceId),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithCacheDir(cacheDir),
		constant.WithUsername(c.Username),
		constant.WithPassword(c.Password),
		constant.WithEndpoint(c.Endpoint),
		constant.WithLogDir(logDir),
		constant.WithLogLevel(level),
		constant.WithLogRollingConfig(&constant.ClientLogRollingConfig{
			MaxSize:   100,
			MaxAge:    3,
			LocalTime: true,
		}))
	return opts, nil
}
