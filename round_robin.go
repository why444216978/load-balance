package load_balance

import "sync/atomic"

type roundRobin struct {
	nodeList []Node //node weight list

	nodeCount int //node count

	maxWeight int //max weight

	minWeight int //min weight

	statistics map[string]int //load balance statistics
}

func NewRoundRobin() LoadBanlance {
	return &roundRobin{}
}

func (r *roundRobin) InitNodeList(nodeList []Node) (err error) {
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
	}

	return
}

func (r *roundRobin) GetNodeAddress() string {
	idx := r.getNodeIndex()
	address := r.nodeList[idx].Node

	if _, ok := r.statistics[address]; !ok {
		r.statistics[address] = 1
		return address
	}

	r.statistics[address] = r.statistics[address] + 1
	return address
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
	return r.statistics
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
