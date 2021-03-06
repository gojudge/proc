package proc

import (
	"fmt"
	"github.com/gogather/com"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Process struct {
	pid    int64
	status string
	pp     *Process
	cp     []*Process
	info   map[string]string
}

func (pro *Process) GetPid() int64 {
	return pro.pid
}

func (pro *Process) GetName() string {
	return pro.info["Name"]
}

func (pro *Process) GetParentProc() *Process {
	return pro.pp
}

func (pro *Process) GetChildrenProc() []*Process {
	return pro.cp
}

// kill process chain with reverse order
func (pro *Process) KillProcChainReverse() {
	if pro.pid == 1 {
		return
	}

	for {
		pro.killNoSonTailProc()
		if len(pro.cp) <= 0 {
			break
		}
	}

	pro.killNoSonTailProc()
}

// kill process which without son process in the tails
func (pro *Process) killNoSonTailProc() {
	cp := pro.cp
	if len(cp) > 0 {
		for i := 0; i < len(cp); i++ {
			cp[i].killNoSonTailProc()
		}
	} else {
		pro.Kill()
	}
}

func (pro *Process) Kill() {
	proc, err := os.FindProcess(int(pro.pid))
	if err != nil {
		return
	}

	proc.Kill()
	pro.pp.removeChildProcessInfoByPid(pro.pid)
}

// remove process info from parents process in cp array
func (pro *Process) removeChildProcessInfoByPid(pid int64) {
	var ncps []*Process
	cps := pro.cp
	for i := 0; i < len(cps); i++ {
		if cps[i].pid != pid {
			ncps = append(ncps, cps[i])
		}
	}
	pro.cp = ncps
}

func (pro *Process) parseProcInfo(content string) {
	lines := strings.Split(content, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		kv := strings.Split(line, ":")
		if len(kv) >= 2 {
			key := com.Strim(strings.TrimSpace(kv[0]))
			value := com.Strim(strings.TrimSpace(kv[1]))
			if key != "" {
				pro.info[key] = value
			}
		}
	}
}

var (
	procMap   map[int64]*Process
	procSlice []Process
)

func reloadProcessTree() {
	scanProc()
	scanRelative()
}

func newProc(pid int64) {
	content, err := com.ReadFileString(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		content = ""
	}

	proc := Process{pid: pid, info: map[string]string{}}
	proc.parseProcInfo(content)
	procSlice = append(procSlice, proc)
	procMap[pid] = &procSlice[len(procSlice)-1]
}

func scanProc() {
	files, _ := ioutil.ReadDir(`/proc/`)
	procMap = map[int64]*Process{}
	for _, file := range files {
		pid, err := strconv.ParseInt(file.Name(), 10, 64)
		if file.IsDir() && err == nil {
			newProc(pid)
		}
	}

}

func scanRelative() {
	for _, proc := range procMap {
		s_ppid := proc.info["PPid"]
		ppid, err := strconv.ParseInt(s_ppid, 10, 64)
		if err == nil {
			_, ok := procMap[ppid]
			if ok {
				proc.pp = procMap[ppid]
				procMap[ppid].cp = append(procMap[ppid].cp, proc)
			}
		}
	}
}

// get process
func GetProc(pid int64) *Process {
	reloadProcessTree()
	p, ok := procMap[pid]
	if ok {
		return p
	} else {
		return nil
	}
}
