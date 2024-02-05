package messages

import (
	"reflect"
	"testing"

	"github.com/bluenviron/goroslib/v2/pkg/msgs/std_msgs"
	"github.com/stretchr/testify/assert"
	"go.viam.com/rdk/logging"
)

func TestGetCallback(t *testing.T) {
	logger := logging.NewTestLogger(t)
	setLastMessage := func(msg map[string]interface{}) {
		t.Logf("Received message: %#v", msg)
	}
	handler := NewMessageHandler(logger, setLastMessage)
	f := handler.getCallback("std_msgs/Time")

	// Check the return to ensure the function accepts only one parameter, and that parameter is *std_msgs/Time{}
	assert.Equal(t, reflect.TypeOf(f).NumIn(), 1, "Function should accept only one parameter")
	assert.Equal(t, reflect.TypeOf(f).In(0), reflect.TypeOf(&std_msgs.Time{}), "Parameter should be *std_msgs/Time{}")
}
