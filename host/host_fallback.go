// +build !darwin,!linux,!freebsd,!openbsd,!solaris,!windows

package host

import "github.com/youshutong2080/systools/internal/common"

func Info() (*InfoStat, error) {
	return nil, common.ErrNotImplementedError
}

func BootTime() (uint64, error) {
	return 0, common.ErrNotImplementedError
}

func Uptime() (uint64, error) {
	return 0, common.ErrNotImplementedError
}

func Users() ([]UserStat, error) {
	return []UserStat{}, common.ErrNotImplementedError
}
