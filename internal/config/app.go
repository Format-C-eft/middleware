package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Build information -ldflags .
var (
	branch     string = "dev" //nolint
	commitHash string = "-"   //nolint
	timeBuild  string = "-"   //nolint
)

// GetInfo - get info
func GetInfo() string {
	return fmt.Sprintf("Branch name = %s \nCommit hash = %s \nTime build = %s\n", branch, commitHash, timeBuild)
}

var cfg *Config

// GetConfigInstance - get config app
func GetConfigInstance() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}

type project struct {
	Name  string `yaml:"name"`
	Debug bool   `yaml:"debug"`
	Token struct {
		Password string `yaml:"password"`
	} `yaml:"token"`
	Branch     string
	CommitHash string
	TimeBuild  string
}

type rest struct {
	Host         string   `yaml:"host"`
	Port         uint32   `yaml:"port"`
	Path         string   `yaml:"path"`
	AccessOrigin []string `yaml:"accessOrigin"`
}

type status struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	VersionPath   string `yaml:"versionPath"`
	LivenessPath  string `yaml:"livenessPath"`
	ReadinessPath string `yaml:"readinessPath"`
}

type metrics struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
	Path string `yaml:"path"`
}

type serverOneC struct {
	Path       string        `yaml:"path"`
	MaxTimeout time.Duration `yaml:"maxTimeout"`
	User       struct {
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
	} `yaml:"user"`
	Routes []string `yaml:"routes"`
}

type jaeger struct {
	Use     bool   `yaml:"use"`
	Service string `yaml:"service"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}

type cacheDB struct {
	Host            string `yaml:"host"`
	DB              int    `yaml:"db"`
	Password        string `yaml:"password"`
	KeysTimeExpires struct {
		Session time.Duration `yaml:"session"`
		User    time.Duration `yaml:"user"`
		Cookie  time.Duration `yaml:"cookie"`
	} `yaml:"keysTimeExpires"`
}

// Config - config service
type Config struct {
	Project  project `yaml:"project"`
	Services struct {
		Rest    rest    `yaml:"rest"`
		Status  status  `yaml:"status"`
		Metrics metrics `yaml:"metrics"`
	} `yaml:"services"`
	Servers struct {
		OneC   serverOneC `yaml:"oneC"`
		Jaeger jaeger     `yaml:"jaeger"`
	} `yaml:"servers"`
	Database struct {
		CacheDB cacheDB `yaml:"cachedb"`
	} `yaml:"database"`
}

// ReadConfigYML - read configurations from file and init instance Config.
func ReadConfigYML(configYML string) error {

	if cfg != nil {
		return nil
	}

	if configYML == "" {
		configYML = "config.yml"
	}

	file, err := os.Open(filepath.Clean(configYML))
	if err != nil {
		return err
	}
	defer file.Close() //nolint: gosec

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	if !strings.HasSuffix(cfg.Servers.OneC.Path, "/") {
		cfg.Servers.OneC.Path += "/"
	}

	if !strings.HasSuffix(cfg.Services.Rest.Path, "/") {
		cfg.Services.Rest.Path += "/"
	}

	cfg.Project.Branch = branch
	cfg.Project.CommitHash = commitHash
	cfg.Project.TimeBuild = timeBuild

	return nil
}
