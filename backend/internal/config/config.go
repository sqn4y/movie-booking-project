package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	Server   Server `json:"server"`
	Database DbConf `json:"database"`
}

type DbConf struct {
	Dialect  string `yaml:"dialect"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	SSLMode  string `yaml:"ssl_mode"`
}

type Server struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

func Load() (*Configuration, error) {
	conf := &Configuration{}
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func (s Server) GetAddr() string {
	return fmt.Sprintf(":%s", s.Port)
}

func (d *DbConf) GetConnectionString() string {
	sslMode := d.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		d.Dialect, d.Username, d.Password, d.Host, d.Port, d.Database, sslMode)
}
