package messages

import (
	"github.com/bluenviron/goroslib/v2/pkg/msg"
	"github.com/bluenviron/goroslib/v2/pkg/msgs/std_msgs"
)

var custom_type_registry = TypeRegistry{
	"ThrottlingStates": func() interface{} { return &ThrottlingStates{} },
	// Add more custom message types here
}

func GetCustomMsgsCallback(handleMessage func(interface{}) error, typeName string) interface{} {
	switch typeName {
	case "ThrottlingStates":
		return func(msg *ThrottlingStates) { handleMessage(msg) }
		// Add more custom message types here
	}
	return nil
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
