package messages

import (
	"reflect"
	"testing"
	"time"

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

func TestConvertStringReadings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/String", map[string]interface{}{"Data": "Hello, World!"})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, "Hello, World!", m.(*std_msgs.String).Data, "Data should be 'Hello, World!'")
}

func TestConvertBoolReadings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/Bool", map[string]interface{}{"Data": true})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, true, m.(*std_msgs.Bool).Data, "Data should be true")
}

func TestConvertInt8Readings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/Int8", map[string]interface{}{"Data": 2})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, int8(2), m.(*std_msgs.Int8).Data, "Data should be 1")
}

func TestConvertInt16Readings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/Int16", map[string]interface{}{"Data": 42})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, int16(42), m.(*std_msgs.Int16).Data, "Data should be 42")
}

func TestConvertInt32Readings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/Int32", map[string]interface{}{"Data": 42})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, int32(42), m.(*std_msgs.Int32).Data, "Data should be 42")
}

func TestConvertInt64Readings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/Int64", map[string]interface{}{"Data": 42})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, int64(42), m.(*std_msgs.Int64).Data, "Data should be 42")
}

func TestConvertUInt8Readings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/UInt8", map[string]interface{}{"Data": 42})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, uint8(42), m.(*std_msgs.UInt8).Data, "Data should be 42")
}

func TestConvertUInt16Readings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/UInt16", map[string]interface{}{"Data": 42})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, uint16(42), m.(*std_msgs.UInt16).Data, "Data should be 42")
}

func TestConvertUInt32Readings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/UInt32", map[string]interface{}{"Data": 42})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, uint32(42), m.(*std_msgs.UInt32).Data, "Data should be 42")
}

func TestConvertUInt64Readings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/UInt64", map[string]interface{}{"Data": 42})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, uint64(42), m.(*std_msgs.UInt64).Data, "Data should be 42")
}

func TestConvertFloat32Readings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/Float32", map[string]interface{}{"Data": 42.0})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, float32(42.0), m.(*std_msgs.Float32).Data, "Data should be 42.0")
}

func TestConvertFloat64Readings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/Float64", map[string]interface{}{"Data": 42.0})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, float64(42.0), m.(*std_msgs.Float64).Data, "Data should be 42.0")
}

func TestConvertTimeReadings(t *testing.T) {
	now := time.Now().UTC()
	m, e := ConvertToRosMsg("std_msgs/Time", map[string]interface{}{"Data": now})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, now, m.(*std_msgs.Time).Data, "Data should be 42")
}

func TestConvertDurationReadings(t *testing.T) {
	d := time.Duration(42)
	m, e := ConvertToRosMsg("std_msgs/Duration", map[string]interface{}{"Data": d})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, d, m.(*std_msgs.Duration).Data, "Data should be 42")
}

func TestConvertColorRGBAReadings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/ColorRGBA", map[string]interface{}{"R": 1.0, "G": 0.0, "B": 0.0, "A": 1.0})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, float32(1.0), m.(*std_msgs.ColorRGBA).R, "R should be 1.0")
	assert.Equal(t, float32(0.0), m.(*std_msgs.ColorRGBA).G, "G should be 0.0")
	assert.Equal(t, float32(0.0), m.(*std_msgs.ColorRGBA).B, "B should be 0.0")
	assert.Equal(t, float32(1.0), m.(*std_msgs.ColorRGBA).A, "A should be 1.0")
}

func TestConvertMultiArrayDimensionReadings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/MultiArrayDimension", map[string]interface{}{"Label": "label", "Size": 42, "Stride": 42})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, "label", m.(*std_msgs.MultiArrayDimension).Label, "Label should be 'label'")
	assert.Equal(t, uint32(42), m.(*std_msgs.MultiArrayDimension).Size, "Size should be 42")
	assert.Equal(t, uint32(42), m.(*std_msgs.MultiArrayDimension).Stride, "Stride should be 42")
}

func TestConvertMultiArrayLayoutReadings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/MultiArrayLayout", map[string]interface{}{"DataOffset": 42, "Dim": []std_msgs.MultiArrayDimension{{Label: "label", Size: 42, Stride: 42}}})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, uint32(42), m.(*std_msgs.MultiArrayLayout).DataOffset, "DataOffset should be 42")
	assert.Equal(t, "label", m.(*std_msgs.MultiArrayLayout).Dim[0].Label, "Label should be 'label'")
	assert.Equal(t, uint32(42), m.(*std_msgs.MultiArrayLayout).Dim[0].Size, "Size should be 42")
	assert.Equal(t, uint32(42), m.(*std_msgs.MultiArrayLayout).Dim[0].Stride, "Stride should be 42")
}

func TestConvertByteReadings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/Byte", map[string]interface{}{"Data": 42})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, int8(42), m.(*std_msgs.Byte).Data, "Data should be 42")
}

func TestConvertByteMultiArrayReadings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/ByteMultiArray", map[string]interface{}{"Data": []int8{42}})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, []int8{42}, m.(*std_msgs.ByteMultiArray).Data, "Data should be []int8{42}")
}

func TestConvertCharReadings(t *testing.T) {
	m, e := ConvertToRosMsg("std_msgs/Char", map[string]interface{}{"Data": 42})
	assert.Nil(t, e, "Error should be nil")
	assert.Equal(t, byte(42), m.(*std_msgs.Char).Data, "Data should be 42")
}
