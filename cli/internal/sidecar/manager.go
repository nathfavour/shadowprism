package sidecar

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/nathfavour/shadowprism/cli/internal/embed"
)

type Manager struct {
        BinaryPath string
        Port       int
        AuthToken  string
        Passphrase string
        cmd        *exec.Cmd
        client     *resty.Client
}

func NewManager(port int, token string, passphrase string) *Manager {
        return &Manager{
                Port:       port,
                AuthToken:  token,
                Passphrase: passphrase,
                client:     resty.New().SetTimeout(2 * time.Second),
        }
}

func (m *Manager) Start(ctx context.Context) error {
        // Extract the embedded binary
        binPath, err := embed.ExtractCore()
        if err != nil {
                return fmt.Errorf("failed to extract embedded core: %w", err)
        }
        m.BinaryPath = binPath

        m.cmd = exec.CommandContext(ctx, m.BinaryPath)
        m.cmd.Env = append(os.Environ(), 
                fmt.Sprintf("SHADOWPRISM_AUTH_TOKEN=%s", m.AuthToken),
                fmt.Sprintf("PRISM_PASSPHRASE=%s", m.Passphrase),
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
