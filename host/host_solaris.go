package host

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/youshutong2080/systools/internal/common"
)

func Info() (*InfoStat, error) {
	result := &InfoStat{
		OS: runtime.GOOS,
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	result.Hostname = hostname

	// Parse versions from output of `uname(1)`
	uname, err := exec.LookPath("/usr/bin/uname")
	if err != nil {
		return nil, err
	}

	out, err := invoke.Command(uname, "-srv")
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(string(out))
	if len(fields) >= 1 {
		result.PlatformFamily = fields[0]
	}
	if len(fields) >= 2 {
		result.KernelVersion = fields[1]
	}
	if len(fields) == 3 {
		result.PlatformVersion = fields[2]
	}

	// Find distribution name from /etc/release
	fh, err := os.Open("/etc/release")
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	sc := bufio.NewScanner(fh)
	if sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch {
		case strings.HasPrefix(line, "SmartOS"):
			result.Platform = "SmartOS"
		case strings.HasPrefix(line, "OpenIndiana"):
			result.Platform = "OpenIndiana"
		case strings.HasPrefix(line, "OmniOS"):
			result.Platform = "OmniOS"
		case strings.HasPrefix(line, "Open Storage"):
			result.Platform = "NexentaStor"
		case strings.HasPrefix(line, "Solaris"):
			result.Platform = "Solaris"
		case strings.HasPrefix(line, "Oracle Solaris"):
			result.Platform = "Solaris"
		default:
			result.Platform = strings.Fields(line)[0]
		}
	}

	// Find the boot time and calculate uptime relative to it
	bootTime, err := BootTime()
	if err != nil {
		return nil, err
	}
	result.BootTime = bootTime
	result.Uptime = uptimeSince(bootTime)

	// Count number of processes based on the number of entries in /proc
	dirs, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	result.Procs = uint64(len(dirs))

	return result, nil
}

var kstatMatch = regexp.MustCompile(`([^\s]+)[\s]+([^\s]*)`)

func BootTime() (uint64, error) {
	kstat, err := exec.LookPath("/usr/bin/kstat")
	if err != nil {
		return 0, err
	}

	out, err := invoke.Command(kstat, "-p", "unix:0:system_misc:boot_time")
	if err != nil {
		return 0, err
	}

	kstats := kstatMatch.FindAllStringSubmatch(string(out), -1)
	if len(kstats) != 1 {
		return 0, fmt.Errorf("expected 1 kstat, found %d", len(kstats))
	}

	return strconv.ParseUint(kstats[0][2], 10, 64)
}

func Uptime() (uint64, error) {
	bootTime, err := BootTime()
	if err != nil {
		return 0, err
	}
	return uptimeSince(bootTime), nil
}

func uptimeSince(since uint64) uint64 {
	return uint64(time.Now().Unix()) - since
}

func Users() ([]UserStat, error) {
	return []UserStat{}, common.ErrNotImplementedError
}

func SensorsTemperatures() ([]TemperatureStat, error) {
	return []TemperatureStat{}, common.ErrNotImplementedError
}
