package cmd

import (
	"encoding/json"
	"os/exec"
)

type vaultConfig struct {
	Command string
	cache   map[string]interface{}
}

func (c *Config) vaultTemplateFunc(key string) interface{} {
	if data, ok := c.Vault.cache[key]; ok {
		return data
	}

	args := []string{"kv", "get", "-format=json", key}
	//nolint:gosec
	cmd := exec.Command(c.Vault.Command, args...)
	cmd.Stdin = c.stdin
	cmd.Stderr = c.stderr
	output, err := c.baseSystem.IdempotentCmdOutput(cmd)
	if err != nil {
		raiseTemplateError(newCmdOutputError(cmd, output, err))
		return nil
	}

	var data interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		raiseTemplateError(newParseCmdOutputError(c.Vault.Command, args, output, err))
		return nil
	}

	if c.Vault.cache == nil {
		c.Vault.cache = make(map[string]interface{})
	}
	c.Vault.cache[key] = data

	return data
}
