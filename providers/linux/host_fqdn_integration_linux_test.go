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

package linux

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

const wantHostname = "debian"
const wantDomainCgo = "cgo"

func TestHost_FQDN_Domain_Cgo(t *testing.T) {
	host, err := newLinuxSystem("").Host()
	if err != nil {
		t.Fatal(fmt.Errorf("could not het host information: %w", err))
	}

	got := host.Info()
	if got.Hostname != wantHostname {
		t.Errorf("got wrong hostname want: %q, got %q", wantHostname, got.Hostname)
	}
	if got.Domain != wantDomainCgo {
		t.Errorf("got wrong domain want: %q, got %q", wantDomainCgo, got.Domain)
	}
	if got.FQDN != fmt.Sprintf("%s.%s", wantHostname, wantDomainCgo) {
		t.Errorf("FQDN shpould not be empty")
	}
}

func TestHost_FQDN_No_Domain_Cgo(t *testing.T) {
	host, err := newLinuxSystem("").Host()
	if err != nil {
		t.Fatal(fmt.Errorf("could not het host information: %w", err))
	}

	got := host.Info()
	if got.Hostname != wantHostname {
		t.Errorf("got wrong hostname want: %s, got %s", wantHostname, got.Hostname)
	}
	if got.Domain != "" {
		t.Errorf("got wrong domain should be empty but got %s", got.Domain)
	}
	wantFQDN := fmt.Sprintf("%s.%s", wantHostname, "lan")
	if got.FQDN != wantFQDN {
		t.Errorf("got wrong FQDN, want: %s, got %s", wantFQDN, got.FQDN)
	}
}

func TestHost_FQDN_Domain_NoCgo(t *testing.T) {
	t.SkipNow()
	host, err := newLinuxSystem("").Host()
	if err != nil {
		t.Fatal(fmt.Errorf("could not het host information: %w", err))
	}

	got := host.Info()
	if got.Hostname != wantHostname {
		t.Errorf("hostname want: %s, got %s", wantHostname, got.Hostname)
	}
	if got.Domain != "" {
		t.Errorf("domain should be empty but got %s", got.Domain)
	}
	if got.FQDN != "" {
		t.Errorf("FQDN should empty, got: %s", got.FQDN)
	}
}

func TestHost_FQDN(t *testing.T) {
	tcs := []struct {
		name string
		cf   container.Config
	}{
		{
			name: "debian Cgo with domain",
			cf: container.Config{
				Hostname:     wantHostname,
				Domainname:   wantDomainCgo,
				AttachStderr: testing.Verbose(),
				AttachStdout: testing.Verbose(),
				WorkingDir:   "/usr/src/elastic/go-sysinfo",
				Image:        "golang:1.19-bullseye",
				Cmd: []string{"go", "test", "-v", "-run",
					"^TestHost_FQDN_Domain_Cgo", "./providers/linux"},
				Tty: false,
			},
		},
		{
			name: "debian Cgo no domain",
			cf: container.Config{
				Hostname:     wantHostname,
				AttachStderr: testing.Verbose(),
				AttachStdout: testing.Verbose(),
				WorkingDir:   "/usr/src/elastic/go-sysinfo",
				Image:        "golang:1.19-bullseye",
				Cmd: []string{"go", "test", "-v", "-run",
					"^TestHost_FQDN_No_Domain_Cgo", "./providers/linux"},
				Tty: false,
			},
		},
		{
			name: "debian no Cgo",
			cf: container.Config{
				Hostname:     wantHostname,
				AttachStderr: testing.Verbose(),
				AttachStdout: testing.Verbose(),
				Env:          []string{"CGO_ENABLED=0"},
				WorkingDir:   "/usr/src/elastic/go-sysinfo",
				Image:        "golang:1.19-bullseye",
				Cmd: []string{"go", "test", "-v", "-run",
					"^TestHost_FQDN_Domain_NoCgo", "./providers/linux"},
				Tty: false,
			},
		},
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			runOnDocker(t, cli, &tc.cf)
		})
	}
}

func runOnDocker(t *testing.T, cli *client.Client, cf *container.Config) {
	ctx := context.Background()

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	wd := pwd + "../../../"

	reader, err := cli.ImagePull(ctx, cf.Image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, cf, &container.HostConfig{
		AutoRemove: false,
		Binds:      []string{wd + ":/usr/src/elastic/go-sysinfo"},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	defer func() {
		err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			Force: true, RemoveVolumes: true})
		if err != nil {
			panic(err)
		}
	}()

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case s := <-statusCh:
		if s.StatusCode != 0 {
			var err error
			if s.Error != nil {
				err = fmt.Errorf("container errored: %s", s.Error.Message)
			}
			t.Errorf("conteiner exited with code %d: error: %v", s.StatusCode, err)
		}
		t.Log("docker starts channel:", s)
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStderr: true, ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
