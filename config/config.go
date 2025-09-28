package config

import (
	"strings"

	"github.com/spf13/viper"
)

// AppConfig 定義應用程式所有需要的設定結構
type AppConfig struct {
	Server   ServerConfig
	Database DBConfig
	Unipile  UnipileConfig
	App      AppURLConfig
}

// ServerConfig 伺服器相關設定
type ServerConfig struct {
	Port      int    `mapstructure:"port"`
	JWTSecret string `yaml:"jwt_secret"`
}

// DBConfig 資料庫相關設定
type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// UnipileConfig Unipile 服務相關設定
type UnipileConfig struct {
	APIKey     string `mapstructure:"api_key"`
	APIBaseURL string `mapstructure:"api_base_url"`
}

// AppURLConfig 應用程式 URL 設定
type AppURLConfig struct {
	ServerURL   string `mapstructure:"server_url"`
	FrontendURL string `mapstructure:"frontend_url"`
}

// LoadConfig 載入 config.yml 檔案
func LoadConfig() (*AppConfig, error) {
	viper.AddConfigPath("./config") // 在當前目錄查找
	viper.SetConfigName("config")   // 配置文件名 (不含擴展名)
	viper.SetConfigType("yml")      // 配置文件類型

	// 允許從環境變數讀取 (例如 SERVER_PORT)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 嘗試讀取配置
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg AppConfig

	// 反序列化到結構體
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
