package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppName string `mapstructure:"appname"`

	// service
	Service struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"service"`

	// notion db
	NotionDB struct {
		ClientSecret string `mapstructure:"client_secret"`
		ClientID     string `mapstructure:"client_id"`
		BaseUrl      string `mapstructure:"base_url"`
		Version      string `mapstructure:"version"`

		// dbs
		DBSecret       string   `mapstructure:"db_secret"`
		ContributorsDB string   `mapstructure:"db_contributors"`
		AchievementDB  string   `mapstructure:"db_achievement"`
		TaskDBs        []string `mapstructure:"task_dbs"`
		WorkloadDBs    []string `mapstructure:"workload_dbs"`
		FinDBs         []string `mapstructure:"finance_dbs"`
	} `mapstructure:"notiondb"`

	// everpay
	Everpay struct {
		Url      string `mapstructure:"everpay_url"`
		ScanUrl  string `mapstructure:"scan_url"`
		PrivKey  string `mapstructure:"eth_prvkey"`
		TokenTag string `mapstructure:"token_tag"`
		AppName  string `mapstructure:"app_name"`
	} `mapstructure:"everpay"`

	// log
	Log struct {
		Level string `mapstructure:"level"`
		File  string `mapstructure:"file_name"`
		Save  bool   `mapstructure:"save_file"`
	} `mapstructure:"log"`
}

func New(file string) *Config {
	viper.SetConfigName(file)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("read config failed: %s", err.Error()))
	}

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("Unmarshal config failed: %s", err.Error()))
	}
	return &config
}
