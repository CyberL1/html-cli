package types

type Config struct {
	Dev ConfigDev `json:"dev"`
}

type ConfigDev struct {
	Root string `json:"root"`
	Port uint16 `json:"port"`
}
