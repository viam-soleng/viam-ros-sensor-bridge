package ros_sensor_publisher

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/bluenviron/goroslib/v2"
	"go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	viamutils "go.viam.com/utils"

	"github.com/viam-soleng/viam-ros-sensor-bridge/messages"
	"github.com/viam-soleng/viam-ros-sensor-bridge/utils"
)

var Model = resource.NewModel(utils.Namespace, "ros", "sensor-publisher")

func init() {
	resource.RegisterComponent(
		generic.API,
		Model,
		resource.Registration[resource.Resource, *RosBridgeConfig]{
			Constructor: NewRosSensorPublisher,
		},
	)
}

func NewRosSensorPublisher(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (resource.Resource, error) {
	logger.Infof("Starting Ros Sensor Publisher Module %v", utils.Version)
	c, cancelFunc := context.WithCancel(context.Background())
	b := RosSensorPublisher{
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

type RosSensorPublisher struct {
	resource.Named
	mu         sync.RWMutex
	wg         sync.WaitGroup
	logger     logging.Logger
	cancelFunc context.CancelFunc
	ctx        context.Context
}

// Close implements resource.Resource.
func (r *RosSensorPublisher) Close(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logger.Info("Closing Ros Sensor Publisher Module")
	r.cancelFunc()
	r.wg.Wait()
	return nil
}

// DoCommand implements resource.Resource.
func (*RosSensorPublisher) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{"ok": 1}, nil
}

// Reconfigure implements resource.Resource.
func (r *RosSensorPublisher) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
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

func (r *RosSensorPublisher) reconfigure(newConf *RosBridgeConfig, deps resource.Dependencies) error {
	r.logger.Debug("Stopping existing readers")

	// call the cancel func to stop all readers
	r.cancelFunc()

	// wait for the readers to stop
	r.wg.Wait()
	r.logger.Debug("Readers stopped")

	// Recreate the context since we cancelled it above
	c, cancelFunc := context.WithCancel(context.Background())
	r.cancelFunc = cancelFunc
	r.ctx = c

	for _, s := range newConf.Sensors {
		r.logger.Debugf("Creating sensor %v", s.Name)
		d, err := deps.Lookup(
			resource.Name{
				API:  sensor.API,
				Name: s.Name,
			})
		if err != nil {
			r.logger.Error(err)
			continue
		}

		r.logger.Debugf("Forking reader %v", s.Name)
		reader := RosReader{
			primaryUri:       newConf.PrimaryUri,
			host:             newConf.Host,
			sensorConfig:     s,
			sensor:           d.(sensor.Sensor),
			logger:           r.logger,
			wg:               &r.wg,
			ctx:              r.ctx,
			requestReconnect: make(chan interface{}, 100),
		}
		viamutils.PanicCapturingGo(reader.read())
	}
	return nil
}

type RosReader struct {
	primaryUri       string
	host             string
	sensorConfig     *SensorConfig
	sensor           sensor.Sensor
	logger           logging.Logger
	wg               *sync.WaitGroup
	ctx              context.Context
	p                *goroslib.Publisher
	n                *goroslib.Node
	mu               sync.Mutex
	requestReconnect chan interface{}
}

func (r *RosReader) onLog(level goroslib.LogLevel, msg string) {
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
	} else {
		r.logger.Errorf("Unknown log level: %v, msg: %v", level, msg)
	}

	if strings.Contains(msg, "got an error") && strings.Contains(msg, "dial tcp") {
		r.logger.Warn("Got a dial tcp error, requesting reconnect")
		r.mu.Lock()
		defer r.mu.Unlock()
		r.requestReconnect <- struct{}{}
	}
}

func (r *RosReader) connect() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Debugf("Shutting down existing publisher %v", r.sensorConfig.Topic)
	if r.p != nil {
		r.p.Close()
	}
	r.p = nil
	// Shutdown any existing nodes
	if r.n != nil {
		r.n.Close()
	}

	r.logger.Debugf("Connecting to %v", r.primaryUri)
	node := utils.GetRosNodeWithRetry(r.logger, r.primaryUri, r.host, r.onLog)
	r.n = node

	// Get the type of the message so we can create the publisher later
	messageType, err := messages.GetMessageType(r.sensorConfig.Type)
	if err != nil {
		r.logger.Error(err)
		return
	}
	r.logger.Debugf("Creating publisher %v", r.sensorConfig.Topic)
	publisher, err := goroslib.NewPublisher(goroslib.PublisherConf{
		Node:  node,
		Topic: r.sensorConfig.Topic,
		Msg:   messageType,
	})
	if err == goroslib.ErrNodeTerminated {
		r.logger.Debugf("Node terminated %v", r.sensor.Name().Name)
		r.logger.Debugf("%v", node == nil)
		return
	}
	if err != nil {
		r.logger.Error(err)
		return
	}
	r.p = publisher
}

func (r *RosReader) read() func() {
	return func() {
		r.logger.Infof("Starting reader %v", r.sensor.Name().Name)
		// increment the waitgroup to make sure we wait for this reader to stop
		r.wg.Add(1)
		defer func() {
			// release the waitgroup when this reader stops
			r.wg.Done()
			r.logger.Debugf("Reader fully stopped %v", r.sensor.Name().Name)
		}()

		r.connect()
		// We need to close the publisher when this reader stops
		defer func() {
			r.logger.Debugf("Closing publisher %v", r.sensorConfig.Topic)
			if r.p != nil {
				r.p.Close()
			}
			r.logger.Debugf("Closing node %v", r.sensor.Name().Name)
			if r.n != nil {
				r.n.Close()
			}
		}()
		if r.sensorConfig.SampleRate == 0 {
			r.logger.Warnf("Sample rate is 0, defaulting to 1Hz %v", r.sensor.Name().Name)
			r.sensorConfig.SampleRate = 1
		}
		lastReconnectRequest := time.Now()
		timer := time.NewTimer(time.Duration(1/r.sensorConfig.SampleRate) * time.Second)
		for {
			select {
			case <-r.requestReconnect:
				if time.Since(lastReconnectRequest) < 1*time.Second {
					continue
				}
				r.logger.Infof("Reconnecting %v", r.sensor.Name().Name)
				r.connect()
				lastReconnectRequest = time.Now()
			case <-r.ctx.Done():
				r.logger.Debugf("Reader recevied shutdown signal %v", r.sensor.Name().Name)
				return
			case <-timer.C:
				r.logger.Debugf("Reading sensor %v", r.sensor.Name().Name)
				readings, err := r.sensor.Readings(r.ctx, map[string]interface{}{})
				if err != nil {
					r.logger.Error(err)
					continue
				}

				d, e := messages.ConvertToRosMsg(r.sensorConfig.Type, readings)
				if e != nil {
					r.logger.Error(e)
					continue
				}
				r.logger.Debugf("Publishing message %v", r.sensor.Name().Name)
				// Only try to write if the publisher is there
				if r.p != nil {
					r.p.Write(d)
				} else {
					r.logger.Warnf("Publisher is nil %v, this could mean ROS isn't responding to connection attempts or we are attempting to reconnect", r.sensor.Name().Name)
				}
				timer.Reset(1 * time.Second)
			}
		}
	}
}
