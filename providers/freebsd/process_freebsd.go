// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

//+build,freebsd,cgo

package freebsd

// #cgo LDFLAGS: -lkvm -lprocstat
//#include <sys/types.h>
//#include <sys/sysctl.h>
//#include <sys/time.h>
//#include <sys/param.h>
//#include <sys/queue.h>
//#include <sys/socket.h>
//#include <sys/user.h>
//
//#include <libprocstat.h>
//#include <string.h>
//struct kinfo_proc getProcInfoAt(struct kinfo_proc *procs, unsigned int index) {
//  return procs[index];
//}
//unsigned int countArrayItems(char **items) {
//  unsigned int i = 0;
//  for (i = 0; items[i] != NULL; ++i);
//  return i;
//}
//char * itemAtIndex(char **items, unsigned int index) {
//  return items[index];
//}
//unsigned int countFileStats(struct filestat_list *head) {
//  unsigned int count = 0;
//  struct filestat *fst;
//  STAILQ_FOREACH(fst, head, next) {
//    ++count;
//  }
//
//  return count;
//}
//void copyFileStats(struct filestat_list *head, struct filestat *out, unsigned int size) {
//  unsigned int index = 0;
//  struct filestat *fst;
//  STAILQ_FOREACH(fst, head, next) {
//    if (!size) {
//      break;
//    }
//    memcpy(out, fst, sizeof(*fst));
//    ++out;
//    --size;
//  }
//}
//
import "C"

import (
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/elastic/go-sysinfo/types"
)

func getProcInfo(op int, arg int) ([]process, error) {
	procstat, err := C.procstat_open_sysctl()

	if procstat == nil {
		return nil, errors.Wrap(err, "failed to open procstat sysctl")
	}
	defer C.procstat_close(procstat)

	var count C.uint = 0
	kprocs, err := C.procstat_getprocs(procstat, C.int(op), C.int(arg), &count)
	if kprocs == nil {
		return nil, errors.Wrap(err, "getprocs failed")
	}
	defer C.procstat_freeprocs(procstat, kprocs)

	procs := make([]process, count)
	var index C.uint
	for index = 0; index < count; index++ {
		proc := C.getProcInfoAt(kprocs, index)
		procs[index].kinfo = proc
		procs[index].pid = int(proc.ki_pid)
	}

	return procs, nil
}

func copyArray(from **C.char) []string {
	if from == nil {
		return nil
	}

	count := C.countArrayItems(from)
	out := make([]string, count)

	for index := C.uint(0); index < count; index++ {
		out[index] = C.GoString(C.itemAtIndex(from, index))
	}

	return out
}

func makeMap(from []string) map[string]string {
	out := make(map[string]string, len(from))

	for _, env := range from {
		parts := strings.Split(env, "=")
		if len(parts) > 1 {
			out[parts[0]] = parts[1]
		}
	}

	return out
}

func getProcEnv(p *process) (map[string]string, error) {
	procstat, err := C.procstat_open_sysctl()

	if procstat == nil {
		return nil, errors.Wrap(err, "failed to open procstat sysctl")
	}
	defer C.procstat_close(procstat)

	env, err := C.procstat_getenvv(procstat, &p.kinfo, 0)
	defer C.procstat_freeenvv(procstat)

	return makeMap(copyArray(env)), err
}

func getProcArgs(p *process) ([]string, error) {
	procstat, err := C.procstat_open_sysctl()

	if procstat == nil {
		return nil, errors.Wrap(err, "failed to open procstat sysctl")
	}
	defer C.procstat_close(procstat)

	args, err := C.procstat_getargv(procstat, &p.kinfo, 0)
	defer C.procstat_freeargv(procstat)

	return copyArray(args), err
}

func getProcPathname(p *process) (string, error) {
	procstat, err := C.procstat_open_sysctl()

	if procstat == nil {
		return "", errors.Wrap(err, "failed to open procstat sysctl")
	}
	defer C.procstat_close(procstat)

	const maxlen = uint(1024)
	out := make([]C.char, maxlen)
	if res, err := C.procstat_getpathname(procstat, &p.kinfo, &out[0], C.ulong(maxlen)); res != 0 {
		return "", err
	}
	return C.GoString(&out[0]), nil
}

func getFileStats(fileStats *C.struct_filestat_list) []C.struct_filestat {
	count := C.countFileStats(fileStats)

	if count < 1 {
		return nil
	}

	out := make([]C.struct_filestat, count)

	C.copyFileStats(fileStats, &out[0], count)
	return out
}

func getProcCWD(p *process) (string, error) {
	procstat, err := C.procstat_open_sysctl()

	if procstat == nil {
		return "", errors.Wrap(err, "failed to open procstat sysctl")
	}
	defer C.procstat_close(procstat)

	fs, err := C.procstat_getfiles(procstat, &p.kinfo, 0)
	if fs == nil {
		return "", errors.Wrap(err, "failed to get files")
	}

	defer C.procstat_freefiles(procstat, fs)

	files := getFileStats(fs)
	for _, f := range files {
		if f.fs_uflags == C.PS_FST_UFLAG_CDIR {
			return C.GoString(f.fs_path), nil
		}
	}

	return "", nil
}

type process struct {
	pid   int
	kinfo C.struct_kinfo_proc
}

func timevalToDuration(tm C.struct_timeval) time.Duration {
	return (time.Duration(tm.tv_sec) * 1000000) + (time.Duration(tm.tv_usec) * 1000)
}

func (p *process) CPUTime() (types.CPUTimes, error) {
	procs, err := getProcInfo(C.KERN_PROC_PID, p.PID())

	if err != nil {
		return types.CPUTimes{}, err
	}
	p.kinfo = procs[0].kinfo

	return types.CPUTimes{
		User:   timevalToDuration(p.kinfo.ki_rusage.ru_utime),
		System: timevalToDuration(p.kinfo.ki_rusage.ru_stime),
	}, nil
}

func (p *process) Info() (types.ProcessInfo, error) {
	procs, err := getProcInfo(C.KERN_PROC_PID, p.PID())

	if err != nil {
		return types.ProcessInfo{}, err
	}
	p.kinfo = procs[0].kinfo

	cwd, err := getProcCWD(p)
	if err != nil {
		return types.ProcessInfo{}, err
	}

	args, err := getProcArgs(p)
	if err != nil {
		return types.ProcessInfo{}, err
	}

	exe, _ := getProcPathname(p)

	return types.ProcessInfo{
		Name:      C.GoString(&p.kinfo.ki_comm[0]),
		PID:       int(p.kinfo.ki_pid),
		PPID:      int(p.kinfo.ki_ppid),
		CWD:       cwd,
		Exe:       exe,
		Args:      args,
		StartTime: time.Unix(int64(p.kinfo.ki_start.tv_sec), int64(p.kinfo.ki_start.tv_usec)*1000),
	}, nil
}

func (p *process) Memory() (types.MemoryInfo, error) {
	procs, err := getProcInfo(C.KERN_PROC_PID, p.PID())

	if err != nil {
		return types.MemoryInfo{}, err
	}
	p.kinfo = procs[0].kinfo

	return types.MemoryInfo{
		Resident: uint64(p.kinfo.ki_rssize),
		Virtual:  uint64(p.kinfo.ki_size),
	}, nil
}

func (p *process) User() (types.UserInfo, error) {
	procs, err := getProcInfo(C.KERN_PROC_PID, p.PID())

	if err != nil {
		return types.UserInfo{}, err
	}

	p.kinfo = procs[0].kinfo

	return types.UserInfo{
		UID:  strconv.FormatUint(uint64(p.kinfo.ki_ruid), 10),
		EUID: strconv.FormatUint(uint64(p.kinfo.ki_uid), 10),
		SUID: strconv.FormatUint(uint64(p.kinfo.ki_svuid), 10),
		GID:  strconv.FormatUint(uint64(p.kinfo.ki_rgid), 10),
		EGID: strconv.FormatUint(uint64(p.kinfo.ki_groups[0]), 10),
		SGID: strconv.FormatUint(uint64(p.kinfo.ki_svgid), 10),
	}, nil
}

func (p *process) PID() int {
	return p.pid
}

func (p *process) OpenHandles() ([]string, error) {
	procstat := C.procstat_open_sysctl()

	if procstat == nil {
		return nil, errors.New("failed to open procstat sysctl")
	}
	defer C.procstat_close(procstat)

	fs := C.procstat_getfiles(procstat, &p.kinfo, 0)
	defer C.procstat_freefiles(procstat, fs)

	files := getFileStats(fs)
	names := make([]string, 0, len(files))

	for _, file := range files {
		if file.fs_uflags == 0 {
			names = append(names, C.GoString(file.fs_path))
		}
	}

	return names, nil
}

func (p *process) OpenHandleCount() (int, error) {
	procstat := C.procstat_open_sysctl()

	if procstat == nil {
		return 0, errors.New("failed to open procstat sysctl")
	}
	defer C.procstat_close(procstat)

	fs := C.procstat_getfiles(procstat, &p.kinfo, 0)
	defer C.procstat_freefiles(procstat, fs)
	return int(C.countFileStats(fs)), nil
}

func (p *process) Environment() (map[string]string, error) {
	return getProcEnv(p)
}

func (s freebsdSystem) Processes() ([]types.Process, error) {
	procs, err := getProcInfo(C.KERN_PROC_PROC, 0)
	out := make([]types.Process, 0, len(procs))

	for _, proc := range procs {
		out = append(out, &process{
			pid:   proc.pid,
			kinfo: proc.kinfo,
		})
	}

	return out, err
}

func (s freebsdSystem) Process(pid int) (types.Process, error) {
	p := process{pid: pid}
	return &p, nil
}

func (s freebsdSystem) Self() (types.Process, error) {
	return s.Process(os.Getpid())
}

const kernCptimeMIB = "kern.cp_time"
const kernClockrateMIB = "kern.clockrate"

func Cptime() (map[string]uint64, error) {
	var clock clockInfo

	if err := sysctlByName(kernClockrateMIB, &clock); err != nil {
		return make(map[string]uint64), errors.Wrap(err, "failed to get kern.clockrate")
	}

	cptime, err := syscall.Sysctl(kernCptimeMIB)
	if err != nil {
		return make(map[string]uint64), errors.Wrap(err, "failed to get kern.cp_time")
	}

	cpMap := make(map[string]uint64)

	times := strings.Split(cptime, " ")
	names := [5]string{"User", "Nice", "System", "IRQ", "Idle"}

	for index, time := range times {
		i, err := strconv.ParseUint(time, 10, 64)

		if err != nil {
			return cpMap, errors.Wrap(err, "error parsing kern.cp_time")
		}

		cpMap[names[index]] = i * uint64(clock.Tick) * 1000
	}

	return cpMap, nil
}
