package ros_sensor_consumer

import (
	"context"
	"sync"

	"github.com/bluenviron/goroslib/v2"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"

	"github.com/viam-soleng/viam-ros-sensor-bridge/messages"
	"github.com/viam-soleng/viam-ros-sensor-bridge/viamrosnode"
)

var Model = resource.NewModel("viam-soleng", "ros", "sensor-consumer")

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
	logger.Info("Starting Ros Sensor Consumer Module v0.0.1")
	c, cancelFunc := context.WithCancel(context.Background())
	b := RosSensorConsumer{
		Named:      conf.ResourceName().AsNamed(),
		logger:     logger,
		cancelFunc: cancelFunc,
		ctx:        c,
	}

	if err := b.Reconfigure(ctx, deps, conf); err != nil {
		return nil, err
	}
	return &b, nil
}

type RosSensorConsumer struct {
	resource.Named
	mu          sync.RWMutex
	logger      logging.Logger
	node        *goroslib.Node
	cancelFunc  context.CancelFunc
	ctx         context.Context
	subscriber  *goroslib.Subscriber
	lastMessage map[string]interface{}
}

// Readings implements resource.Sensor.
func (r *RosSensorConsumer) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastMessage, nil
}

// Close implements resource.Resource.
func (r *RosSensorConsumer) Close(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logger.Info("Closing Ros Sensor Consumer Module")
	r.cleanup()
	return nil
}

// DoCommand implements resource.Resource.
func (*RosSensorConsumer) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{"ok": 1}, nil
}

// Reconfigure implements resource.Resource.
func (r *RosSensorConsumer) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logger.Debug("Reconfiguring Docker Manager Module")

	newConf, err := resource.NativeConfig[*RosBridgeConfig](conf)
	if err != nil {
		return err
	}

	// In case the module has changed name
	r.Named = conf.ResourceName().AsNamed()

	return r.reconfigure(newConf, deps)
}

func (r *RosSensorConsumer) reconfigure(newConf *RosBridgeConfig, deps resource.Dependencies) error {

	r.cleanup()

	var err error
	r.logger.Debug("Creating new node")
	r.node, err = viamrosnode.GetInstance(newConf.PrimaryUri)
	if err != nil {
		return err
	}

	r.logger.Debugf("Creating sensor %v", newConf.Sensor.Topic)
	handler := messages.NewMessageHandler(r.logger, r.setLastMessage)
	conf, err := handler.GetSubscriberConfigWithHandler(newConf.Sensor.Type)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	conf.Node = r.node
	conf.Topic = newConf.Sensor.Topic

	subscriber, err := goroslib.NewSubscriber(*conf)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	r.subscriber = subscriber

	return nil
}

func (r *RosSensorConsumer) cleanup() {
	r.logger.Debug("Stopping existing consumers")
	if r.subscriber != nil {
		r.subscriber.Close()
	}
	r.logger.Debug("Readers stopped")

	// Close the node if it exists, this is because we need to recreate it with the new primary uri
	// TODO: We can be more selective here, compare the old and new URI?
	if r.node != nil {
		r.logger.Debug("Closing node")
		viamrosnode.ShutdownNodes()
	}
}

func (r *RosSensorConsumer) setLastMessage(m map[string]interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logger.Debugf("Setting last message %v", m)
	r.lastMessage = m
}
