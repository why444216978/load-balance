package load_balance

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
