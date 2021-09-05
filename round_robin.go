package load_balance

import (
	"sort"
	"sync"
	"sync/atomic"
)

type roundRobin struct {
	nodeList   []Node         //node weight list
	nodeCount  int            //node count
	maxWeight  int            //max weight
	minWeight  int            //min weight
	statistics map[string]int //load balance statistics
	lock       sync.RWMutex
}

func NewRoundRobin() LoadBanlance {
	return &roundRobin{}
}

func (r *roundRobin) InitNodeList(nodeList []Node) (err error) {
	sort.Slice(nodeList, func(i, j int) bool {
		return nodeList[i].Weight > nodeList[j].Weight
	})

	r.nodeList = nodeList
	r.nodeCount = len(r.nodeList)
	if r.nodeCount <= 0 {
		err = ErrNodeEmpty
		return
	}

	r.statistics = make(map[string]int)

	r.maxWeight = 0
	r.minWeight = 1
	for _, v := range r.nodeList {
		r.maxWeight = r.maxInt(r.maxWeight, v.Weight)
		r.minWeight = r.minInt(r.minWeight, v.Weight)
		r.statistics[v.Node] = 0
	}

	return
}

func (r *roundRobin) GetNodeAddress() string {
	idx := r.getNodeIndex()
	node := r.nodeList[idx].Node
	r.addCall(node)

	return node
}

func (r *roundRobin) getNodeIndex() int {
	count := requestAtomic()

	if r.minWeight == r.maxWeight {
		return int(count % uint64(r.nodeCount))
	}

	currentWeight := int(count % uint64(r.maxWeight))
	weightNode := make(map[string]int)
	for _, v := range r.nodeList {
		if v.Weight >= currentWeight {
			weightNode[v.Node] = v.Weight
		}
	}
	length := len(weightNode)

	if length == 1 {
		return 0
	}

	return int(count % uint64(length))
}

func (r *roundRobin) GetStatistics() map[string]int {
	r.lock.RLock()
	defer r.lock.RUnlock()
	res := r.statistics
	return res
}

func (r *roundRobin) addCall(node string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.statistics[node] = r.statistics[node] + 1
}

func (r *roundRobin) maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (r *roundRobin) minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var count uint64 = 0

func requestAtomic() uint64 {
	return atomic.AddUint64(&count, 1)
}
