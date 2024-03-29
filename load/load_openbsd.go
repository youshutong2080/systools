// +build openbsd

package load

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/youshutong2080/systools/internal/common"
)

func Avg() (*AvgStat, error) {
	values, err := common.DoSysctrl("vm.loadavg")
	if err != nil {
		return nil, err
	}

	load1, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return nil, err
	}
	load5, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return nil, err
	}
	load15, err := strconv.ParseFloat(values[2], 64)
	if err != nil {
		return nil, err
	}

	ret := &AvgStat{
		Load1:  float64(load1),
		Load5:  float64(load5),
		Load15: float64(load15),
	}

	return ret, nil
}

// Misc returnes miscellaneous host-wide statistics.
// darwin use ps command to get process running/blocked count.
// Almost same as Darwin implementation, but state is different.
func Misc() (*MiscStat, error) {
	bin, err := exec.LookPath("ps")
	if err != nil {
		return nil, err
	}
	out, err := invoke.Command(bin, "axo", "state")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")

	ret := MiscStat{}
	for _, l := range lines {
		if strings.Contains(l, "R") {
			ret.ProcsRunning++
		} else if strings.Contains(l, "D") {
			ret.ProcsBlocked++
		}
	}

	return &ret, nil
}
