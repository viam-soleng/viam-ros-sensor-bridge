package ros_sensor_subscriber

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/bluenviron/goroslib/v2"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	viamutils "go.viam.com/utils"

	"github.com/viam-soleng/viam-ros-sensor-bridge/messages"
	"github.com/viam-soleng/viam-ros-sensor-bridge/utils"
)

var Model = resource.NewModel(utils.Namespace, "ros", "sensor-subscriber")

func init() {
	resource.RegisterComponent(
		sensor.API,
		Model,
		resource.Registration[sensor.Sensor, *RosBridgeConfig]{
			Constructor: NewRosSensorConsumer,
		},
	)
}

func NewRosSensorConsumer(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (sensor.Sensor, error) {
	logger.Infof("Starting Ros Sensor Consumer Module %v", utils.Version)
	c, cancelFunc := context.WithCancel(context.Background())
	b := RosSensorSubscriber{
		Named:            conf.ResourceName().AsNamed(),
		logger:           logger,
		cancelFunc:       cancelFunc,
		ctx:              c,
		requestReconnect: make(chan bool, 100),
	}

	if err := b.Reconfigure(ctx, deps, conf); err != nil {
		return nil, err
	}
	return &b, nil
}

type RosSensorSubscriber struct {
	resource.Named
	mu               sync.RWMutex
	logger           logging.Logger
	node             *goroslib.Node
	cancelFunc       context.CancelFunc
	ctx              context.Context
	subscriber       *goroslib.Subscriber
	lastMessage      map[string]interface{}
	conf             *RosBridgeConfig
	requestReconnect chan bool
}

// Readings implements resource.Sensor.
func (r *RosSensorSubscriber) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.lastMessage == nil {
		return map[string]interface{}{}, nil
	}
	return r.lastMessage, nil
}

// Close implements resource.Resource.
func (r *RosSensorSubscriber) Close(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cancelFunc()
	r.logger.Info("Closing ROS Sensor Subscriber")
	r.cleanup()
	return nil
}

// DoCommand implements resource.Resource.
func (*RosSensorSubscriber) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{"ok": 1}, nil
}

// Reconfigure implements resource.Resource.
func (r *RosSensorSubscriber) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logger.Info("Reconfiguring ROS Sensor Subscriber")

	newConf, err := resource.NativeConfig[*RosBridgeConfig](conf)
	if err != nil {
		return err
	}

	// In case the module has changed name
	r.Named = conf.ResourceName().AsNamed()
	r.conf = newConf
	viamutils.PanicCapturingGo(r.reconnectHandler())
	r.requestReconnect <- true
	r.logger.Info("Reconfigured ROS Sensor Subscriber")
	return nil
}

func (r *RosSensorSubscriber) reconnectHandler() func() {
	return func() {
		lastReconnectRequest := time.Now()
		for {
			select {
			case <-r.ctx.Done():
				return
			case val := <-r.requestReconnect:
				// we want to be able to force a reconfigure
				if time.Since(lastReconnectRequest) < 1*time.Second && !val {
					continue
				}
				r.logger.Info("Reconnecting to ROS Sensor Subscriber")
				r.connect()
				lastReconnectRequest = time.Now()
			}
		}
	}
}

func (r *RosSensorSubscriber) onLog(level goroslib.LogLevel, msg string) {
	if level == goroslib.LogLevelFatal {
		r.logger.Fatal(msg)
	} else if level == goroslib.LogLevelError {
		r.logger.Error(msg)
	} else if level == goroslib.LogLevelWarn {
		r.logger.Warn(msg)
	} else if level == goroslib.LogLevelInfo {
		r.logger.Info(msg)
	} else if level == goroslib.LogLevelDebug {
		r.logger.Debug(msg)
	}

	if strings.Contains(msg, "got an error") && strings.Contains(msg, "tcp") {
		r.logger.Warn("Got a tcp error, reconnecting")
		r.requestReconnect <- false
	}
}

func (r *RosSensorSubscriber) connect() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cleanup()
	node := utils.GetRosNodeWithRetry(r.logger, r.conf.PrimaryUri, r.conf.Host, r.onLog)

	r.node = node

	r.logger.Infof("Creating ROS Subscriber %v", r.conf.Sensor.Topic)
	handler := messages.NewMessageHandler(r.logger, r.setLastMessage)
	conf, err := handler.GetSubscriberConfigWithHandler(r.conf.Sensor.Type)
	if err != nil {
		r.logger.Errorf("Failed to get subscriber config: %v", err)
		return
	}
	conf.Node = r.node
	conf.Topic = r.conf.Sensor.Topic

	subscriber, err := goroslib.NewSubscriber(*conf)
	if err != nil {
		r.logger.Errorf("Failed to create subscriber: %v", err)
		return
	}
	r.subscriber = subscriber
	r.logger.Infof("Created ROS Subscriber %v", r.conf.Sensor.Topic)
}

func (r *RosSensorSubscriber) cleanup() {
	r.logger.Debug("Stopping existing consumers")
	if r.subscriber != nil {
		r.subscriber.Close()
	}
	r.logger.Debug("Readers stopped")

	if r.node != nil {
		r.logger.Debug("Closing node")
		r.node.Close()
	}
}

func (r *RosSensorSubscriber) setLastMessage(m map[string]interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	m["Timestamp"] = time.Now().UTC().UnixMilli()
	r.logger.Debugf("Setting last message %v", m)
	r.lastMessage = m
}
