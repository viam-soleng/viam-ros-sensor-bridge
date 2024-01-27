package ros_sensor_publisher

import "errors"

type RosBridgeConfig struct {
	PrimaryUri string          `json:"primary_uri"`
	Sensors    []*SensorConfig `json:"sensors"`
}

type SensorConfig struct {
	Topic      string  `json:"topic"`
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	SampleRate float64 `json:"sample_rate"`
}

func (cfg *RosBridgeConfig) Validate(path string) ([]string, error) {
	// NodeName will get default value if string is empty
	if cfg.PrimaryUri == "" {
		return nil, errors.New("primary_uri is required")
	}

	if cfg.Sensors == nil {
		return nil, errors.New("sensors is required")
	}

	for _, sensor := range cfg.Sensors {
		if sensor.Topic == "" {
			return nil, errors.New("topic is required")
		}
		if sensor.Name == "" {
			return nil, errors.New("sensor name is required")
		}
		if sensor.Type == "" {
			return nil, errors.New("sensor type is required")
		}
	}

	return nil, nil
}
