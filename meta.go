package configuration

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path"

	pathEx "github.com/illidaris/file/path"
	"github.com/spf13/viper"
)

const (
	KEY_IS_REMOTE = "isRemote"
)

var DefaultCenter IConfigurationCenter

func IsRemote() bool {
	b := viper.GetBool(KEY_IS_REMOTE)
	return b
}

func Get(key string) interface{} {
	if IsRemote() {
		return DefaultCenter.Get(key)
	}
	return viper.Get(key)
}
func LoadConfig(configPath string) error {
	b, err := pathEx.ExistOrNot(configPath)
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("[config]%s has no find", configPath)
	}
	viper.SetConfigFile(configPath)
	// If a config file is found, read it in.
	if readErr := viper.ReadInConfig(); readErr != nil {
		return readErr
	}
	// enable remote config
	if IsRemote() {
		nacosConfig := &SimpleConfig{}
		if err := viper.UnmarshalKey("nacos", nacosConfig); err != nil {
			return err
		}
		DefaultCenter, err = NewSimpleNacos(nacosConfig)
		if err != nil {
			return err
		}
		if err := DefaultCenter.AddConfigListener(nacosConfig.Service.ServiceName, nacosConfig.Service.GroupName); err != nil {
			return err
		}
	}
	viper.WatchConfig()
	return nil
}

func Init() error {
	p, _ := os.Getwd()
	baseP := path.Join(p, "config", "config.yml")
	err := LoadConfig(baseP)
	if err != nil {
		println(err)
	}
	return err
}

func GetIPX() string {
	ip, _ := GetIP()
	return ip
}

func GetIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("ip no found")
}