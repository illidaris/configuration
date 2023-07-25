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
	KEY_IS_REMOTE     = "isRemote"
	KEY_TMP_PORT_PATH = "./tmp/port.txt"
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
func LoadConfig(configPath string, callback func(string, string, string, string)) error {
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
		if DefaultCenter == nil {
			return errors.New("DefaultCenter is nil")
		}
		if err := DefaultCenter.AddConfigListener(nacosConfig.Service.ServiceName, nacosConfig.Service.GroupName, callback); err != nil {
			return err
		}
		for _, other := range nacosConfig.Service.Others {
			otherErr := DefaultCenter.AddConfigListener(other.ServiceName, other.GroupName, callback)
			if otherErr != nil {
				println(otherErr)
			}
		}
	}
	viper.WatchConfig()
	return nil
}

func Init(configPath string) error {
	if configPath == "" {
		p, _ := os.Getwd()
		configPath = path.Join(p, "config", "config.yml")
	}
	err := LoadConfig(configPath, nil)
	if err != nil {
		println(err.Error())
	}
	return err
}

func SetILogger(log ILogger) {
	defaultLogger = log
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
