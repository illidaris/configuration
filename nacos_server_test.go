package configuration

import (
	"testing"

	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetServive(t *testing.T) {
	convey.Convey("TestGetServive", t, func() {
		convey.Convey("success", func() {
			center, err := NewSimpleNacos(&SimpleConfig{
				Service: ServiceInfo{},
				ClientConfig: &SimpleClientConfig{
					NamespaceId: "f9341555-a32d-492b-a604-611fb379e588",
					AppName:     "common.op.demo.microserver",
					Username:    "dev_common_op_demo",
					Password:    "123456",
				},
				ServerConfigs: []SimpleServerConfig{
					{
						IpAddr:   "192.168.97.71",
						Port:     8848,
						GrpcPort: 9848,
					},
				},
			})
			convey.So(err, convey.ShouldBeNil)
			n := center.(*NacosCenter)
			s, err := n.GetService(vo.GetServiceParam{ServiceName: "common.op.demo.microserver", GroupName: "OP"})
			convey.So(err, convey.ShouldBeNil)
			println(s.Name)
			s1, err := n.GetAllService(vo.GetAllServiceInfoParam{NameSpace: "f9341555-a32d-492b-a604-611fb379e588", GroupName: "OP"})
			convey.So(err, convey.ShouldBeNil)
			println(s1.Count)
			s2, err := n.SelectAllInstances(vo.SelectAllInstancesParam{ServiceName: "common.op.demo.microserver", GroupName: "OP"})
			convey.So(err, convey.ShouldBeNil)
			println(s2[0].Ip)
			s3, err := n.SelectInstances(vo.SelectInstancesParam{ServiceName: "common.op.demo.microserver", GroupName: "OP"})
			convey.So(err, convey.ShouldBeNil)
			println(s3[0].Enable)
			s4, err := n.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{ServiceName: "common.op.demo.microserver", GroupName: "OP"})
			convey.So(err, convey.ShouldBeNil)
			println(s4.Ip)
			n.Subscribe(&vo.SubscribeParam{ServiceName: "common.op.demo.microserver", GroupName: "OP", SubscribeCallback: func(services []model.Instance, err error) {
				println(services[0].Ip)
			}})
			convey.So(err, convey.ShouldBeNil)
			// <-time.After(time.Hour * 1)
		})
	})
}
