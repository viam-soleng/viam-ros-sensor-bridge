package messages

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/bluenviron/goroslib/v2"
	"github.com/bluenviron/goroslib/v2/pkg/msg"
	"github.com/bluenviron/goroslib/v2/pkg/msgs/std_msgs"
	"go.viam.com/rdk/logging"
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

var type_registries = []TypeRegistry{std_msgs_registry, custom_registry}

type MessageHandler struct {
	logger         logging.Logger
	setLastMessage func(map[string]interface{})
}

func NewMessageHandler(logger logging.Logger, setLastMessage func(map[string]interface{})) *MessageHandler {
	return &MessageHandler{logger: logger, setLastMessage: setLastMessage}
}

func GetMessageType(typeName string) (interface{}, error) {
	for _, registry := range type_registries {
		creator, ok := registry[typeName]
		if ok {
			return creator(), nil
		}
	}

	return nil, ErrTypeNotFound
}

func GetStdMsgsCallback(handleMessage func(interface{}) error, typeName string) interface{} {
	switch typeName {
	case "std_msgs/Header":
		return func(msg *std_msgs.Header) { handleMessage(msg) }
	case "std_msgs/String":
		return func(msg *std_msgs.String) { handleMessage(msg) }
	case "std_msgs/Bool":
		return func(msg *std_msgs.Bool) { handleMessage(msg) }
	case "std_msgs/Int8":
		return func(msg *std_msgs.Int8) { handleMessage(msg) }
	case "std_msgs/Int16":
		return func(msg *std_msgs.Int16) { handleMessage(msg) }
	case "std_msgs/Int32":
		return func(msg *std_msgs.Int32) { handleMessage(msg) }
	case "std_msgs/Int64":
		return func(msg *std_msgs.Int64) { handleMessage(msg) }
	case "std_msgs/UInt8":
		return func(msg *std_msgs.UInt8) { handleMessage(msg) }
	case "std_msgs/UInt16":
		return func(msg *std_msgs.UInt16) { handleMessage(msg) }
	case "std_msgs/UInt32":
		return func(msg *std_msgs.UInt32) { handleMessage(msg) }
	case "std_msgs/UInt64":
		return func(msg *std_msgs.UInt64) { handleMessage(msg) }
	case "std_msgs/Float32":
		return func(msg *std_msgs.Float32) { handleMessage(msg) }
	case "std_msgs/Float64":
		return func(msg *std_msgs.Float64) { handleMessage(msg) }
	case "std_msgs/Time":
		return func(msg *std_msgs.Time) { handleMessage(msg) }
	case "std_msgs/Duration":
		return func(msg *std_msgs.Duration) { handleMessage(msg) }
	case "std_msgs/ColorRGBA":
		return func(msg *std_msgs.ColorRGBA) { handleMessage(msg) }
	case "std_msgs/MultiArrayDimension":
		return func(msg *std_msgs.MultiArrayDimension) { handleMessage(msg) }
	case "std_msgs/MultiArrayLayout":
		return func(msg *std_msgs.MultiArrayLayout) { handleMessage(msg) }
	case "std_msgs/Byte":
		return func(msg *std_msgs.Byte) { handleMessage(msg) }
	case "std_msgs/ByteMultiArray":
		return func(msg *std_msgs.ByteMultiArray) { handleMessage(msg) }
	case "std_msgs/Char":
		return func(msg *std_msgs.Char) { handleMessage(msg) }
	case "std_msgs/Empty":
		return func(msg *std_msgs.Empty) { handleMessage(msg) }
	}
	return nil
}

func GetCustomMsgsCallback(handleMessage func(interface{}), typeName string) interface{} {
	switch typeName {
	case "ThrottlingStates":
		return func(msg *ThrottlingStates) { handleMessage(msg) }
	}
	return nil
}

func (h *MessageHandler) getCallback(typeName string) interface{} {
	handler := h.handleMessage
	if strings.HasPrefix(typeName, "std_msgs/") {
		return GetStdMsgsCallback(handler, typeName)
	}
	return nil
}

func (h *MessageHandler) GetSubscriberConfigWithHandler(typeName string) (*goroslib.SubscriberConf, error) {
	handler := h.getCallback(typeName)
	if handler != nil {
		return &goroslib.SubscriberConf{Callback: handler}, nil
	}
	return nil, ErrTypeNotFound
}

func (h *MessageHandler) handleMessage(msg interface{}) error {
	h.logger.Debugf("Converting message %#v", msg)
	m, err := convertFromRosMsg(msg)
	if err != nil {
		h.logger.Error(err)
		return err
	}
	h.logger.Debugf("Converted message %#v", m)
	h.setLastMessage(m)
	return nil
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

func convertFromRosMsg(msg interface{}) (map[string]interface{}, error) {
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
