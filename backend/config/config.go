package config

type Config struct {
	Path        string
	Port        string `toml:"port"`
	OpenBrowser bool   `toml:"open_browser"`
	Updates     bool   `toml:"updates"`
	UpdatesUrl  string `toml:"updates_url"`
}

func NewConfig(path string) *Config {
	return &Config{
		Path:        path,
		Port:        ":4000",
		OpenBrowser: true,
		Updates:     false,
	}
}
