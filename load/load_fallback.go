// +build !darwin,!linux,!freebsd,!openbsd,!windows

package load

import "github.com/youshutong2080/systools/internal/common"

func Avg() (*AvgStat, error) {
	return nil, common.ErrNotImplementedError
}

func Misc() (*MiscStat, error) {
	return nil, common.ErrNotImplementedError
}
