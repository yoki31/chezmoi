package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/coreos/go-semver/semver"

	"github.com/twpayne/chezmoi/v2/pkg/chezmoi"
)

type withSessionTokenType bool

const (
	withSessionToken    withSessionTokenType = true
	withoutSessionToken withSessionTokenType = false
)

type unsupportedVersionError struct {
	version *semver.Version
}

func (e unsupportedVersionError) Error() string {
	return fmt.Sprintf("%s: unsupported version", e.version)
}

var onepasswordVersionRx = regexp.MustCompile(`^(\d+\.\d+\.\d+\S*)`)

type onepasswordConfig struct {
	Command       string
	Prompt        bool
	version       *semver.Version
	versionErr    error
	environFunc   func() []string
	outputCache   map[string][]byte
	sessionTokens map[string]string
}

type onepasswordArgs struct {
	item    string
	vault   string
	account string
	args    []string
}

type onepasswordItemV1 struct {
	Details struct {
		Fields   []map[string]interface{} `json:"fields"`
		Sections []struct {
			Fields []map[string]interface{} `json:"fields,omitempty"`
		} `json:"sections"`
	} `json:"details"`
}

type onepasswordItemV2 struct {
	Fields []map[string]interface{} `json:"fields"`
}

func (c *Config) onepasswordTemplateFunc(userArgs ...string) map[string]interface{} {
	version, err := c.onepasswordVersion()
	if err != nil {
		raiseTemplateError(err)
		return nil
	}

	var baseArgs []string
	switch {
	case version.Major == 1:
		baseArgs = []string{"get", "item"}
	case version.Major >= 2:
		baseArgs = []string{"item", "get", "--format", "json"}
	default:
		raiseTemplateError(unsupportedVersionError{
			version: version,
		})
		return nil
	}

	args, err := newOnepasswordArgs(baseArgs, userArgs)
	if err != nil {
		raiseTemplateError(err)
		return nil
	}

	output, err := c.onepasswordOutput(args, withSessionToken)
	if err != nil {
		raiseTemplateError(err)
		return nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		raiseTemplateError(newParseCmdOutputError(c.Onepassword.Command, args.args, output, err))
		return nil
	}
	return data
}

func (c *Config) onepasswordDetailsFieldsTemplateFunc(userArgs ...string) map[string]interface{} {
	version, err := c.onepasswordVersion()
	if err != nil {
		raiseTemplateError(err)
		return nil
	}

	switch {
	case version.Major == 1:
		item, err := c.onepasswordItemV1(userArgs)
		if err != nil {
			raiseTemplateError(err)
			return nil
		}

		result := make(map[string]interface{})
		for _, field := range item.Details.Fields {
			if designation, ok := field["designation"].(string); ok {
				result[designation] = field
			}
		}
		return result

	case version.Major >= 2:
		item, err := c.onepasswordItemV2(userArgs)
		if err != nil {
			raiseTemplateError(err)
			return nil
		}

		result := make(map[string]interface{})
		for _, field := range item.Fields {
			if _, ok := field["section"]; ok {
				continue
			}
			if label, ok := field["label"].(string); ok {
				result[label] = field
			}
		}
		return result

	default:
		raiseTemplateError(unsupportedVersionError{
			version: version,
		})
		return nil
	}
}

func (c *Config) onepasswordDocumentTemplateFunc(userArgs ...string) string {
	version, err := c.onepasswordVersion()
	if err != nil {
		raiseTemplateError(err)
		return ""
	}

	var baseArgs []string
	switch {
	case version.Major == 1:
		baseArgs = []string{"get", "document"}
	case version.Major >= 2:
		baseArgs = []string{"document", "get"}
	default:
		raiseTemplateError(unsupportedVersionError{
			version: version,
		})
		return ""
	}

	args, err := newOnepasswordArgs(baseArgs, userArgs)
	if err != nil {
		raiseTemplateError(err)
		return ""
	}

	output, err := c.onepasswordOutput(args, withSessionToken)
	if err != nil {
		raiseTemplateError(err)
		return ""
	}
	return string(output)
}

func (c *Config) onepasswordItemFieldsTemplateFunc(userArgs ...string) map[string]interface{} {
	version, err := c.onepasswordVersion()
	if err != nil {
		raiseTemplateError(err)
		return nil
	}

	switch {
	case version.Major == 1:
		item, err := c.onepasswordItemV1(userArgs)
		if err != nil {
			raiseTemplateError(err)
			return nil
		}

		result := make(map[string]interface{})
		for _, section := range item.Details.Sections {
			for _, field := range section.Fields {
				if t, ok := field["t"].(string); ok {
					result[t] = field
				}
			}
		}
		return result

	case version.Major >= 2:
		item, err := c.onepasswordItemV2(userArgs)
		if err != nil {
			raiseTemplateError(err)
			return nil
		}

		result := make(map[string]interface{})
		for _, field := range item.Fields {
			if _, ok := field["section"]; !ok {
				continue
			}
			if label, ok := field["label"].(string); ok {
				result[label] = field
			}
		}
		return result

	default:
		raiseTemplateError(unsupportedVersionError{
			version: version,
		})
		return nil
	}
}

// onepasswordGetOrRefreshSessionToken will return the current session token if
// the token within the environment is still valid. Otherwise it will ask the
// user to sign in and get the new token.
func (c *Config) onepasswordGetOrRefreshSessionToken(args *onepasswordArgs) (string, error) {
	if !c.Onepassword.Prompt {
		return "", nil
	}

	// Check if there's already a valid session token cached in this run for
	// this account.
	sessionToken, ok := c.Onepassword.sessionTokens[args.account]
	if ok {
		return sessionToken, nil
	}

	// If no account has been given then look for any session tokens in the
	// environment.
	if args.account == "" {
		var environ []string
		if c.Onepassword.environFunc != nil {
			environ = c.Onepassword.environFunc()
		} else {
			environ = os.Environ()
		}
		sessionToken = onepasswordUniqueSessionToken(environ)
		if sessionToken != "" {
			return sessionToken, nil
		}
	}

	var commandArgs []string
	if args.account == "" {
		commandArgs = []string{"signin", "--raw"}
	} else {
		sessionToken = os.Getenv("OP_SESSION_" + args.account)
		commandArgs = []string{"signin", args.account, "--raw"}
	}
	if sessionToken != "" {
		commandArgs = append([]string{"--session", sessionToken}, commandArgs...)
	}

	//nolint:gosec
	cmd := exec.Command(c.Onepassword.Command, commandArgs...)
	cmd.Stdin = c.stdin
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	output, err := c.baseSystem.IdempotentCmdOutput(cmd)
	if err != nil {
		return "", newCmdOutputError(cmd, output, err)
	}
	sessionToken = strings.TrimSpace(string(output))

	// Cache the session token in memory, so we don't try to refresh it again
	// for this run for this account.
	if c.Onepassword.sessionTokens == nil {
		c.Onepassword.sessionTokens = make(map[string]string)
	}
	c.Onepassword.sessionTokens[args.account] = sessionToken

	return sessionToken, nil
}

func (c *Config) onepasswordItemV1(userArgs []string) (*onepasswordItemV1, error) {
	args, err := newOnepasswordArgs([]string{"get", "item"}, userArgs)
	if err != nil {
		return nil, err
	}

	output, err := c.onepasswordOutput(args, withSessionToken)
	if err != nil {
		return nil, err
	}

	var item onepasswordItemV1
	if err := json.Unmarshal(output, &item); err != nil {
		return nil, newParseCmdOutputError(c.Onepassword.Command, args.args, output, err)
	}
	return &item, nil
}

func (c *Config) onepasswordItemV2(userArgs []string) (*onepasswordItemV2, error) {
	args, err := newOnepasswordArgs([]string{"item", "get", "--format", "json"}, userArgs)
	if err != nil {
		return nil, err
	}

	output, err := c.onepasswordOutput(args, withSessionToken)
	if err != nil {
		return nil, err
	}

	var item onepasswordItemV2
	if err := json.Unmarshal(output, &item); err != nil {
		return nil, newParseCmdOutputError(c.Onepassword.Command, args.args, output, err)
	}
	return &item, nil
}

func (c *Config) onepasswordOutput(args *onepasswordArgs, withSessionToken withSessionTokenType) ([]byte, error) {
	key := strings.Join(args.args, "\x00")
	if output, ok := c.Onepassword.outputCache[key]; ok {
		return output, nil
	}

	commandArgs := args.args
	if withSessionToken {
		sessionToken, err := c.onepasswordGetOrRefreshSessionToken(args)
		if err != nil {
			return nil, err
		}
		if sessionToken != "" {
			commandArgs = append([]string{"--session", sessionToken}, commandArgs...)
		}
	}

	//nolint:gosec
	cmd := exec.Command(c.Onepassword.Command, commandArgs...)
	cmd.Stdin = c.stdin
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	output, err := c.baseSystem.IdempotentCmdOutput(cmd)
	if err != nil {
		return nil, newCmdOutputError(cmd, output, err)
	}

	if c.Onepassword.outputCache == nil {
		c.Onepassword.outputCache = make(map[string][]byte)
	}
	c.Onepassword.outputCache[key] = output

	return output, nil
}

func (c *Config) onepasswordVersion() (*semver.Version, error) {
	if c.Onepassword.version != nil || c.Onepassword.versionErr != nil {
		return c.Onepassword.version, c.Onepassword.versionErr
	}

	args := &onepasswordArgs{
		args: []string{"--version"},
	}
	output, err := c.onepasswordOutput(args, withoutSessionToken)
	if err != nil {
		c.Onepassword.versionErr = err
		return nil, c.Onepassword.versionErr
	}

	m := onepasswordVersionRx.FindSubmatch(output)
	if m == nil {
		c.Onepassword.versionErr = fmt.Errorf("%q: cannot extract version", bytes.TrimSpace(output))
		return nil, c.Onepassword.versionErr
	}

	version, err := semver.NewVersion(string(m[1]))
	if err != nil {
		c.Onepassword.versionErr = fmt.Errorf("%q: cannot parse version: %w", m[1], err)
		return nil, c.Onepassword.versionErr
	}

	c.Onepassword.version = version
	return c.Onepassword.version, c.Onepassword.versionErr
}

func newOnepasswordArgs(baseArgs, userArgs []string) (*onepasswordArgs, error) {
	if len(userArgs) < 1 || 3 < len(userArgs) {
		return nil, fmt.Errorf("expected 1, 2, or 3 arguments, got %d", len(userArgs))
	}

	a := &onepasswordArgs{
		args: baseArgs,
	}
	a.item = userArgs[0]
	a.args = append(a.args, a.item)
	if len(userArgs) > 1 {
		a.vault = userArgs[1]
		a.args = append(a.args, "--vault", a.vault)
	}
	if len(userArgs) > 2 {
		a.account = userArgs[2]
		a.args = append(a.args, "--account", a.account)
	}
	return a, nil
}

// onepasswordUniqueSessionToken will look for any session tokens in the
// environment. If it finds exactly one then it will return it.
func onepasswordUniqueSessionToken(environ []string) string {
	var token string
	for _, env := range environ {
		key, value, found := chezmoi.CutString(env, "=")
		if found && strings.HasPrefix(key, "OP_SESSION_") {
			if token != "" {
				return ""
			}
			token = value
		}
	}
	return token
}
