// +build linux

package load

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/youshutong2080/systools/internal/common"
)

func Avg() (*AvgStat, error) {
	filename := common.HostProc("loadavg")
	line, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	values := strings.Fields(string(line))

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
		Load1:  load1,
		Load5:  load5,
		Load15: load15,
	}

	return ret, nil
}

// Misc returnes miscellaneous host-wide statistics.
// Note: the name should be changed near future.
func Misc() (*MiscStat, error) {
	filename := common.HostProc("stat")
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	ret := &MiscStat{}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		v, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			continue
		}
		switch fields[0] {
		case "procs_running":
			ret.ProcsRunning = int(v)
		case "procs_blocked":
			ret.ProcsBlocked = int(v)
		case "ctxt":
			ret.Ctxt = int(v)
		default:
			continue
		}

	}

	return ret, nil
}
