package types

type GoInfo struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	MaxProcs int    `json:"max_procs"`
	Version  string `json:"version"`
}
