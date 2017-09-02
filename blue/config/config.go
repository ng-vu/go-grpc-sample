package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/ng-vu/go-grpc-sample/base/l"
)

var ll = l.New()

// Config ...
type Config struct {
	Redis    Redis    `yaml:"redis"`
	Mongo    Mongo    `yaml:"mongo"`
	Postgres Postgres `yaml:"postgres"`

	AgencyService  AgencyService  `yaml:"agency_service"`
	PartnerService PartnerService `yaml:"partner_service"`
	SAdminService  SAdminService  `yaml:"sadmin_service"`

	DocumentPath string `yaml:"document_path"`
	Development  bool   `yaml:"development"`
}

// AgencyService ...
type AgencyService struct {
	GRPC GRPC
	HTTP HTTP
}

// PartnerService ...
type PartnerService struct {
	GRPC GRPC
	HTTP HTTP
}

// SAdminService ...
type SAdminService struct {
	GRPC GRPC
	HTTP HTTP

	MagicToken string `yaml:"magic_token"`
}

// HTTP ...
type HTTP struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Listen ...
func (c HTTP) Listen() string {
	return c.Host + ":" + c.Port
}

// GRPC ...
type GRPC struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Listen ...
func (c GRPC) Listen() string {
	return c.Host + ":" + c.Port
}

// Redis ...
type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// ConnectionString ...
func (c Redis) ConnectionString() string {
	s := ""
	if c.Username == "" || c.Password == "" {
		s = fmt.Sprintf("redis://%s:%s", c.Host, c.Port)
	} else {
		s = fmt.Sprintf("redis://%s:%s@%s:%s", c.Username, c.Password, c.Host, c.Port)
	}
	return s
}

// Postgres ...
type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
	Timeout  int    `yaml:"timeout"`
}

// ConnectionString ...
func (c Postgres) ConnectionString() string {
	sslmode := c.SSLMode
	if c.SSLMode == "" {
		sslmode = "disable"
	}
	if c.Timeout == 0 {
		c.Timeout = 30
	}

	s := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&connect_timeout=%v", c.Username, c.Password, c.Host, c.Port, c.Database, sslmode, c.Timeout)
	return s
}

// Mongo ...
type Mongo struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// ConnectionString returns connection string for mongodb
func (c Mongo) ConnectionString() string {
	conn := "mongodb://"
	if c.Username != "" {
		conn += fmt.Sprintf("%s:%s@", c.Username, c.Password)
	}
	conn += c.Address
	return conn
}

// Default returns default config
func Default() *Config {
	cfg := &Config{
		Redis: Redis{
			Host:     "redis",
			Port:     "6379",
			Username: "",
			Password: "",
		},
		Postgres: Postgres{
			Host:     "postgres",
			Port:     "5432",
			Username: "postgres",
			Password: "postgrespass",
			Database: "blue_dev",
			SSLMode:  "",
			Timeout:  15,
		},
		Mongo: Mongo{
			Address:  "mongo:27017",
			Database: "blue_dev",
			Username: "",
			Password: "",
		},
		AgencyService: AgencyService{
			GRPC: GRPC{Host: "", Port: "3001"},
			HTTP: HTTP{Host: "", Port: "3000"},
		},
		PartnerService: PartnerService{
			GRPC: GRPC{Host: "", Port: "3101"},
			HTTP: HTTP{Host: "", Port: "3100"},
		},
		SAdminService: SAdminService{
			GRPC:       GRPC{Host: "", Port: "3901"},
			HTTP:       HTTP{Host: "", Port: "3900"},
			MagicToken: "9xGnyrmKrhznoyt.ADMIN.mo9gaEbzaGxbWdR8dXIAo=",
		},
		Development: true,
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		ll.Warn("No GOPATH")
	} else {
		cfg.DocumentPath = filepath.Join(gopath, "src/github.com/ng-vu/go-grpc-sample/doc")
	}

	return cfg
}

// DefaultTest returns default config for testing
func DefaultTest() *Config {
	cfg := Default()
	cfg.Postgres.Database = "blue_test"
	cfg.Mongo.Database = "blue_test"

	return cfg
}

// Load loads config from file
func Load(configPath string) (cfg *Config, err error) {
	if configPath == "" {
		return Default(), nil
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return
}

func init() {
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		ll.Fatal("Unable to load timezone", l.Error(err))
	}

	time.Local = location
	ll.Info("Set default timezone", l.String("location", location.String()))
}
