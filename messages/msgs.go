package messages

import (
	"encoding/json"
	"errors"

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

func GetMessageType(typeName string) (interface{}, error) {
	for _, registry := range type_registries {
		creator, ok := registry[typeName]
		if ok {
			return creator(), nil
		}
	}

	return nil, ErrTypeNotFound
}

type HandlerRegistry map[string]func(*MessageHandler) *goroslib.SubscriberConf

var std_msgs_handler_registry = HandlerRegistry{
	"std_msgs/Header": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.Header) { h.handleMessage(msg) }}
	},
	"std_msgs/String": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.String) { h.handleMessage(msg) }}
	},
	"std_msgs/Bool": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.Bool) { h.handleMessage(msg) }}
	},
	"std_msgs/Int8": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.Int8) { h.handleMessage(msg) }}
	},
	"std_msgs/Int16": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.Int16) { h.handleMessage(msg) }}
	},
	"std_msgs/Int32": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.Int32) { h.handleMessage(msg) }}
	},
	"std_msgs/Int64": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.Int64) { h.handleMessage(msg) }}
	},
	"std_msgs/UInt8": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.UInt8) { h.handleMessage(msg) }}
	},
	"std_msgs/UInt16": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.UInt16) { h.handleMessage(msg) }}
	},
	"std_msgs/UInt32": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.UInt32) { h.handleMessage(msg) }}
	},
	"std_msgs/UInt64": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.UInt64) { h.handleMessage(msg) }}
	},
	"std_msgs/Float32": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.Float32) { h.handleMessage(msg) }}
	},
	"std_msgs/Float64": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.Float64) { h.handleMessage(msg) }}
	},
	"std_msgs/Time": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.Time) { h.handleMessage(msg) }}
	},
	"std_msgs/Duration": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.Duration) { h.handleMessage(msg) }}
	},
	"std_msgs/ColorRGBA": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.ColorRGBA) { h.handleMessage(msg) }}
	},
	"std_msgs/MultiArrayDimension": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.MultiArrayDimension) { h.handleMessage(msg) }}
	},
	"std_msgs/MultiArrayLayout": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.MultiArrayLayout) { h.handleMessage(msg) }}
	},
	"std_msgs/Byte": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.Byte) { h.handleMessage(msg) }}
	},
	"std_msgs/ByteMultiArray": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.ByteMultiArray) { h.handleMessage(msg) }}
	},
	"std_msgs/Char": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.Char) { h.handleMessage(msg) }}
	},
	"std_msgs/Empty": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *std_msgs.Empty) { h.handleMessage(msg) }}
	},
}

var custom_handler_registry = HandlerRegistry{
	"ThrottlingStates": func(h *MessageHandler) *goroslib.SubscriberConf {
		return &goroslib.SubscriberConf{Callback: func(msg *ThrottlingStates) { h.handleMessage(msg) }}
	},
}

var handler_registries = []HandlerRegistry{std_msgs_handler_registry, custom_handler_registry}

type MessageHandler struct {
	logger         logging.Logger
	setLastMessage func(map[string]interface{})
}

func NewMessageHandler(logger logging.Logger, setLastMessage func(map[string]interface{})) *MessageHandler {
	return &MessageHandler{logger: logger, setLastMessage: setLastMessage}
}

func (h *MessageHandler) handleMessage(msg interface{}) error {
	h.logger.Debugf("Converting message %#v", msg)
	m, err := ConvertFromRosMsg(msg)
	if err != nil {
		h.logger.Error(err)
		return err
	}
	h.logger.Debugf("Converted message %#v", m)
	h.setLastMessage(m)
	return nil
}

func (h *MessageHandler) GetSubscriberConfigWithHandler(typeName string) (*goroslib.SubscriberConf, error) {
	for _, registry := range handler_registries {
		creator, ok := registry[typeName]
		if ok {
			return creator(h), nil
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
