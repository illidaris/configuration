package configuration

import (
	"testing"

	_ "github.com/smartystreets/goconvey/convey"
)

func TestXxx(t *testing.T) {
	NewSimpleNacos(&SimpleConfig{
		Service:       ServiceInfo{},
		ClientConfig:  &SimpleClientConfig{},
		ServerConfigs: []SimpleServerConfig{},
	})

}
