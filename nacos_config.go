package configuration

import (
	"bytes"
	"context"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

// AddConfigListener add config listen from remote nacos center
func (n *NacosCenter) AddConfigListener(id, group string, callback func(string, string, string, string)) error {
	// step01 get config
	data, err := n.GetContent(group, id)
	if err != nil {
		return err
	}
	if err := n.SetConfig("", group, id, data, callback); err != nil {
		return err
	}
	// step01 listen change
	return n.ListenContent(group, id, func(namespace, group, dataId, data string) {
		if err := n.SetConfig(namespace, group, dataId, data, callback); err != nil {
			n.logger.ErrorCtx(context.TODO(), err.Error())
		}
	})
}

// SetConfig set viper config
func (n *NacosCenter) SetConfig(namespace, group, dataId, data string, callback func(string, string, string, string)) error {
	n.logger.InfoCtx(context.TODO(), fmt.Sprintf("[%s|%s|%s]%s", namespace, group, dataId, data))
	if callback != nil {
		callback(namespace, group, dataId, data)
	}
	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader([]byte(data))); err != nil {
		return err
	}
	if _, ok := n.settingMap[group]; !ok {
		n.settingMap[group] = map[string]*viper.Viper{}
	}
	n.settingMap[group][dataId] = v
	return nil
}

// Get get kv
func (n *NacosCenter) Get(key string) interface{} {
	data := n.getMeta(n.GroupName, n.ServiceName)
	if data == nil {
		return nil
	}
	return data.Get(key)
}

// GetByID get kv
func (n *NacosCenter) GetByID(group, id, key string) interface{} {
	data := n.getMeta(group, id)
	if data == nil {
		return nil
	}
	return data.Get(key)
}

// getMeta get kv
func (n *NacosCenter) getMeta(group, id string) *viper.Viper {
	if _, ok := n.settingMap[group]; !ok {
		return nil
	}
	if _, ok := n.settingMap[group][id]; !ok {
		return nil
	}
	return n.settingMap[group][id]
}

func (n *NacosCenter) GetContent(group, id string) (string, error) {
	param := vo.ConfigParam{
		DataId: id,
		Group:  group,
	}
	return n.ConfigClient.GetConfig(param)
}

func (n *NacosCenter) ListenContent(group, id string, f func(namespace, group, dataId, data string)) error {
	param := vo.ConfigParam{
		DataId:   id,
		Group:    group,
		OnChange: f,
	}
	return n.ConfigClient.ListenConfig(param)
}
