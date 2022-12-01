module github.com/twpayne/chezmoi/v2

go 1.16

require (
	filippo.io/age v1.0.0
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20220113124808-70ae35bab23f // indirect
	github.com/bmatcuk/doublestar/v4 v4.0.2
	github.com/bradenhilton/mozillainstallhash v1.0.0
	github.com/charmbracelet/glamour v0.5.0
	github.com/coreos/go-semver v0.3.0
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/go-git/go-git/v5 v5.4.2
	github.com/godbus/dbus/v5 v5.0.6 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/go-github/v42 v42.0.0
	github.com/google/gops v0.3.22
	github.com/google/renameio/v2 v2.0.0
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0
	github.com/muesli/combinator v0.3.0
	github.com/muesli/termenv v0.11.0 // indirect
	github.com/pelletier/go-toml v1.9.5
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/rogpeppe/go-internal v1.8.1
	github.com/rs/zerolog v1.26.1
	github.com/sergi/go-diff v1.1.0
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/spf13/afero v1.9.2
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.14.0
	github.com/stretchr/testify v1.8.1
	github.com/twpayne/go-pinentry v0.2.0
	github.com/twpayne/go-shell v0.3.1
	github.com/twpayne/go-vfs/v4 v4.1.0
	github.com/twpayne/go-xdg/v6 v6.0.0
	github.com/xanzy/ssh-agent v0.3.1 // indirect
	github.com/zalando/go-keyring v0.1.1
	go.etcd.io/bbolt v1.3.6
	go.uber.org/multierr v1.8.0
	golang.org/x/oauth2 v0.0.0-20221014153046-6fdb5e3db783
	golang.org/x/sync v0.1.0
	golang.org/x/sys v0.0.0-20220908164124-27713097b956
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
	howett.net/plist v1.0.0
)

exclude github.com/sergi/go-diff v1.2.0 // Produces incorrect diffs
