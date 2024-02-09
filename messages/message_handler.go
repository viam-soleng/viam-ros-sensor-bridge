package messages

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/bluenviron/goroslib/v2"
	"go.viam.com/rdk/logging"
)

var ErrTypeNotFound = errors.New("type not found")

type TypeRegistry map[string]func() interface{}
type SerializerTypeRegistry map[string]func(map[string]interface{}) interface{}

type MessageHandler struct {
	logger         logging.Logger
	setLastMessage func(map[string]interface{})
}

func NewMessageHandler(logger logging.Logger, setLastMessage func(map[string]interface{})) *MessageHandler {
	return &MessageHandler{logger: logger, setLastMessage: setLastMessage}
}

func (h *MessageHandler) getCallback(typeName string) interface{} {
	handler := h.handleMessage
	if strings.HasPrefix(typeName, "std_msgs/") {
		return GetStdMsgsCallback(handler, typeName)
	} else {
		return GetCustomMsgsCallback(handler, typeName)
	}
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

var type_registries = []TypeRegistry{std_msgs_registry, custom_type_registry}

func GetMessageType(typeName string) (interface{}, error) {
	for _, registry := range type_registries {
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
	b, err := json.Marshal(data)
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
