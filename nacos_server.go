package configuration

import (
	"errors"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// ============================ register ============================

func (n *NacosCenter) RegisterMine() error {
	param := vo.RegisterInstanceParam{
		Ip:          n.IP,
		Port:        n.GetPort(),
		Weight:      n.Weight,
		Enable:      n.Enable,
		Healthy:     n.Healthy,
		ClusterName: n.ClusterName,
		ServiceName: n.ServiceName,
		GroupName:   n.GroupName,
		Ephemeral:   n.Ephemeral,
	}
	success, err := n.RegisterService(param)
	if err != nil {
		return err
	}
	if !success {
		return errors.New("RegisterServiceInstance failed")
	}
	return nil
}

func (n *NacosCenter) RegisterService(param vo.RegisterInstanceParam) (bool, error) {
	return n.NamingClient.RegisterInstance(param)
}

func (n *NacosCenter) DeRegisterMine() error {
	param := vo.DeregisterInstanceParam{
		Ip:          n.IP,
		Port:        n.GetPort(),
		Cluster:     n.ClusterName,
		ServiceName: n.ServiceName,
		GroupName:   n.GroupName,
	}
	success, err := n.DeRegisterService(param)
	if err != nil {
		return err
	}
	if !success {
		return errors.New("DeRegisterServiceInstance failed")
	}
	return nil
}

func (n *NacosCenter) DeRegisterService(param vo.DeregisterInstanceParam) (bool, error) {
	return n.NamingClient.DeregisterInstance(param)
}

func (n *NacosCenter) UpdateMine() error {
	param := vo.UpdateInstanceParam{
		Ip:          n.IP,
		Port:        n.GetPort(),
		Weight:      n.Weight,
		Enable:      n.Enable,
		Healthy:     n.Healthy,
		ClusterName: n.ClusterName,
		ServiceName: n.ServiceName,
		GroupName:   n.GroupName,
		Ephemeral:   n.Ephemeral,
	}
	success, err := n.UpdateService(param)
	if err != nil {
		return err
	}
	if !success {
		return errors.New("UpdateServiceInstance failed")
	}
	return nil
}

func (n *NacosCenter) UpdateService(param vo.UpdateInstanceParam) (bool, error) {
	return n.NamingClient.UpdateInstance(param)
}

// =========================== other ===========================

func (n *NacosCenter) BatchRegisterService(param vo.BatchRegisterInstanceParam) (bool, error) {
	success, err := n.NamingClient.BatchRegisterInstance(param)
	if !success || err != nil {
		panic("BatchRegisterServiceInstance failed!" + err.Error())
	}
	fmt.Printf("BatchRegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
	return success, err
}

// =========================== discover ===========================

func (n *NacosCenter) DiscoverInstanceOne(group, service string, clusters ...string) (string, error) {
	instance, err := n.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		GroupName:   group,
		ServiceName: service,
		Clusters:    clusters,
	})
	if err != nil {
		return "", err
	}
	if instance == nil {
		return "", errors.New("no found instance")
	}
	return fmt.Sprintf("%s:%d", instance.Ip, instance.Port), nil
}

func (n *NacosCenter) GetService(param vo.GetServiceParam) (model.Service, error) {
	return n.NamingClient.GetService(param)
}

func (n *NacosCenter) GetAllService(param vo.GetAllServiceInfoParam) (model.ServiceList, error) {
	return n.NamingClient.GetAllServicesInfo(param)
}

func (n *NacosCenter) SelectAllInstances(param vo.SelectAllInstancesParam) ([]model.Instance, error) {
	return n.NamingClient.SelectAllInstances(param)
}

func (n *NacosCenter) SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error) {
	return n.NamingClient.SelectInstances(param)
}

func (n *NacosCenter) SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (*model.Instance, error) {
	return n.NamingClient.SelectOneHealthyInstance(param)
}

func (n *NacosCenter) Subscribe(param *vo.SubscribeParam) {
	n.NamingClient.Subscribe(param)
}

func (n *NacosCenter) UnSubscribe(param *vo.SubscribeParam) {
	n.NamingClient.Unsubscribe(param)
}
