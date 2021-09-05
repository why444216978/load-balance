package load_balance

import "errors"

type BalanceType string

const (
	balanceTypeRandom     BalanceType = "ramdom"
	balanceTypeRoundRobin BalanceType = "round_robin"
)

type Node struct {
	Node   string
	Weight int
}

type nodeOffset struct {
	Node        string
	Weight      int
	OffsetStart int
	OffsetEnd   int
}

type LoadBanlance interface {
	InitNodeList(nodeList []Node) (err error)

	GetNodeAddress() string

	GetStatistics() map[string]int
}

func New(typ BalanceType) (LoadBanlance, error) {
	switch typ {
	case balanceTypeRandom:
		return NewRandom(), nil
	case balanceTypeRoundRobin:
		return NewRoundRobin(), nil
	}

	return nil, errors.New("type error")
}
