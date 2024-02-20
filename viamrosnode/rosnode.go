package viamrosnode

import (
	"strconv"
	"strings"
	"sync"

	"github.com/bluenviron/goroslib/v2"
)

var lock *sync.Mutex = &sync.Mutex{}
var i int = 0

func getInsance(primary string, nodeConfig goroslib.NodeConf) (*goroslib.Node, error) {
	node, err := goroslib.NewNode(nodeConfig)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func GetInsanceWithLogCallback(primary string, host string, logCallback func(level goroslib.LogLevel, msg string)) (*goroslib.Node, error) {
	lock.Lock()
	defer lock.Unlock()
	defer func() { i = i + 1 }()
	nodeConfig := goroslib.NodeConf{
		Name:            strings.Join([]string{"viamrosnode_", primary, strconv.Itoa(i)}, ""),
		MasterAddress:   primary,
		Host:            host,
		LogDestinations: goroslib.LogDestinationCallback,
		OnLog:           logCallback,
	}
	return getInsance(primary, nodeConfig)
}

func GetInstance(primary string, host string) (*goroslib.Node, error) {
	lock.Lock()
	defer lock.Unlock()
	defer func() { i = i + 1 }()
	nodeConfig := goroslib.NodeConf{
		Name:          strings.Join([]string{"viamrosnode_", primary, strconv.Itoa(i)}, ""),
		MasterAddress: primary,
		Host:          host,
	}
	return getInsance(primary, nodeConfig)
}
