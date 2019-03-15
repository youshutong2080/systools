// +build darwin
// +build !cgo

package host

import "github.com/youshutong2080/systools/internal/common"

func SensorsTemperatures() ([]TemperatureStat, error) {
	return []TemperatureStat{}, common.ErrNotImplementedError
}
