package conf

type Config struct {
	APIPort    int    `yaml:"ApiPort"`
	DBHost     string `yaml:"DBHost"`
	DBLogin    string `yaml:"DBLogin"`
	DBPassword string `yaml:"DBPassword"`
}
