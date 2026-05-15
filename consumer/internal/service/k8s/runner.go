package k8s

import (
    "context"
    "fmt"
    "io"
    "time"

    "github.com/kirmala/code_runner/consumer/internal/domain"
    batchv1 "k8s.io/api/batch/v1"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/resource"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

var translatorImages = map[domain.Translator]string{
    domain.PythonTranslator: "python:3.11-slim",
    domain.GppTranslator:    "gcc:13",
    domain.ClangTranslator:  "silkeh/clang:16",
}

var translatorCommands = map[domain.Translator][]string{
    domain.PythonTranslator: {"sh", "-c", "echo '$CODE' > /tmp/code.py && python3 /tmp/code.py"},
    domain.GppTranslator:    {"sh", "-c", "echo '$CODE' > /tmp/code.cpp && g++ /tmp/code.cpp -o /tmp/code && /tmp/code"},
    domain.ClangTranslator:  {"sh", "-c", "echo '$CODE' > /tmp/code.cpp && clang /tmp/code.cpp -o /tmp/code && /tmp/code"},
}

type Runner struct {
    clientset *kubernetes.Clientset
    namespace string
}

func NewRunner(namespace string) (*Runner, error) {
    // InClusterConfig читает токен и сертификат которые k8s
    // автоматически монтирует в каждый под — не нужен kubeconfig
    cfg, err := rest.InClusterConfig()
    if err != nil {
        return nil, fmt.Errorf("getting in-cluster config: %w", err)
    }

    clientset, err := kubernetes.NewForConfig(cfg)
    if err != nil {
        return nil, fmt.Errorf("creating clientset: %w", err)
    }

    return &Runner{clientset: clientset, namespace: namespace}, nil
}

func (r *Runner) Run(ctx context.Context, task domain.Task) (domain.Task, error) {
    image, ok := translatorImages[task.Translator]
    if !ok {
        return domain.Task{}, fmt.Errorf("unknown translator: %v", task.Translator)
    }

    cmd, ok := translatorCommands[task.Translator]
    if !ok {
        return domain.Task{}, fmt.Errorf("unknown translator command: %v", task.Translator)
    }

    jobName := fmt.Sprintf("runner-%s", task.Id.String())

    // ttlSeconds — сколько секунд Job хранится после завершения
    // потом k8s сам его удалит
    ttl := int32(60)
    backoff := int32(0) // не перезапускать при ошибке

    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name:      jobName,
            Namespace: r.namespace,
        },
        Spec: batchv1.JobSpec{
            TTLSecondsAfterFinished: &ttl,
            BackoffLimit:            &backoff,
            Template: corev1.PodTemplateSpec{
                Spec: corev1.PodSpec{
                    RestartPolicy: corev1.RestartPolicyNever,
                    Containers: []corev1.Container{
                        {
                            Name:    "runner",
                            Image:   image,
                            Command: cmd,
                            Env: []corev1.EnvVar{
                                // передаём код через env чтобы избежать проблем
                                // с кавычками и спецсимволами в shell
                                {Name: "CODE", Value: task.Code},
                            },
                            Resources: corev1.ResourceRequirements{
                                Limits: corev1.ResourceList{
                                    // ограничения как в DockerRunner
                                    corev1.ResourceMemory: resource.MustParse("128Mi"),
                                    corev1.ResourceCPU:    resource.MustParse("500m"),
                                },
                            },
                        },
                    },
                },
            },
        },
    }

    _, err := r.clientset.BatchV1().Jobs(r.namespace).Create(ctx, job, metav1.CreateOptions{})
    if err != nil {
        return domain.Task{}, fmt.Errorf("creating job: %w", err)
    }

    // ждём завершения Job
    if err := r.waitForJob(ctx, jobName); err != nil {
        return domain.Task{}, err
    }

    // читаем логи пода который создал Job
    output, err := r.getJobLogs(ctx, jobName)
    if err != nil {
        return domain.Task{}, err
    }

    task.Result = output
    task.Status = domain.StatusCompleted
    return task, nil
}

func (r *Runner) waitForJob(ctx context.Context, jobName string) error {
    // таймаут 15 секунд на выполнение кода
    ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
    defer cancel()

    for {
        select {
        case <-ctx.Done():
            return fmt.Errorf("job timed out")
        case <-time.After(500 * time.Millisecond):
            job, err := r.clientset.BatchV1().Jobs(r.namespace).Get(ctx, jobName, metav1.GetOptions{})
            if err != nil {
                return fmt.Errorf("getting job status: %w", err)
            }
            if job.Status.Succeeded > 0 {
                return nil
            }
            if job.Status.Failed > 0 {
                return fmt.Errorf("job failed")
            }
        }
    }
}

func (r *Runner) getJobLogs(ctx context.Context, jobName string) (string, error) {
    // находим под созданный этим Job по лейблу
    pods, err := r.clientset.CoreV1().Pods(r.namespace).List(ctx, metav1.ListOptions{
        LabelSelector: fmt.Sprintf("job-name=%s", jobName),
    })
    if err != nil || len(pods.Items) == 0 {
        return "", fmt.Errorf("finding job pod: %w", err)
    }

    req := r.clientset.CoreV1().Pods(r.namespace).GetLogs(pods.Items[0].Name, &corev1.PodLogOptions{})
    logs, err := req.Stream(ctx)
    if err != nil {
        return "", fmt.Errorf("streaming logs: %w", err)
    }
    defer logs.Close()

    output, err := io.ReadAll(logs)
    if err != nil {
        return "", fmt.Errorf("reading logs: %w", err)
    }

    return string(output), nil
}