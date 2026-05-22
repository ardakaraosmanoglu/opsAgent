package executor

import (
    "bytes"
    "context"
    "errors"
    "os/exec"
    "time"
)

type Result struct {
    Stdout   string
    Stderr   string
    ExitCode int
    Duration time.Duration
}

type Runner struct {
    Timeout       time.Duration
    MaxOutputSize int
}

func NewRunner(timeoutSeconds int, maxOutputSizeKB int) *Runner {
    return &Runner{
        Timeout:       time.Duration(timeoutSeconds) * time.Second,
        MaxOutputSize: maxOutputSizeKB * 1024,
    }
}

func (r *Runner) Run(ctx context.Context, command string) (*Result, error) {
    start := time.Now()

    ctx, cancel := context.WithTimeout(ctx, r.Timeout)
    defer cancel()

    cmd := exec.CommandContext(ctx, "bash", "-lc", command)

    var stdout bytes.Buffer
    var stderr bytes.Buffer

    cmd.Stdout = &limitedBuffer{Buffer: &stdout, Limit: r.MaxOutputSize}
    cmd.Stderr = &limitedBuffer{Buffer: &stderr, Limit: r.MaxOutputSize}

    err := cmd.Run()

    if ctx.Err() == context.DeadlineExceeded {
        return nil, errors.New("command timed out")
    }

    exitCode := 0
    if err != nil {
        if exitError, ok := err.(*exec.ExitError); ok {
            exitCode = exitError.ExitCode()
        } else {
            return nil, err
        }
    }

    return &Result{
        Stdout:   stdout.String(),
        Stderr:   stderr.String(),
        ExitCode: exitCode,
        Duration: time.Since(start),
    }, nil
}

type limitedBuffer struct {
    Buffer *bytes.Buffer
    Limit  int
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
    if b.Buffer.Len()+len(p) > b.Limit {
        remaining := b.Limit - b.Buffer.Len()
        if remaining > 0 {
            b.Buffer.Write(p[:remaining])
        }
        return len(p), nil
    }
    return b.Buffer.Write(p)
}
