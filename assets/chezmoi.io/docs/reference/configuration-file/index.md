# Configuration file

chezmoi searches for its configuration file according to the [XDG Base
Directory
Specification](https://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html)
and supports [JSON](https://www.json.org/json-en.html),
[TOML](https://github.com/toml-lang/toml), and [YAML](https://yaml.org/). The
basename of the config file is `chezmoi`, and the first config file found is
used.

In most installations, the config file will be read from
`~/.config/chezmoi/chezmoi.<format>`, where `<format>` is one of `json`,
`toml`, or `yaml`. The config file can be set explicitly with the `--config`
command line option. By default, the format is detected based on the extension
of the config file name, but can be overridden with the `--config-format`
command line option.


## Examples

=== "JSON"

    ```json title="~/.config/chezmoi/chezmoi.json"
    {
        "sourceDir": "/home/user/.dotfiles",
        "git": {
            "autoPush": true
        }
    }
    ```

=== "TOML"

    ```toml title="~/.config/chezmoi/chezmoi.toml"
    sourceDir = "/home/user/.dotfiles"
    [git]
        autoPush = true
    ```

=== "YAML"

    ```yaml title="~/.config/chezmoi/chezmoi.yaml"
    sourceDir: /home/user/.dotfiles
    git:
        autoPush: true
    ```
