package load_balance

import (
	"fmt"
	"testing"
)

func TestRandom(t *testing.T) {
	nodes := []Node{
		Node{
			Node: "127.0.0.1",
		},
		Node{
			Node: "127.0.0.2",
		},
		Node{
			Node: "127.0.0.3",
		},
	}

	load := NewRandom()

	if err := load.InitNodeList(nodes); err != nil {
		panic(err)
	}

	i := 1
	for {
		if i > 100 {
			break
		}
		target := load.GetNodeAddress()
		_ = target
		i++
	}
	fmt.Println(load.GetStatistics())
}
