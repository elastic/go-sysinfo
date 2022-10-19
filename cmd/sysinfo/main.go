package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
)

const (
	hostTmpl = `Host Info:

Architecture:   {{ .Info.Architecture }}
Boot Time:      {{ .Info.BootTime }} ({{ .Info.BootTime | ago }} ago)
{{ with .Info.Containerized -}}
Containerized:  {{ .Info.Containerized }}
{{ end -}}
Hostname:       {{ .Info.Hostname }}
Unique ID:      {{ .Info.UniqueID }}
OS:             {{ .Info.OS }}
Kernel Version: {{ .Info.KernelVersion }}
Timezone:       {{ .Info.Timezone }} (Offset: {{ .Info.TimezoneOffsetSec }})
IP Addresses:   {{ .Info.IPs | join " " }}
MAC Addresses:  {{ .Info.MACs | join " " }}
`
	cpuTmpl = `CPU Usage:

User    {{ .User }}
System  {{ .System }}
Idle    {{ .Idle }}
IOWait  {{ .IOWait }}
IRQ     {{ .IRQ }}
Nice    {{ .Nice }}
SoftIRQ {{ .SoftIRQ }}
Steal   {{ .Steal }}
`
	memTmpl = `Memory Usage:

Total        {{ .Total | bytes }}
Used         {{ .Used | bytes }}
Available    {{ .Available | bytes }}
Free         {{ .Free | bytes }}
VirtualTotal {{ .VirtualTotal | bytes }}
VirtualUsed  {{ .VirtualUsed | bytes }}
VirtualFree  {{ .VirtualFree | bytes }}
`
	loadTmpl = `Load Average: {{ .One }} {{ .Five }} {{ .Fifteen }}
`
)

func main() {
	host, err := sysinfo.Host()
	if err != nil {
		log.Fatalf("error getting host info: %s", err)
	}
	if err := printHostInfo(os.Stdout, host); err != nil {
		log.Fatalf("error printing host info: %s", err)
	}
}

func printHostInfo(w io.Writer, host types.Host) error {
	if err := render(hostTmpl, host, w); err != nil {
		return fmt.Errorf("error rendering host information: %w", err)
	}

	cpuTimes, err := host.CPUTime()
	if err != nil {
		return fmt.Errorf("error getting host cpu times: %w", err)
	}
	fmt.Fprintln(w)
	if err := render(cpuTmpl, cpuTimes, w); err != nil {
		return fmt.Errorf("error rendering host cpu time: %w", err)
	}

	memory, err := host.Memory()
	if err != nil {
		return fmt.Errorf("error getting host memory: %w", err)
	}
	fmt.Fprintln(w)
	if err := render(memTmpl, memory, w); err != nil {
		return fmt.Errorf("error rendering host memory: %w", err)
	}

	load, err := host.LoadAverage()
	if err != nil {
		return fmt.Errorf("error getting host load: %w", err)
	}
	fmt.Fprintln(w)
	if err := render(loadTmpl, load, w); err != nil {
		return fmt.Errorf("error rendering host load: %w", err)
	}

	return nil
}

func ago(t time.Time) string {
	return time.Since(t).String()
}

func join(sep string, xs []string) string {
	return strings.Join(xs, sep)
}

func render(tmpl string, ctx interface{}, w io.Writer) error {
	funcs := template.FuncMap{
		"ago":   ago,
		"join":  join,
		"bytes": humanize.Bytes,
	}

	t, err := template.New("tmpl").Funcs(funcs).Parse(tmpl)
	if err != nil {
		return err
	}

	return t.Execute(w, ctx)
}
