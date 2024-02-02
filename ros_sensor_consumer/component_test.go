package ros_sensor_consumer

import (
	"os"
	"testing"

	"go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/resource"
)

func component_test_setup(t *testing.T) (resource.Config, resource.Dependencies) {
	os.Setenv("VIAM_MODULE_DATA", os.TempDir())
	cfg := resource.Config{
		Name:                "foo",
		Model:               Model,
		API:                 generic.API,
		ConvertedAttributes: &RosBridgeConfig{},
	}

	return cfg, resource.Dependencies{}
}
