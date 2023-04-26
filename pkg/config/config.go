package config

import (
	"errors"
	"io/fs"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	API         *API        `mapstructure:"api"`
	Logging     *Logging    `mapstructure:"logging"`
	PostgreSQL  *PostgreSQL `mapstructure:"postgresql"`
	TG          *TG         `mapstructure:"telegram"`
}

type API struct {
	WorkerPort int `mapstructure:"worker_port"`
	UIPort     int `mapstructure:"ui_port"`
}

type TG struct {
	Token string `mapstructure:"token"`
	Host  string `mapstructure:"host"`
}

type PostgreSQL struct {
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

	var confFileName string
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
	c.Environment = env
	c.loadEnv()
	return c, nil
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
	// c.PostgreSQL.Host = c.envVariables("PSQL_HOST")
	// c.PostgreSQL.Port = c.envIntVariables("PSQL_PORT")
	c.PostgreSQL.Username = c.envVariables("PSQL_USERNAME")
	c.PostgreSQL.Password = c.envVariables("PSQL_PASSWORD")
	// c.PostgreSQL.DBName = c.envVariables("PSQL_DBNAME")
	c.TG.Token = c.envVariables("TG_TOKEN")
}
