package configuration

import (
	"bytes"
	"encoding/base64"
	"os"
	"path"

	"github.com/illidaris/aphrodite/pkg/encrypter"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var (
	mode = "debug"
)

type Config vo.NacosClientParam

type ServiceInfo struct {
	GroupName   string       `yaml:"groupname" json:"groupname"`
	ServiceName string       `yaml:"servicename" json:"servicename"`
	ClusterName string       `yaml:"clustername" json:"clustername"`
	IP          string       `yaml:"ip" json:"ip"`
	Port        uint64       `yaml:"port" json:"port"`
	RPort       uint64       `yaml:"rport" json:"rport"`
	Weight      float64      `yaml:"weight" json:"weight"`
	Ephemeral   bool         `yaml:"ephemeral" json:"ephemeral"`
	Enable      bool         `yaml:"enable" json:"enable"`
	Healthy     bool         `yaml:"healthy" json:"healthy"`
	Others      []ConfigName `yaml:"others" json:"others"`
	Sentinel    ConfigName   `yaml:"sentinel" json:"sentinel"`
}

type ConfigName struct {
	ExtConfigFile string `yaml:"extconfigfile" json:"extconfigfile"`
	GroupName     string `yaml:"groupname" json:"groupname"`
	ServiceName   string `yaml:"servicename" json:"servicename"`
}

type SimpleConfig struct {
	Service       ServiceInfo          `yaml:"service" json:"service"`
	ClientConfig  *SimpleClientConfig  `yaml:"clientconfig" json:"clientconfig"`   // optional
	ServerConfigs []SimpleServerConfig `yaml:"serverconfigs" json:"serverconfigs"` // optional
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
	IpAddr   string `yaml:"ipaddr" json:"ipaddr"`     // the nacos server address
	Port     uint64 `yaml:"port" json:"port"`         // nacos server port
	GrpcPort uint64 `yaml:"grpcport" json:"grpcport"` // nacos server grpc port, default=server port + 1000, this is not required
}

func (c *SimpleServerConfig) ToNacosServerConfig() constant.ServerConfig {
	return constant.ServerConfig{
		IpAddr:   c.IpAddr,
		Port:     c.Port,
		GrpcPort: c.GrpcPort,
	}
}

type SimpleClientConfig struct {
	NamespaceId string `yaml:"namespaceid" json:"namespaceid"` // the namespaceId of Nacos.When namespace is public, fill in the blank string here.
	AppName     string `yaml:"appname" json:"appname"`         // the appName
	Endpoint    string `yaml:"endpoint" json:"endpoint"`       // the endpoint for get Nacos server addresses
	Username    string `yaml:"username" json:"username"`       // the username for nacos auth
	Password    string `yaml:"password" json:"password"`       // the password for nacos auth
}

func (c *SimpleClientConfig) GetRealPwd() string {
	if secretKey == "" {
		return c.Password
	}
	pwdBs := []byte{}
	rawBs, err := base64.StdEncoding.DecodeString(c.Password)
	if err != nil {
		return c.Password
	}
	err = encrypter.DecryptStream(bytes.NewBuffer(rawBs), bytes.NewBuffer(pwdBs), []byte(secretKey))
	if err != nil {
		return c.Password
	}
	return string(pwdBs)
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
		constant.WithPassword(c.GetRealPwd()),
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
