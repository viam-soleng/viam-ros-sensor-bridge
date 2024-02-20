package utils

import (
	"time"

	"github.com/bluenviron/goroslib/v2"
	"go.viam.com/rdk/logging"

	"github.com/viam-soleng/viam-ros-sensor-bridge/viamrosnode"
)

func GetRosNodeWithRetry(logger logging.Logger, primaryUri string, host string, logCallback func(level goroslib.LogLevel, msg string)) *goroslib.Node {
	logger.Debug("Creating new node")

	now := time.Now()
	var node *goroslib.Node
	var err error
	errCount := 0
	for {
		node, err = viamrosnode.GetInsanceWithLogCallback(primaryUri, host, logCallback)
		if node != nil && err == nil {
			break
		}

		// If we fail to create the node, we will retry
		logger.Debugf("Failed to create node: %v", err)
		errCount++
		if errCount%10 == 0 {
			logger.Warnf("Failed to create node after 10 attempts: %v", err)
		}

		time.Sleep(1 * time.Second)
	}

	logger.Debugf("Node created in %v", time.Since(now))
	return node
}
