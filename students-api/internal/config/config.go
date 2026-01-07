package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Address string `yaml:"address" env-required:"true"`
}

// Using struct tags for envconfig and yaml parsing
//env-default:"production" // setting it production as default so that while deploying if we forget to set it, it doesn't run in debug mode

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

// MustLoad loads the configuration from environment variables and a YAML file. It is required to succeed, otherwise our application cannot run.
func MustLoad() *Config { // Make sure we don't return error because if config loading fails, app cannot run
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		// Here we define a command-line flag for config path. Like go run cmd/students_api/main.go -config-path=config.yaml
		flags := flag.String("config", "", "Path to the configuration file")
		flag.Parse()
		configPath = *flags // dereference the pointer to get the actual string value

		if configPath == "" {
			log.Fatal("Confiq path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at path: %s", configPath) // Fatalf is for formatted log message
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg) // ReadConfig reads from the file and environment variables. Need to pass memory address of cfg
	if err != nil {
		log.Fatalf("Failed to read config: %v", err.Error()) // .Error() converts error to string
	}

	return &cfg

}
