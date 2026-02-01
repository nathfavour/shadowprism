package sidecar

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
)

type Manager struct {
	BinaryPath string
	Port       int
	AuthToken  string
	cmd        *exec.Cmd
	client     *resty.Client
}

func NewManager(port int, token string) *Manager {
	return &Manager{
		Port:      port,
		AuthToken: token,
		client:    resty.New().SetTimeout(2 * time.Second),
	}
}

func (m *Manager) Start(ctx context.Context) error {
	// In a real build, we would extract the embedded binary here.
	// For development, we look for the compiled binary in the core directory.
	binaryName := "shadowprism-core"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	// Look for binary in ../core/target/debug/ or current dir
	cwd, _ := os.Getwd()
	possiblePaths := []string{
		filepath.Join(cwd, "..", "core", "target", "debug", binaryName),
		filepath.Join(cwd, "core", "target", "debug", binaryName),
		filepath.Join(cwd, binaryName),
		filepath.Join("/app", binaryName),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			m.BinaryPath = path
			break
		}
	}

	if m.BinaryPath == "" {
		return fmt.Errorf("core binary not found. please run 'cargo build' in /core")
	}

	m.cmd = exec.CommandContext(ctx, m.BinaryPath)
	m.cmd.Env = append(os.Environ(), 
		fmt.Sprintf("SHADOWPRISM_AUTH_TOKEN=%s", m.AuthToken),
		fmt.Sprintf("PORT=%d", m.Port),
	)

	// Redirect logs for debugging (could be piped to a TUI buffer later)
	m.cmd.Stdout = os.Stdout
	m.cmd.Stderr = os.Stderr

	if err := m.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start core: %w", err)
	}

	// Wait for health check
	return m.waitForReady(ctx)
}

func (m *Manager) waitForReady(ctx context.Context) error {
	cm, err := NewConfigManager()
	if err != nil {
		return err
	}
	client := NewCoreClient(cm.GetSocketPath(), m.AuthToken)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			_, err := client.GetStatus()
			if err == nil {
				return nil
			}
			// Check if process died early
			if m.cmd.ProcessState != nil && m.cmd.ProcessState.Exited() {
				return fmt.Errorf("core process exited prematurely")
			}
		}
	}
}

func (m *Manager) Stop() error {
	if m.cmd != nil && m.cmd.Process != nil {
		return m.cmd.Process.Kill()
	}
	return nil
}
