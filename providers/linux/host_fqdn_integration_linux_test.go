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

//go:build integration

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

func TestHost_FQDN(t *testing.T) {
	// TODO: read GO_VERSION and set the image accordingly
	const image = "golang:1.19-bullseye"
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
				Image:        image,
				Cmd: []string{
					"go", "test", "-v",
					"-tags", "integration,docker",
					"-run", "^TestHost_FQDN_Domain_Cgo",
					"./providers/linux"},
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
				Image:        image,
				Cmd: []string{
					"go", "test", "-v", "-count", "1",
					"-tags", "integration,docker",
					"-run", "^TestHost_FQDN_No_Domain_Cgo",
					"./providers/linux"},
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
				Image:        image,
				Cmd: []string{
					"go", "test", "-v",
					"-tags", "integration,docker",
					"-run", "^TestHost_FQDN_Domain_NoCgo",
					"./providers/linux"},
				Tty: false,
			},
		},
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
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
		t.Fatalf("could not get current directory: %v", err)
	}
	wd := pwd + "../../../"

	reader, err := cli.ImagePull(ctx, cf.Image, types.ImagePullOptions{})
	if err != nil {
		t.Fatalf("failed to pull image %s: %v", cf.Image, err)
	}
	defer reader.Close()
	io.Copy(os.Stderr, reader)

	resp, err := cli.ContainerCreate(ctx, cf, &container.HostConfig{
		AutoRemove: false,
		Binds:      []string{wd + ":/usr/src/elastic/go-sysinfo"},
	}, nil, nil, "")
	if err != nil {
		t.Fatalf("could not create docker conteiner: %v", err)
	}
	defer func() {
		err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			Force: true, RemoveVolumes: true})
		if err != nil {
			t.Logf("WARNING: could not remove docker container: %v", err)
		}
	}()

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		t.Fatalf("could not start docker container: %v", err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			// Not using fatal as we might be able to recover the container
			// logs.
			t.Errorf("docker ContainerWait failed: %v", err)
		}
	case s := <-statusCh:
		if s.StatusCode != 0 {
			msg := fmt.Sprintf("container exited with status code %d", s.StatusCode)
			if s.Error != nil {
				msg = fmt.Sprintf("%s: error: %s", msg, s.Error.Message)
			}

			// Not using fatal as we might be able to recover the container
			// logs.
			t.Errorf(msg)
		}
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStderr: true, ShowStdout: true})
	if err != nil {
		t.Fatalf("could not get container logs: %v", err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
