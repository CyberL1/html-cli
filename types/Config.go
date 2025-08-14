package types

type Config struct {
	Dev   ConfigDev   `json:"dev"`
	Build ConfigBuild `json:"build"`
}

type ConfigDev struct {
	Port uint16 `json:"port"`
}

type ConfigBuild struct {
	Directory string `json:"directory"`
}
