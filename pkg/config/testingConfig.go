package config

func NewTestConfig() *Config {
	conf := &Config{
		Environment: "test",
		Logging:     &Logging{LogLevel: "debug"},
		Postgres: &Postgres{
			OPDB: &PostgresSQLConfig{
				Alias:    "opdb",
				Host:     "localhost",
				Port:     5432,
				Username: "postgres",
				Password: "postgres",
				DBName:   "telegram_bot_test",
			},
		},
		TG: &TG{
			Token: "",
			Host:  "",
		},
	}

	return conf
}
