package redis

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	TTL      int64  `yaml:"TTL"`
}
