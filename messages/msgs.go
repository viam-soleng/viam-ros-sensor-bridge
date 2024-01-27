package messages

import (
	"encoding/json"
	"errors"

	"github.com/bluenviron/goroslib/v2/pkg/msg"
	"github.com/bluenviron/goroslib/v2/pkg/msgs/sensor_msgs"
	"github.com/bluenviron/goroslib/v2/pkg/msgs/std_msgs"
)

var ErrTypeNotFound = errors.New("type not found")
var ErrFieldNotFound = errors.New("field not found")
var ErrNotFloat64 = errors.New("field not float64")
var ErrFieldNotBool = errors.New("field not bool")

func parseFloat64WithDefault(data map[string]interface{}, field string, defaultValue float64) (float64, error) {
	variance, err := parseFloat64(data, field)
	if err != nil {
		return defaultValue, nil
	}
	return variance, nil
}

func parseFloat64(data map[string]interface{}, field string) (float64, error) {
	raw, ok := data[field]
	if !ok {
		return 0, ErrFieldNotFound
	}
	val, ok := raw.(float64)
	if !ok {
		return 0, ErrNotFloat64
	}
	return val, nil
}

func parseBool(data map[string]interface{}, field string) (bool, error) {
	raw, ok := data[field]
	if !ok {
		return false, ErrFieldNotFound
	}
	val, ok := raw.(bool)
	if !ok {
		return false, ErrFieldNotBool
	}
	return val, nil
}

func GetMessageType(typeName string) (interface{}, error) {
	switch typeName {
	case "":
		return &std_msgs.String{}, nil
	case "sensor_msgs/Temperature":
		return &sensor_msgs.Temperature{}, nil
	case "sensor_msgs/FluidPressure":
		return &sensor_msgs.FluidPressure{}, nil
	case "ThrottlingStates":
		return &ThrottlingStates{}, nil
	default:
		return nil, ErrTypeNotFound
	}
}

func ConvertMessage(sensorName string, typeName string, data map[string]interface{}) (interface{}, error) {
	switch typeName {
	case "ThrottlingStates":
		return convertThrottlingStates(data)
	case "sensor_msgs/Temperature":
		return convertTemperature(data)
	case "sensor_msgs/FluidPressure":
		return convertFluidPressure(data)
	case "std_msgs/String":
		var d string
		for _, value := range data {
			d = value.(string)
		}
		return &std_msgs.String{Data: d}, nil
	case "json":
		j, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		return &std_msgs.String{Data: string(j)}, nil
	default:
		return nil, ErrTypeNotFound
	}
}

func convertTemperature(data map[string]interface{}) (*sensor_msgs.Temperature, error) {
	o := &sensor_msgs.Temperature{}

	var err error
	if o.Temperature, err = parseFloat64(data, "temperature"); err != nil {
		return nil, err
	}

	if o.Variance, err = parseFloat64WithDefault(data, "variance", 0); err != nil {
		return nil, err
	}

	return o, nil
}

func convertFluidPressure(data map[string]interface{}) (*sensor_msgs.FluidPressure, error) {
	o := &sensor_msgs.FluidPressure{}

	var err error
	if o.FluidPressure, err = parseFloat64(data, "pressure"); err != nil {
		return nil, err
	}

	if o.Variance, err = parseFloat64WithDefault(data, "variance", 0); err != nil {
		return nil, err
	}

	return o, nil
}

type ThrottlingStates struct {
	msg.Package                  `ros:"sample_msgs"`
	Header                       std_msgs.Header `rosname:"header"`
	Undervoltage                 bool            `rosname:"undervoltage"`
	ArmFrequentlyCapped          bool            `rosname:"arm_frequently_capped"`
	Throttled                    bool            `rosname:"throttled"`
	SoftTemperatureLimitActive   bool            `rosname:"soft_temperature_limit_active"`
	UndervoltageOccurred         bool            `rosname:"undervoltage_occurred"`
	ArmFrequentlyCappedOccurred  bool            `rosname:"arm_frequently_capped_occurred"`
	ThrottlingOccurred           bool            `rosname:"throttling_occurred"`
	SoftTemperatureLimitOccurred bool            `rosname:"soft_temperature_limit_occurred"`
}

func convertThrottlingStates(data map[string]interface{}) (*ThrottlingStates, error) {
	o := &ThrottlingStates{}
	var err error
	if o.Undervoltage, err = parseBool(data, "undervolt"); err != nil {
		return nil, err
	}

	if o.ArmFrequentlyCapped, err = parseBool(data, "arm_frequency_capped"); err != nil {
		return nil, err
	}

	if o.Throttled, err = parseBool(data, "currently_throttled"); err != nil {
		return nil, err
	}

	if o.SoftTemperatureLimitActive, err = parseBool(data, "soft_temp_limit_active"); err != nil {
		return nil, err
	}

	if o.UndervoltageOccurred, err = parseBool(data, "under_volt_occurred"); err != nil {
		return nil, err
	}

	if o.ArmFrequentlyCappedOccurred, err = parseBool(data, "arm_frequency_cap_occured"); err != nil {
		return nil, err
	}

	if o.ThrottlingOccurred, err = parseBool(data, "throttling_occurred"); err != nil {
		return nil, err
	}

	if o.SoftTemperatureLimitOccurred, err = parseBool(data, "soft_temp_limit_occurred"); err != nil {
		return nil, err
	}

	return o, nil
}
