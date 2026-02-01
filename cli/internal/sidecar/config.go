package sidecar

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/crypto/nacl/secretbox"
)

type ConfigManager struct {
	HomeDir string
	Key     [32]byte
}

func NewConfigManager() (*ConfigManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(home, ".shadowprism")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return nil, err
	}

	cm := &ConfigManager{HomeDir: appDir}
	if err := cm.initKey(); err != nil {
		return nil, err
	}

	return cm, nil
}

// initKey generates a machine-specific key so we don't need user prompts.
func (cm *ConfigManager) initKey() error {
	// We use a combination of hostname and OS-specific markers to derive a key
	hostname, _ := os.Hostname()
	machineID := fmt.Sprintf("%s-%s-%s", hostname, runtime.GOOS, runtime.GOARCH)
	
	hash := sha256.Sum256([]byte(machineID))
	copy(cm.Key[:], hash[:])
	return nil
}

func (cm *ConfigManager) SaveSecret(name string, value string) error {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return err
	}

	encrypted := secretbox.Seal(nonce[:], []byte(value), &nonce, &cm.Key)
	
	path := filepath.Join(cm.HomeDir, name+".enc")
	return os.WriteFile(path, []byte(hex.EncodeToString(encrypted)), 0600)
}

func (cm *ConfigManager) LoadSecret(name string) (string, error) {
	path := filepath.Join(cm.HomeDir, name+".enc")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	encrypted, err := hex.DecodeString(string(data))
	if err != nil {
		return "", err
	}

	if len(encrypted) < 24 {
		return "", fmt.Errorf("invalid secret")
	}

	var nonce [24]byte
	copy(nonce[:], encrypted[:24])
	
	decrypted, ok := secretbox.Open(nil, encrypted[24:], &nonce, &cm.Key)
	if !ok {
		return "", fmt.Errorf("decryption failed")
	}

	return string(decrypted), nil
}

func (cm *ConfigManager) GetSocketPath() string {
	return filepath.Join(cm.HomeDir, "engine.sock")
}
