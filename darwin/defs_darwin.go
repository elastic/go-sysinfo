// +build ignore

package darwin

/*
#include <libproc.h>
*/
import "C"

type processState uint32

const (
	stateSIDL processState = iota + 1
	stateRun
	stateSleep
	stateStop
	stateZombie
)

const argMax = C.ARG_MAX

type bsdInfo C.struct_proc_bsdinfo

type procTaskInfo C.struct_proc_taskinfo

type procTaskAllInfo C.struct_proc_taskallinfo

type vinfoStat C.struct_vinfo_stat

type fsid C.struct_fsid

type vnodeInfo C.struct_vnode_info

type vnodeInfoPath C.struct_vnode_info_path

type procVnodePathInfo C.struct_proc_vnodepathinfo
