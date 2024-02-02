package messages

import (
	"encoding/json"
	"errors"

	"github.com/bluenviron/goroslib/v2/pkg/msg"
	"github.com/bluenviron/goroslib/v2/pkg/msgs/std_msgs"
)

var ErrTypeNotFound = errors.New("type not found")

type TypeRegistry map[string]func() interface{}

var std_msgs_registry = TypeRegistry{
	"std_msgs/Header":              func() interface{} { return &std_msgs.Header{} },
	"std_msgs/String":              func() interface{} { return &std_msgs.String{} },
	"std_msgs/Bool":                func() interface{} { return &std_msgs.Bool{} },
	"std_msgs/Int8":                func() interface{} { return &std_msgs.Int8{} },
	"std_msgs/Int16":               func() interface{} { return &std_msgs.Int16{} },
	"std_msgs/Int32":               func() interface{} { return &std_msgs.Int32{} },
	"std_msgs/Int64":               func() interface{} { return &std_msgs.Int64{} },
	"std_msgs/UInt8":               func() interface{} { return &std_msgs.UInt8{} },
	"std_msgs/UInt16":              func() interface{} { return &std_msgs.UInt16{} },
	"std_msgs/UInt32":              func() interface{} { return &std_msgs.UInt32{} },
	"std_msgs/UInt64":              func() interface{} { return &std_msgs.UInt64{} },
	"std_msgs/Float32":             func() interface{} { return &std_msgs.Float32{} },
	"std_msgs/Float64":             func() interface{} { return &std_msgs.Float64{} },
	"std_msgs/Time":                func() interface{} { return &std_msgs.Time{} },
	"std_msgs/Duration":            func() interface{} { return &std_msgs.Duration{} },
	"std_msgs/ColorRGBA":           func() interface{} { return &std_msgs.ColorRGBA{} },
	"std_msgs/MultiArrayDimension": func() interface{} { return &std_msgs.MultiArrayDimension{} },
	"std_msgs/MultiArrayLayout":    func() interface{} { return &std_msgs.MultiArrayLayout{} },
	"std_msgs/Byte":                func() interface{} { return &std_msgs.Byte{} },
	"std_msgs/ByteMultiArray":      func() interface{} { return &std_msgs.ByteMultiArray{} },
	"std_msgs/Char":                func() interface{} { return &std_msgs.Char{} },
	"std_msgs/Empty":               func() interface{} { return &std_msgs.Empty{} },
}

var custom_registry = TypeRegistry{
	"ThrottlingStates": func() interface{} { return &ThrottlingStates{} },
}

var registries = []TypeRegistry{std_msgs_registry, custom_registry}

func GetMessageType(typeName string) (interface{}, error) {
	for _, registry := range registries {
		creator, ok := registry[typeName]
		if ok {
			return creator(), nil
		}
	}

	return nil, ErrTypeNotFound
}

func ConvertToRosMsg(typeName string, data map[string]interface{}) (interface{}, error) {
	t, err := GetMessageType(typeName)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}
	return t, err
}

func ConvertFromRosMsg(msg interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	var mapResult map[string]interface{}
	err = json.Unmarshal(data, &mapResult)
	if err != nil {
		return nil, err
	}
	return mapResult, nil
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
