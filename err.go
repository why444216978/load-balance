package load_balance

import "errors"

var (
	ErrNodeEmpty   = errors.New("node is empty")
	ErrTotalWeight = errors.New("totalWeight = 0")
)
