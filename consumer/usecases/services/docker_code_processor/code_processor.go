package dockerCodeProcessor

import (
	"code_processor/consumer/models"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type CodeProcessor struct {
	cli       *client.Client
	imageName string
}

func NewCodeProcessor(imageName string) (*CodeProcessor, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		return nil, fmt.Errorf("creating client: %w", err)
	}

	// buildContext := "."

	// buildContextTar, err := archive.TarWithOptions(buildContext, &archive.TarOptions{})
	// if err != nil {
	// 	return nil, fmt.Errorf("creating tar archive of build context: %w", err)
	// }
	// defer buildContextTar.Close()

	// buildOptions := types.ImageBuildOptions{
	// 	Dockerfile: "Dockerfile",
	// 	Tags:       []string{imageName},
	// 	Remove:     true,
	// }

	// buildResponse, err := cli.ImageBuild(context.Background(), buildContextTar, buildOptions)
	// if err != nil {
	// 	return nil, fmt.Errorf("building image: %w", err)
	// }
	// defer buildResponse.Body.Close()

	return &CodeProcessor{cli: cli, imageName: imageName}, nil
}

func (r *CodeProcessor) Process(task models.Task) (*models.Task, error) {

	var config container.Config
	config.Image = r.imageName
	switch task.Translator {
	case "python":
		config.Cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > /tmp/code.py && python3 /tmp/code.py", task.Code)}
	case "g++":
		config.Cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > /tmp/code.cpp && g++ /tmp/code.cpp -o /tmp/code && /tmp/code", task.Code)}
	case "clang":
		config.Cmd = []string{"sh", "-c", fmt.Sprintf("echo '%s' > /tmp/code.cpp && clang /tmp/code.cpp -o /tmp/code && /tmp/code", task.Code)}
	default:

		return nil, fmt.Errorf("unsupported translator: %s", task.Translator)
	}

	resp, err := r.cli.ContainerCreate(
		context.Background(),
		&config,
		nil, // HostConfig
		nil, // NetworkingConfig
		nil, // Platform
		"",  // Container name
	)
	if err != nil {
		return nil, fmt.Errorf("creating container: %w", err)
	}

	defer func() {
		err := r.cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{})
		if err != nil {
			fmt.Printf("removing container: %v\n", err)
		}
	}()

	if err := r.cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("starting container: %w", err)
	}

	statusCh, errCh := r.cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, fmt.Errorf("running container: %w", err)
		}
	case <-statusCh:
	}

	out, err := r.cli.ContainerLogs(context.Background(), resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return nil, fmt.Errorf("collecting container logs: %w", err)
	}
	defer out.Close()

	output, err := io.ReadAll(out)
	if err != nil {
		return nil, fmt.Errorf("reading container logs: %w", err)
	}

	task.Result = string(output)
	task.Status = "ready"
	return &task, nil
}
