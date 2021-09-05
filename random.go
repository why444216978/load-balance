package load_balance

import (
	"math/rand"
	"sort"
	"sync"
)

type random struct {
	nodeList    []Node         // node weight list
	nodeCount   int            //node count
	offsetList  []nodeOffset   //node offset
	totalWeight int            //total weight
	forCount    int            //search offset count
	sameWeight  bool           //all node weight is same
	statistics  map[string]int //load balance statistics
	lock        sync.RWMutex
}

// New new load balance object
func NewRandom() LoadBanlance {
	return &random{}
}

// Init init
func (r *random) InitNodeList(nodeList []Node) (err error) {
	r.nodeCount = len(nodeList)
	if r.nodeCount <= 0 {
		err = ErrNodeEmpty
		return
	}

	r.nodeList = nodeList
	r.offsetList = make([]nodeOffset, 0)
	r.totalWeight = 0
	r.forCount = 0
	r.sameWeight = true
	r.statistics = make(map[string]int)

	lastWeight := 0
	for k, v := range nodeList {
		tmp := nodeOffset{}

		if k == 0 {
			tmp = nodeOffset{
				Node:        v.Node,
				Weight:      v.Weight,
				OffsetStart: 0,
				OffsetEnd:   v.Weight,
			}
			lastWeight = v.Weight
		} else {
			tmp = nodeOffset{
				Node:        v.Node,
				Weight:      v.Weight,
				OffsetStart: r.totalWeight + 1,
				OffsetEnd:   r.totalWeight + v.Weight,
			}
			if lastWeight != v.Weight {
				r.sameWeight = false
			}
		}

		r.totalWeight = r.totalWeight + v.Weight
		r.offsetList = append(r.offsetList, tmp)

		r.statistics[v.Node] = 0
	}

	if r.totalWeight < 0 {
		err = ErrTotalWeight
		return
	}

	if r.sameWeight {
		return
	}

	//sort by weight
	sort.Slice(r.offsetList, func(i, j int) bool {
		return r.offsetList[i].Weight > r.offsetList[j].Weight
	})

	return
}

// GetNodeAddress get node address
func (r *random) GetNodeAddress() string {
	node := ""
	if r.sameWeight {
		idx := rand.Intn(r.nodeCount)
		node = r.nodeList[idx].Node
	} else {
		idx := rand.Intn(r.totalWeight)
		for _, v := range r.offsetList {
			r.forCount = r.forCount + 1
			if idx >= v.OffsetStart && idx <= v.OffsetEnd {
				node = v.Node
				break
			}
		}
	}

	r.addCall(node)

	return node
}

// GetStatistics get load balance statistic
func (r *random) GetStatistics() map[string]int {
	r.lock.RLock()
	defer r.lock.RUnlock()
	res := r.statistics
	return res
}

// addCall
func (r *random) addCall(node string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.statistics[node] = r.statistics[node] + 1
}
