package ros_sensor_subscriber

import "errors"

type RosBridgeConfig struct {
	PrimaryUri string        `json:"primary_uri"`
	Host       string        `json:"host"`
	Sensor     *SensorConfig `json:"sensor"`
}

type SensorConfig struct {
	Topic string `json:"topic"`
	Type  string `json:"message_type"`
}

func (cfg *RosBridgeConfig) Validate(path string) ([]string, error) {
	// NodeName will get default value if string is empty
	if cfg.PrimaryUri == "" {
		return nil, errors.New("primary_uri is required")
	}
	if cfg.Sensor == nil {
		return nil, errors.New("sensor is required")
	} else {
		if cfg.Sensor.Topic == "" {
			return nil, errors.New("topic is required")
		}
		if cfg.Sensor.Type == "" {
			return nil, errors.New("sensor type is required")
		}
	}
	return nil, nil
}
