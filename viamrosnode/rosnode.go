package viamrosnode

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/bluenviron/goroslib/v2"
)

var lock *sync.Mutex
var i int
var nodes map[string]*goroslib.Node

func init() {
	lock = &sync.Mutex{}
	i = 0
	nodes = make(map[string]*goroslib.Node)
}

func GetInstance(primary string) (*goroslib.Node, error) {
	lock.Lock()
	defer lock.Unlock()
	node, ok := nodes[primary]
	if ok {
		return node, nil
	} else {
		node, err := goroslib.NewNode(goroslib.NodeConf{
			Name:          strings.Join([]string{"viamrosnode_", primary, strconv.Itoa(i)}, ""),
			MasterAddress: primary,
		})
		if err != nil {
			return nil, err
		}

		nodes[primary] = node
		i = i + 1
		return node, nil
	}
}

func ShutdownNodes() {
	lock.Lock()
	defer lock.Unlock()
	for primary, node := range nodes {
		fmt.Printf("Closing %s", primary)
		node.Close()
	}
	// We have to recreate this map because we can't modify the map while iterating over it
	nodes = make(map[string]*goroslib.Node)
}
