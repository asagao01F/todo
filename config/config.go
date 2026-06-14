package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// 構造体の定義（以前提示したものと同様）
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	API      APIConfig      `mapstructure:"api"`
	Database DatabaseConfig `mapstructure:"database"`
}

type AppConfig struct {
	Env  string `mapstructure:"env"`
	Name string `mapstructure:"name"`
}

type APIConfig struct {
	Port           int `mapstructure:"port"`
	TimeoutSeconds int `mapstructure:"timeout_seconds"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
}

type PostgresConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	SSLMode      string `mapstructure:"ssl_mode"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// LoadConfig は環境（APP_ENV）に応じて設定ファイルを読み込みます
func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath("./config") // 実行ディレクトリからの相対パス

	// 1. まずデフォルトの共通設定 (config.yaml) を読み込む
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read default config: %w", err)
	}

	// 2. 環境変数「APP_ENV」の値を確認する（なければ "local" とみなす）
	// 環境変数は「TODO_APP_ENV=production」のようにプレフィックスをつけても読めるように設定
	v.SetEnvPrefix("TODO")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	env := v.GetString("APP_ENV")
	if env == "" {
		env = "local" // デフォルトはローカル
	}

	// 3. 環境ごとの設定ファイル (config.local.yaml など) を重ねて読み込む
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	if err := v.MergeInConfig(); err != nil {
		// 本番等でファイルがない場合はエラーにするが、ローカルでファイルがない場合は警告のみにするなど調整可能
		return nil, fmt.Errorf("failed to merge env config (%s): %w", env, err)
	}

	// 4. 構造体にマッピングする
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}