package identity

import "time"

var (
	setting *Config

	TimeNow = time.Now
)

func init() {
	Init(Config{
		Token: TokenConfig{
			AccessInterval:  1 * time.Hour,
			RefreshInterval: 7 * 24 * time.Hour,
		},
	})
}

type Config struct {
	Token TokenConfig
}

func Init(cfg Config) {
	setting = &cfg
}
