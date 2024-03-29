// +build darwin

package host

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/youshutong2080/systools/internal/common"
	"github.com/youshutong2080/systools/process"
	"log"
	"regexp"
	"errors"
)

// from utmpx.h
const USER_PROCESS = 7

func Info() (*InfoStat, error) {
	ret := &InfoStat{
		OS:             runtime.GOOS,
		PlatformFamily: "darwin",
	}

	hostname, err := os.Hostname()
	if err == nil {
		ret.Hostname = hostname
	}

	uname, err := exec.LookPath("uname")
	if err == nil {
		out, err := invoke.Command(uname, "-r")
		if err == nil {
			ret.KernelVersion = strings.ToLower(strings.TrimSpace(string(out)))
		}
	}

	platform, family, pver, err := PlatformInformation()
	if err == nil {
		ret.Platform = platform
		ret.PlatformFamily = family
		ret.PlatformVersion = pver
	}

	system, role, err := Virtualization()
	if err == nil {
		ret.VirtualizationSystem = system
		ret.VirtualizationRole = role
	}

	boot, err := BootTime()
	if err == nil {
		ret.BootTime = boot
		ret.Uptime = uptime(boot)
	}

	procs, err := process.Pids()
	if err == nil {
		ret.Procs = uint64(len(procs))
	}

	values, err := common.DoSysctrl("kern.uuid")
	if err == nil && len(values) == 1 && values[0] != "" {
		ret.HostID = strings.ToLower(values[0])
	}

	return ret, nil
}

func BootTime() (uint64, error) {
	if cachedBootTime != 0 {
		return cachedBootTime, nil
	}
	values, err := common.DoSysctrl("kern.boottime")
	if err != nil {
		return 0, err
	}
	// ex: { sec = 1392261637, usec = 627534 } Thu Feb 13 12:20:37 2014
	v := strings.Replace(values[2], ",", "", 1)
	boottime, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}
	cachedBootTime = uint64(boottime)

	return cachedBootTime, nil
}

func uptime(boot uint64) uint64 {
	return uint64(time.Now().Unix()) - boot
}

func Uptime() (uint64, error) {
	boot, err := BootTime()
	if err != nil {
		return 0, err
	}
	return uptime(boot), nil
}

func Users() ([]UserStat, error) {
	utmpfile := "/var/run/utmpx"
	var ret []UserStat

	file, err := os.Open(utmpfile)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return ret, err
	}

	u := Utmpx{}
	entrySize := int(unsafe.Sizeof(u))
	count := len(buf) / entrySize

	for i := 0; i < count; i++ {
		b := buf[i*entrySize : i*entrySize+entrySize]

		var u Utmpx
		br := bytes.NewReader(b)
		err := binary.Read(br, binary.LittleEndian, &u)
		if err != nil {
			continue
		}
		if u.Type != USER_PROCESS {
			continue
		}
		user := UserStat{
			User:     common.IntToString(u.User[:]),
			Terminal: common.IntToString(u.Line[:]),
			Host:     common.IntToString(u.Host[:]),
			Started:  int(u.Tv.Sec),
		}
		ret = append(ret, user)
	}

	return ret, nil

}

func PlatformInformation() (string, string, string, error) {
	platform := ""
	family := ""
	pver := ""

	sw_vers, err := exec.LookPath("sw_vers")
	if err != nil {
		return "", "", "", err
	}
	uname, err := exec.LookPath("uname")
	if err != nil {
		return "", "", "", err
	}

	out, err := invoke.Command(uname, "-s")
	if err == nil {
		platform = strings.ToLower(strings.TrimSpace(string(out)))
	}

	out, err = invoke.Command(sw_vers, "-productVersion")
	if err == nil {
		pver = strings.ToLower(strings.TrimSpace(string(out)))
	}

	return platform, family, pver, nil
}

func Virtualization() (string, string, error) {
	system := ""
	role := ""

	return system, role, nil
}


func Cmdexec(command string) (string,error) {
	cmd := exec.Command("/bin/sh", "-c", command) //调用Command函数
	var out bytes.Buffer //缓冲字节
	cmd.Stdout = &out //标准输出
	err := cmd.Run() //运行指令 ，做判断
	if err != nil {
		return "",err
	}
	//fmt.Printf("%s", out.String()) //输出执行结果
	return out.String(),err
}

func CmdexecIostat(disk string) ([]string,error){
	cmdOut,err :=Cmdexec("iostat -d "+disk)
	if err != nil{
		return nil,err
	}
	cmdOuts := strings.Split(cmdOut,"\n")
	if len(cmdOuts) >1 {
		val := cmdOuts[2]
		val = strings.TrimSpace(val)
		reg := regexp.MustCompile(`[ ]+`)
		val = reg.ReplaceAllString(val, " ")
		log.Println(val)
		vals := strings.Split(val, " ")
		return vals,err
	}else {
		return nil,errors.New("cmdOuts len is not greater than 1!")
	}
}

func SystemUptime() (days, hours, mins int64, err error) {
	var  up_time uint64
	up_time, err = Uptime()
	if err != nil {
		return
	}


	secStr := strconv.FormatInt(int64(up_time), 10)
	var secF float64
	secF, err = strconv.ParseFloat(secStr, 64)
	if err != nil {
		return
	}

	minTotal := secF / 60.0
	hourTotal := minTotal / 60.0

	days = int64(hourTotal / 24.0)
	hours = int64(hourTotal) - days*24
	mins = int64(minTotal) - (days * 60 * 24) - (hours * 60)

	return
}