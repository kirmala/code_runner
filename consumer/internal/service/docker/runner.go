package docker

import (
	"context"
	"fmt"
	"io"
	"time"

	"regexp"
	"strings"
	"unicode"

	"github.com/kirmala/code_runner/consumer/cmd/app/config"
	"github.com/kirmala/code_runner/consumer/internal/domain"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type Runner struct {
	cli       *client.Client
	imageName string
	resource  config.ContainerResource
}

func cleanContainerOutput(output string) string {
	ansiEsc := regexp.MustCompile(`\x1B[@-_][0-?]*[ -/]*[@-~]`)
	output = ansiEsc.ReplaceAllString(output, "")

	// 2. Remove non-printable characters (keep newlines and tabs)
	var cleaned strings.Builder
	for _, r := range output {
		if unicode.IsPrint(r) || r == '\n' || r == '\t' || r == '\r' {
			cleaned.WriteRune(r)
		}
	}
	output = cleaned.String()

	// 3. Ensure valid UTF-8
	output = strings.ToValidUTF8(output, "")

	// 4. Trim and ensure non-empty
	output = strings.TrimSpace(output)
	if output == "" {
		output = "[empty output after cleaning]"
	}
	return output
}

func NewRunner(imageName string, clientVersion string, resource config.ContainerResource) (*Runner, error) {
	cli, err := client.New(client.WithAPIVersion(clientVersion))
	if err != nil {
		return nil, fmt.Errorf("creating client: %w", err)
	}

	_, err = cli.ImageInspect(context.Background(), imageName)
	if err != nil {
		return nil, fmt.Errorf("inspecting image: %w", err)
	}

	return &Runner{cli: cli, resource: resource, imageName: imageName}, nil
}

func (r *Runner) Run(ctx context.Context, task domain.Task) (domain.Task, error) {
	var cmd []string
	switch task.Translator {
	case domain.PythonTranslator:
		cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > /tmp/code.py && python3 /tmp/code.py", task.Code)}
	case domain.GppTranslator:
		cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > /tmp/code.cpp && g++ /tmp/code.cpp -o /tmp/code && /tmp/code", task.Code)}
	case domain.ClangTranslator:
		cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > /tmp/code.cpp && clang /tmp/code.cpp -o /tmp/code && /tmp/code", task.Code)}
	}

	resp, err := r.cli.ContainerCreate(
		ctx,
		client.ContainerCreateOptions{
			Config: &container.Config{
				Image: r.imageName,
				Cmd:   cmd,
			},
			HostConfig: &container.HostConfig{
				Resources: container.Resources{
					Memory:    r.resource.Memory,
					NanoCPUs:  r.resource.NanoCPUs,
					PidsLimit: r.resource.PidsLimit,
				},
				NetworkMode:    "none",
				ReadonlyRootfs: false,
			},
		},
	)

	if err != nil {
		return domain.Task{}, fmt.Errorf("creating container: %w", err)
	}

	defer func() {
		_, _ = r.cli.ContainerRemove(ctx, resp.ID, client.ContainerRemoveOptions{})
	}()

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Second))

	defer cancel()

	if _, err := r.cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		return domain.Task{}, fmt.Errorf("starting container: %w", err)
	}

	res := r.cli.ContainerWait(ctx, resp.ID, client.ContainerWaitOptions{Condition: container.WaitConditionNotRunning})
	select {
	case err := <-res.Error:
		if err != nil {
			return domain.Task{}, fmt.Errorf("running container: %w", err)
		}
	case <-res.Result:
	}

	out, err := r.cli.ContainerLogs(ctx, resp.ID, client.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return domain.Task{}, fmt.Errorf("collecting container logs: %w", err)
	}
	defer func() { _ = out.Close() }()

	output, err := io.ReadAll(out)
	if err != nil {
		return domain.Task{}, fmt.Errorf("reading container logs: %w", err)
	}
	cleanOutput := cleanContainerOutput(string(output))

	task.Result = cleanOutput
	task.Status = domain.StatusCompleted

	return task, nil
}
