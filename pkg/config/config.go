package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Environment   string
	APIs          *APIs          `mapstructure:"api"`
	Logging       *Logging       `mapstructure:"logging"`
	Postgres      *Postgres      `mapstructure:"postgresql"`
	KafkaProducer *KafkaProducer `mapstructure:"kafka_producer"`
	KafkaConsumer *KafkaConsumer `mapstructure:"kafka_consumer"`
	TG            *TG            `mapstructure:"telegram"`
}

type APIs struct {
	Worker *API `mapstructure:"worker"`
	UI     *API `mapstructure:"ui"`
}

type API struct {
	Port              int    `mapstructure:"port"`
	StaticContentPath string `mapstructure:"static_content_path"`
}

type TG struct {
	Token string `mapstructure:"token"`
	Host  string `mapstructure:"host"`
}

type Postgres struct {
	OPDB *PostgresSQLConfig `mapstructure:"opdb"`
}

type KafkaConsumer struct {
	Brokers             string        `mapstructure:"brokers"`
	Topic               string        `mapstructure:"topic"`
	SessionTimeout      int           `mapstructure:"session_timeout"`
	PollTimeout         time.Duration `mapstructure:"poll_timeout"`
	Group               string        `mapstructure:"group"`
	BatchSize           int           `mapstructure:"batch_size"`
	ConnectionRetries   int           `mapstructure:"connection_retries"`
	AutoOffsetReset     string        `mapstructure:"auto_offset_reset"`
	OffsetCommitRetries int           `mapstructure:"offset_commit_retries"`
}

type KafkaProducer struct {
	ConnectionName string        `mapstructure:"connection_name"`
	Brokers        string        `mapstructure:"brokers"`
	Topic          string        `mapstructure:"topic"`
	TimeOut        time.Duration `mapstructure:"timeout"`
	Partitions     int           `mapstructure:"partitions"`
	PushTimeout    time.Duration `mapstructure:"push_timeout"`
}

type PostgresSQLConfig struct {
	Alias    string `mapstructure:"alias"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type Logging struct {
	LogLevel string `mapstructure:"loglevel"`
}

func LoadConf(env string) (Config, error) {
	var c Config
	c.loadDefault()
	var confFileName string
	if env == "prod" {
		confFileName = "prod"
	}

	if env == "dev" {
		confFileName = "dev"
		err := c.readEnvironment("../dev.env")
		if err != nil {
			return Config{}, err
		}
	}

	v := viper.New()
	v.SetConfigName(confFileName)
	v.AddConfigPath("../")
	if err := v.ReadInConfig(); err != nil {
		return Config{}, err
	}

	if err := v.Unmarshal(&c); err != nil {
		return Config{}, err
	}

	// Setup environment. Used for adjusting logic for different instances(dev & prod)

	fmt.Println(c.Postgres.OPDB.Alias)
	c.Environment = env
	c.loadEnv()
	return c, nil
}

func (c *Config) loadDefault() {
	var opdb = &PostgresSQLConfig{}
	var postgres = &Postgres{}
	c.Postgres = postgres
	c.Postgres.OPDB = opdb
}

// envVariables assigns variables from Environment
func (c *Config) envVariables(key string) string {
	return os.Getenv(key)
}

func (c *Config) envIntVariables(key string) int {
	value := os.Getenv(key)
	if value == "" {
		return 0
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return int(intValue)
}

// readEnvironment reads the first existing env file from the list
func (c *Config) readEnvironment(files ...string) error {
	for _, f := range files {
		err := godotenv.Load(f)
		if err == nil {
			return nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return nil
}

func (c *Config) loadEnv() {
	c.Postgres.OPDB.Host = c.envVariables("PSQL_HOST")
	c.Postgres.OPDB.Port = c.envIntVariables("PSQL_PORT")
	c.Postgres.OPDB.Username = c.envVariables("PSQL_USERNAME")
	c.Postgres.OPDB.Password = c.envVariables("PSQL_PASSWORD")
	c.Postgres.OPDB.DBName = c.envVariables("PSQL_DBNAME")
	c.TG.Token = c.envVariables("TG_TOKEN")
}
