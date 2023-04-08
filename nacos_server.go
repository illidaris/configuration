package configuration

import (
	"errors"
	"fmt"

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

func (n *NacosCenter) BatchRegisterService(param vo.BatchRegisterInstanceParam) {
	success, err := n.NamingClient.BatchRegisterInstance(param)
	if !success || err != nil {
		panic("BatchRegisterServiceInstance failed!" + err.Error())
	}
	fmt.Printf("BatchRegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}

// =========================== discover ===========================

func (n *NacosCenter) GetService(param vo.GetServiceParam) {
	service, err := n.NamingClient.GetService(param)
	if err != nil {
		panic("GetService failed!" + err.Error())
	}
	fmt.Printf("GetService,param:%+v, result:%+v \n\n", param, service)
}

func (n *NacosCenter) GetAllService(param vo.GetAllServiceInfoParam) {
	service, err := n.NamingClient.GetAllServicesInfo(param)
	if err != nil {
		panic("GetAllService failed!")
	}
	fmt.Printf("GetAllService,param:%+v, result:%+v \n\n", param, service)
}

func (n *NacosCenter) SelectAllInstances(param vo.SelectAllInstancesParam) {
	instances, err := n.NamingClient.SelectAllInstances(param)
	if err != nil {
		panic("SelectAllInstances failed!" + err.Error())
	}
	fmt.Printf("SelectAllInstance,param:%+v, result:%+v \n\n", param, instances)
}

func (n *NacosCenter) SelectInstances(param vo.SelectInstancesParam) {
	instances, err := n.NamingClient.SelectInstances(param)
	if err != nil {
		panic("SelectInstances failed!" + err.Error())
	}
	fmt.Printf("SelectInstances,param:%+v, result:%+v \n\n", param, instances)
}

func (n *NacosCenter) SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) {
	instances, err := n.NamingClient.SelectOneHealthyInstance(param)
	if err != nil {
		panic("SelectOneHealthyInstance failed!")
	}
	fmt.Printf("SelectOneHealthyInstance,param:%+v, result:%+v \n\n", param, instances)
}

func (n *NacosCenter) Subscribe(param *vo.SubscribeParam) {
	n.NamingClient.Subscribe(param)
}

func (n *NacosCenter) UnSubscribe(param *vo.SubscribeParam) {
	n.NamingClient.Unsubscribe(param)
}
