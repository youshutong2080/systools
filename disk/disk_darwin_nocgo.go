// +build darwin
// +build !cgo

package disk

import "github.com/youshutong2080/systools/internal/common"

func IOCounters(names ...string) (map[string]IOCountersStat, error) {
	return nil, common.ErrNotImplementedError
}
