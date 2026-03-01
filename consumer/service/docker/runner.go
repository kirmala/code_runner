package docker

import (
	"code_processor/http_server/models"
	"context"
	"fmt"
	"io"
	"time"

	"regexp"
	"strings"
	"unicode"

	"code_processor/consumer/cmd/app/config"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Runner struct {
	cli       *client.Client
	imageName string
	resource config.ContainerResource
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
	cli, err := client.NewClientWithOpts(client.WithVersion(clientVersion))
	if err != nil {
		return nil, fmt.Errorf("creating client: %w", err)
	}

	_, err = cli.ImageInspect(context.Background(), imageName)
	if err != nil {
		return nil, fmt.Errorf("inspecting image: %w", err)
	}

	return &Runner{cli: cli, resource: resource, imageName: imageName}, nil
}

func (r *Runner) Run(ctx context.Context, task models.Task) (models.Task, error) {
	var cmd []string
	switch task.Translator {
	case models.PythonTranslator:
		cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > /tmp/code.py && python3 /tmp/code.py", task.Code)}
	case models.GppTranslator:
		cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > /tmp/code.cpp && g++ /tmp/code.cpp -o /tmp/code && /tmp/code", task.Code)}
	case models.ClangTranslator:
		cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > /tmp/code.cpp && clang /tmp/code.cpp -o /tmp/code && /tmp/code", task.Code)}
	default:
		return models.Task{}, models.ErrUnknownTranslator
	}

	resp, err := r.cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: r.imageName,
			Cmd:   cmd,
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory:    r.resource.Memory,
				NanoCPUs:  r.resource.NanoCPUs,
				PidsLimit: r.resource.PidsLimit,
			},
			NetworkMode:    "none",
			ReadonlyRootfs: false,
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		return models.Task{}, fmt.Errorf("creating container: %w", err)
	}

	defer func() {
		_ = r.cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})
	}()

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Second))

	defer cancel()

	if err := r.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return models.Task{}, fmt.Errorf("starting container: %w", err)
	}

	statusCh, errCh := r.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return models.Task{}, fmt.Errorf("running container: %w", err)
		}
	case <-statusCh:
	}

	out, err := r.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return models.Task{}, fmt.Errorf("collecting container logs: %w", err)
	}
	defer func() { _ = out.Close() }()

	output, err := io.ReadAll(out)
	if err != nil {
		return models.Task{}, fmt.Errorf("reading container logs: %w", err)
	}
	cleanOutput := cleanContainerOutput(string(output))

	task.Result = cleanOutput
	task.Status = models.StatusCompleted

	return task, nil
}
