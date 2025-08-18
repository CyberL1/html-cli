package constants

import "html-cli/types"

var (
	Version string
	Config  types.Config
)

var DefaultConfig = types.Config{
	Dev: types.ConfigDev{
		Port: 8080,
	},
	Build: types.ConfigBuild{
		Directory: "build",
	},
}
