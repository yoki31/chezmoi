package cmd

import (
	"errors"
	"io/fs"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/bradenhilton/mozillainstallhash"
	"howett.net/plist"

	"github.com/twpayne/chezmoi/v2/pkg/chezmoi"
)

type ioregData struct {
	value map[string]interface{}
}

func (c *Config) fromYamlTemplateFunc(s string) interface{} {
	var data interface{}
	if err := chezmoi.FormatYAML.Unmarshal([]byte(s), &data); err != nil {
		raiseTemplateError(err)
		return nil
	}
	return data
}

func (c *Config) includeTemplateFunc(filename string) string {
	var absPath chezmoi.AbsPath
	if filepath.IsAbs(filename) {
		var err error
		absPath, err = chezmoi.NewAbsPathFromExtPath(filename, c.homeDirAbsPath)
		if err != nil {
			raiseTemplateError(err)
		}
	} else {
		absPath = c.SourceDirAbsPath.JoinString(filename)
	}
	contents, err := c.fileSystem.ReadFile(absPath.String())
	if err != nil {
		raiseTemplateError(err)
		return ""
	}
	return string(contents)
}

func (c *Config) ioregTemplateFunc() map[string]interface{} {
	if runtime.GOOS != "darwin" {
		return nil
	}

	if c.ioregData.value != nil {
		return c.ioregData.value
	}

	command := "ioreg"
	args := []string{"-a", "-l"}
	cmd := exec.Command(command, args...)
	output, err := c.baseSystem.IdempotentCmdOutput(cmd)
	if err != nil {
		raiseTemplateError(newCmdOutputError(cmd, output, err))
		return nil
	}

	var value map[string]interface{}
	if _, err := plist.Unmarshal(output, &value); err != nil {
		raiseTemplateError(newParseCmdOutputError(command, args, output, err))
		return nil
	}
	c.ioregData.value = value
	return value
}

func (c *Config) joinPathTemplateFunc(elem ...string) string {
	return filepath.Join(elem...)
}

func (c *Config) lookPathTemplateFunc(file string) string {
	switch path, err := exec.LookPath(file); {
	case err == nil:
		return path
	case errors.Is(err, exec.ErrNotFound):
		return ""
	case errors.Is(err, fs.ErrNotExist):
		return ""
	default:
		raiseTemplateError(err)
		return ""
	}
}

func (c *Config) mozillaInstallHashTemplateFunc(path string) string {
	mozillaInstallHash, err := mozillainstallhash.MozillaInstallHash(path)
	if err != nil {
		raiseTemplateError(err)
		return ""
	}
	return mozillaInstallHash
}

func (c *Config) outputTemplateFunc(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := c.baseSystem.IdempotentCmdOutput(cmd)
	if err != nil {
		raiseTemplateError(newCmdOutputError(cmd, output, err))
		return ""
	}
	// FIXME we should be able to return output directly, but
	// github.com/Masterminds/sprig's trim function only accepts strings
	return string(output)
}

func (c *Config) statTemplateFunc(name string) interface{} {
	switch fileInfo, err := c.fileSystem.Stat(name); {
	case err == nil:
		return map[string]interface{}{
			"name":    fileInfo.Name(),
			"size":    fileInfo.Size(),
			"mode":    int(fileInfo.Mode()),
			"perm":    int(fileInfo.Mode().Perm()),
			"modTime": fileInfo.ModTime().Unix(),
			"isDir":   fileInfo.IsDir(),
		}
	case errors.Is(err, fs.ErrNotExist):
		return nil
	default:
		raiseTemplateError(err)
		return nil
	}
}

func (c *Config) toYamlTemplateFunc(data interface{}) string {
	yaml, err := chezmoi.FormatYAML.Marshal(data)
	if err != nil {
		raiseTemplateError(err)
		return ""
	}
	return string(yaml)
}

func raiseTemplateError(err error) {
	panic(err)
}
