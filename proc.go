package proc

import (
	"github.com/gogather/com"
	"io/ioutil"
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

func (pro *Process) GetPid() {
	return pro.pid
}

func (pro *Process) GetParentProc() *Process {
	return pro.pp
}

func (pro *Process) GetChildrenProc() []*Process {
	return pro.cp
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
	procMap map[int64]Process
)

func reloadProcessTree() {
	scanProc()
	scanRelative()
}

func newProc(pid int64) {
	content := com.ReadFileString(fmt.Sprintf("/proc/%d/status", pid))
	proc := Process{pid: pid}
	proc.parseProcInfo(content)
	procMap[pid] = proc
}

func scanProc() {
	files, _ := ioutil.ReadDir(`/proc/`)
	procMap = map[int64]Process{}
	for _, file := range files {
		pid, err := strconv.ParseInt(file.Name(), 10, 64)
		if file.IsDir() && err == nil {
			procMap[pid] = newProc(pid)
		}
	}
}

func scanRelative() {
	for key, proc := range procMap {
		s_ppid := proc.info["PPid"]
		ppid, err := strconv.ParseInt(s_ppid, 10, 64)
		pproc := procMap[ppid]
		if pproc != nil {
			proc.pp = &pproc
			pproc.cp = append(pproc.cp, &proc)
		}
	}
}

// get process
func GetProc(pid int64) *Process {
	reloadProcessTree()
	return &procMap[pid]
}
