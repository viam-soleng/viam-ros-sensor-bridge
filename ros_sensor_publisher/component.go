package ros_sensor_publisher

import (
	"context"
	"sync"
	"time"

	"github.com/bluenviron/goroslib/v2"
	"go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	viamutils "go.viam.com/utils"

	"github.com/viam-soleng/viam-ros-sensor-bridge/messages"
	"github.com/viam-soleng/viam-ros-sensor-bridge/viamrosnode"
)

var Model = resource.NewModel("viam-soleng", "ros", "sensor-publisher")

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
	logger.Info("Starting Ros Sensor Publisher Module v0.0.1")
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
	node       *goroslib.Node
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
	if r.node != nil {
		r.logger.Info("Closing node")
		r.node.Close()
	}
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

	// Close the node if it exists, this is because we need to recreate it with the new primary uri
	// TODO: We can be more selective here, compare the old and new URI?
	if r.node != nil {
		r.logger.Debug("Closing node")
		viamrosnode.ShutdownNodes()
	}

	var err error
	r.logger.Debug("Creating new node")
	r.node, err = viamrosnode.GetInstance(newConf.PrimaryUri)
	if err != nil {
		return err
	}

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
		viamutils.PanicCapturingGo(reader(s, d.(sensor.Sensor), r.node, r.logger, &r.wg, r.ctx))
	}
	return nil
}

func reader(s *SensorConfig, sensor sensor.Sensor, node *goroslib.Node, logger logging.Logger, wg *sync.WaitGroup, ctx context.Context) func() {
	return func() {
		logger.Infof("Starting reader %v", sensor.Name().Name)

		// increment the waitgroup to make sure we wait for this reader to stop
		wg.Add(1)
		defer func() {
			// release the waitgroup when this reader stops
			wg.Done()
			logger.Debugf("Reader fully stopped %v", sensor.Name().Name)
		}()

		// Get the type of the message so we can create the publisher later
		messageType, err := messages.GetMessageType(s.Type)
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Debugf("Creating publisher %v", s.Topic)
		publisher, err := goroslib.NewPublisher(goroslib.PublisherConf{
			Node:  node,
			Topic: s.Topic,
			Msg:   messageType,
		})
		if err == goroslib.ErrNodeTerminated {
			logger.Debugf("Node terminated %v", sensor.Name().Name)
			logger.Debugf("%v", node == nil)
			return
		}
		if err != nil {
			logger.Error(err)
			return
		}
		// We need to close the publisher when this reader stops
		defer func() {
			logger.Debugf("Closing publisher %v", s.Topic)
			publisher.Close()
		}()

		timer := time.NewTimer(time.Duration(1/s.SampleRate) * time.Second)
		for {
			select {
			case <-ctx.Done():
				logger.Debugf("Reader recevied shutdown signal %v", sensor.Name().Name)
				return
			case <-timer.C:
				logger.Debugf("Reading sensor %v", sensor.Name().Name)
				readings, err := sensor.Readings(ctx, map[string]interface{}{})
				if err != nil {
					logger.Error(err)
					continue
				}

				d, e := messages.ConvertToRosMsg(s.Type, readings)
				if e != nil {
					logger.Error(e)
					continue
				}
				logger.Debugf("Publishing message %v", sensor.Name().Name)
				publisher.Write(d)
				timer.Reset(1 * time.Second)
			}
		}
	}
}
