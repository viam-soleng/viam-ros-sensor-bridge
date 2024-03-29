package main

import (
	"context"

	"go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/utils"

	"github.com/viam-soleng/viam-ros-sensor-bridge/ros_sensor_publisher"
	"github.com/viam-soleng/viam-ros-sensor-bridge/ros_sensor_subscriber"
	module_utils "github.com/viam-soleng/viam-ros-sensor-bridge/utils"
)

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs(module_utils.LoggerName))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) (err error) {
	custom_module, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}

	err = custom_module.AddModelFromRegistry(ctx, generic.API, ros_sensor_publisher.Model)
	if err != nil {
		return err
	}

	err = custom_module.AddModelFromRegistry(ctx, sensor.API, ros_sensor_subscriber.Model)
	if err != nil {
		return err
	}

	err = custom_module.Start(ctx)
	defer custom_module.Close(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
